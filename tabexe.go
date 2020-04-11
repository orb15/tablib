package tablib

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"tablib/dice"
	"tablib/table"
	res "tablib/tableresult"
	"tablib/util"
	"time"
)

const (
	//max number of table subcalls before we punt on possible recursion
	//TODO: this should ultimately be a config param
	defaultMaxCallDepth = 25
)

type workPackage struct {
	nameSvc   nameResolver
	table     *table.Table
	operation string
	count     int
	pickCount int
}

type executionEngine struct {
	callDepth int //number of table calls - prevent malicious or inadvertent circular refs with a hammer
	rnd       *rand.Rand
}

func newExecutionEngine() *executionEngine {
	return &executionEngine{
		callDepth: 0,
		rnd:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (ee *executionEngine) execute(wp *workPackage, tr *res.TableResult) {

	//quick check on sanity
	if wp.count > defaultMaxCallDepth {
		tr.AddLog(fmt.Sprintf("Too many rolls requested, max is: %d", defaultMaxCallDepth))
		return
	}

	for i := 1; i <= wp.count; i++ {
		generated := ee.executeInternal(wp, tr)
		tr.AddResult(generated)
		ee.callDepth = 0 //calldepth resets since this is a new roll/pick attempt
	}
}

//this function is recursively called to expand table refs. It is called initially
//based on a roll or pick request, but if that request's result contains a
//tableref, this function is called to fulfill that reference. This process
//happes recursively until all tablerefs are resolved. See 'expandAllTableRefs()
//below to see how this function is called recursively
func (ee *executionEngine) executeInternal(wp *workPackage, tr *res.TableResult) string {
	//TODO: eventually may want to call a goroutine and use contexts for timeout
	//but skip that for now and see how things perform. Might be important in public
	//servers to timeout an ill-behaved lua script or recursive table issue

	var generated string
	switch wp.operation {
	case table.OpRoll:
		tr.AddLog(fmt.Sprintf("Executing Roll on table: %s", wp.table.Definition.Name))
		generated = ee.executeRoll(wp, tr)
	case table.OpPick:
		tr.AddLog(fmt.Sprintf("Executing Pick %d on table: %s ", wp.pickCount, wp.table.Definition.Name))
		generated = ee.executePick(wp, tr)
	}

	return generated
}

//picks n unique rows from a flat table
func (ee *executionEngine) executePick(wp *workPackage, tr *res.TableResult) string {

	//check call depth - will rolling here push us over?
	if !ee.checkCallDepth(tr) {
		return "Call depth exceeded!"
	}

	//picking on range tables is not allowed
	if wp.table.Definition.TableType == table.TypeRange {
		tr.AddLog(fmt.Sprintf("Pick requested on ranged table: %s", wp.table.Definition.Name))
		return "Pick on range table not allowed"
	}

	//if asking for more picks than content, return content and a warning
	if wp.pickCount >= len(wp.table.RawContent) {
		tr.AddLog(fmt.Sprintf("Pick %d on table: %s requested but it has only %d entries",
			wp.count, wp.table.Definition.Name, len(wp.table.RawContent)))
		return strings.Join(wp.table.RawContent, ",")
	}

	//create a tracking slice to track picked values
	//TODO: Can this be made more efficient - probably but efficiency will be in
	//part determined by the len(content) and count.
	type pick struct {
		v string
		p bool
	}
	pickSlice := make([]*pick, 0, wp.pickCount)
	for _, c := range wp.table.RawContent {
		pickSlice = append(pickSlice, &pick{v: c, p: false})
	}

	//do the picks
	remaining := wp.pickCount
	length := len(pickSlice)
	outSlice := make([]string, 0, wp.pickCount)
	for remaining > 0 {
		picked := ee.rnd.Intn(length)
		if !pickSlice[picked].p {
			pickSlice[picked].p = true
			remaining--
			outSlice = append(outSlice, pickSlice[picked].v)
		}
	}
	buf := strings.Join(outSlice, ",")
	return ee.expandAllRefs(buf, wp, tr) //recurse in case this generate table refs
}

//randomly selects a row from a flat or range table
func (ee *executionEngine) executeRoll(wp *workPackage, tr *res.TableResult) string {

	//check call depth - will rolling here push us over?
	if !ee.checkCallDepth(tr) {
		return ""
	}

	//roll on the table
	rolledValue := ee.rollDice(wp.table.Definition.DiceParsed)
	tr.AddLog(fmt.Sprintf("Rolled: %d", rolledValue))

	//interpret the roll based on table type
	var buf string
	switch wp.table.Definition.TableType {
	case table.TypeFlat:
		buf = wp.table.RawContent[rolledValue-1]
	case table.TypeRange:
		buf = ee.rangeResultFromRoll(wp, rolledValue)
	}
	return ee.expandAllRefs(buf, wp, tr)
}

//works in conjunction with executeInternal to recursively expand all tablerefs
func (ee *executionEngine) expandAllRefs(buf string, wp *workPackage, tr *res.TableResult) string {

	//buf may need tablerefs expanded, check that first
	bufParts, exists := util.FindNextTableRef(buf)
	if !exists { //expansion not needed
		return buf
	}

	//here if need to expand at least one tableref
	var sb strings.Builder
	for exists { //there is at least one table ref remaining
		sb.Reset()
		sb.WriteString(bufParts[0]) //everything up to the first reference

		//need to recurse here so set up the new work package's common elements
		nextWp := &workPackage{
			nameSvc: wp.nameSvc,
		}
		safeAndSane := false //sanity checker - programming mistake trap

		//what type of table ref do we have - build rest of workPkg...
		if table.ExternalCalledPattern.MatchString(bufParts[1]) {
			extMatches := table.ExternalCalledPattern.FindStringSubmatch(bufParts[1])
			nextWp.count = 1 //always roll once per external tables
			nextWp.operation = table.OpRoll
			tableRef, err := wp.nameSvc.tableForName(extMatches[1])
			if err != nil {
				tr.AddLog(fmt.Sprintf("%v", err))
				sb.WriteString(fmt.Sprintf(" --BADREF: %s--", bufParts[1]))
				return sb.String()
			}
			nextWp.table = tableRef
			safeAndSane = true
		}
		if table.InlineCalledPattern.MatchString(bufParts[1]) {
			extMatches := table.InlineCalledPattern.FindStringSubmatch(bufParts[1])
			nextWp.count = 1 //always roll once on internal tables
			nextWp.operation = table.OpRoll
			tablename := util.BuildFullName(wp.table.Definition.Name, extMatches[1])
			tableRef, err := wp.nameSvc.tableForName(tablename)
			if err != nil {
				tr.AddLog(fmt.Sprintf("%v", err))
				sb.WriteString(fmt.Sprintf(" --BADREF: %s--", bufParts[1]))
				return sb.String()
			}
			nextWp.table = tableRef
			safeAndSane = true
		}
		if table.PickCalledPattern.MatchString(bufParts[1]) {
			extMatches := table.PickCalledPattern.FindStringSubmatch(bufParts[1])
			nextWp.pickCount, _ = strconv.Atoi(extMatches[1]) //no err per regex
			nextWp.count = 1                                  //always roll once on pick requests
			nextWp.operation = table.OpPick
			tableRef, err := wp.nameSvc.tableForName(extMatches[2])
			if err != nil {
				tr.AddLog(fmt.Sprintf("%v", err))
				sb.WriteString(fmt.Sprintf(" --BADREF: %s--", bufParts[1]))
				return sb.String()
			}
			nextWp.table = tableRef
			safeAndSane = true
		}

		//should never happen but check anyway - if we fail here, tests and
		//table parsing logic have gone wrong - fix yer code!
		if !safeAndSane {
			msg := fmt.Sprintf("Unexpected table ref. This should NEVER happen: %s", bufParts[1])
			tr.AddLog(msg)
			return msg
		}

		//recurse to expand the first ref found in bufParts
		generated := ee.executeInternal(nextWp, tr)

		//capture the results of the recursion
		sb.WriteString(generated)

		//tack on that part of the buffer that hasnt been checked for references
		sb.WriteString(bufParts[2])

		//assemble the string - do we still have tablerefs to be expanded?
		replacedString := sb.String()
		bufParts, exists = util.FindNextTableRef(replacedString)
	}
	return sb.String()
}

//use the result of a roll to determine which ranged content item should be returned
func (ee *executionEngine) rangeResultFromRoll(wp *workPackage, roll int) string {
	for _, rc := range wp.table.RangeContent {
		if roll >= rc.Low && roll <= rc.High {
			return rc.Content
		}
	}

	//we should never get here. This is a fallback message to let us know something
	//in the code is seriously borked. Something is wrong with the dice roller, dice
	//parser or the ranged table parser for this to happen. Also check your tests,
	//something got missed there too
	msg := fmt.Sprintf("In ee.rangeResultFromRoll, roll of %d exceeded bounds of table: %s",
		roll, wp.table.Definition.Name)
	return msg
}

//executes a dice roll as specified in the dice parsed result
func (ee *executionEngine) rollDice(dpr []*dice.ParseResult) int {
	//this algo could be made faster. It is written this way for debugging and
	//clarity. Optimizing for speed is not a requirement of this lib at least for now

	//step 1: roll each die expression and push those results independently
	//into a slice. Furthermore, gather the operation meant to join each die result
	//to the one following

	//store the result of each die expression in a slice
	values := make([]int, len(dpr), len(dpr)) //each roll's
	ops := make([]string, len(dpr), len(dpr)) //its assoc operator
	i := 0
	for _, d := range dpr {
		//handle constant at end eg 2d6 + 3
		if d.DieType == 0 { //a die type of 0 means Count holds the constant
			values[i] = d.Count
			ops[i] = d.Operator //must be "none"
			i++
			continue
		}
		thisRoll := 0
		for j := 1; j <= d.Count; j++ {
			rndval := ee.rnd.Intn(d.DieType) + 1
			thisRoll += rndval
		}
		values[i] = thisRoll
		ops[i] = d.Operator
		i++
	}

	//step2: utilize the operators between each die expression to calc the total
	total := 0
	for i := 0; i < len(values); i++ {
		if i == 0 { //the first value rolled is always the initial total
			total += values[0]
			continue
		}
		switch ops[i-1] {
		case "+":
			total += values[i]
		case "-":
			total -= values[i]
		case "*":
			total = total * values[i]
		}
	}
	return total
}

//determine if we have had too many table refs and need to punt. This is a
//brute force block to circular table dependencies (malicious or otherwise) as
//only so many roll or picks are allowed before we stop resolving lookups
func (ee *executionEngine) checkCallDepth(tr *res.TableResult) bool {
	ee.callDepth++
	if ee.callDepth > defaultMaxCallDepth {
		tr.AddLog(fmt.Sprintf("Unable to roll on table, max call depth of: %d exceeded", defaultMaxCallDepth))
		return false
	}
	return true
}
