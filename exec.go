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

	lua "github.com/yuin/gopher-lua"
)

const (
	//max number of table subcalls before we punt on possible recursion
	//this should ultimately be a config param
	defaultMaxCallDepth = 25
)

type workPackage struct {
	repo      TableRepository
	table     *table.Table
	script    *lua.FunctionProto
	operation string
	count     int
}

type executionEngine struct {
	callDepth int //number of table calls - prevent recursion with a hammer
	rnd       *rand.Rand
}

func newExecutionEngine() *executionEngine {
	return &executionEngine{
		callDepth: 0,
		rnd:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (ee *executionEngine) execute(wp *workPackage, tr *res.TableResult) {
	ee.executeInternal(wp, tr)
}

func (ee *executionEngine) executeInternal(wp *workPackage, tr *res.TableResult) {
	//TODO: eventually may want to call a goroutine and use contexts for timeout
	//but skip that for now and see how things perform. Might be important in public
	//servers to timeout an ill-behaved lua script or recursive table issue

	switch wp.operation {
	case "roll":
		tr.AddLog(fmt.Sprintf("Executing Roll on table: %s count: %d", wp.table.Definition.Name, wp.count))
		for i := 1; i <= wp.count; i++ {
			ee.executeRoll(wp, tr)
		}
	case "pick":
		tr.AddLog(fmt.Sprintf("NOOP: Executing Pick on table: %s count: %d", wp.table.Definition.Name, wp.count))
	case "script":
		tr.AddLog(fmt.Sprintf("NOOP: Executing Script: FIXMYNAME count: %d", wp.count))
	}
}

func (ee *executionEngine) executeRoll(wp *workPackage, tr *res.TableResult) {

	//check call depth - will rolling here push us over?
	ee.callDepth++
	if ee.callDepth > defaultMaxCallDepth {
		tr.AddLog(fmt.Sprintf("Unable to roll on table, max call depth of: %d exceeded", defaultMaxCallDepth))
	}

	//roll on the table
	var rolledValue int
	var buf string
	switch wp.table.Definition.TableType {
	case "flat": //flat tables need a dice parse result, range tables already have one
		dpr := &dice.ParseResult{
			Count:   1,
			DieType: len(wp.table.RawContent),
		}
		dp := make([]*dice.ParseResult, 1, 1)
		dp[0] = dpr
		rolledValue = ee.rollDice(dp)
		tr.AddLog(fmt.Sprintf("Rolled: %d", rolledValue))
		buf = wp.table.RawContent[rolledValue-1]
	case "range":
		rolledValue = ee.rollDice(wp.table.Definition.DiceParsed)
		tr.AddLog(fmt.Sprintf("Rolled: %d", rolledValue))
		buf = ee.rangeResultFromRoll(wp, rolledValue)
	}

	//at this point, we have a random string stored in the buf string - but it may
	//need expansion for each table it references
	bufParts, exists := util.FindNextTableRef(buf)
	if !exists { //expansion not needed
		tr.AddResult(buf)
		return
	}

	//here if need to expand at least one tableref
	var sb strings.Builder
	for exists { //there is at least one table ref remaining
		sb.WriteString(bufParts[0]) //everything up to the first reference

		//need to recurse here - build new work pkg & create new TableResult
		//On return, need to add the string from TableResult to the String Builder
		//need to add the Log entries from the TabbleResult to this TableResult

		nextWp := &workPackage{
			repo:   wp.repo,
			script: nil,
		}
		safeAndSane := false

		//what type of table ref do we have - build rest of workPkg...
		if table.ExternalCalledPattern.MatchString(bufParts[1]) {
			extMatches := table.ExternalCalledPattern.FindStringSubmatch(bufParts[1])
			nextWp.count = 1 //always roll once per external tables
			nextWp.operation = "roll"
			tableRef, err := wp.repo.TableForName(extMatches[1])
			if err != nil {
				tr.AddLog(fmt.Sprintf("%v", err))
				return
			}
			nextWp.table = tableRef
			safeAndSane = true
		}
		if table.InlineCalledPattern.MatchString(bufParts[1]) {
			extMatches := table.InlineCalledPattern.FindStringSubmatch(bufParts[1])
			nextWp.count = 1 //always roll once per internal tables
			nextWp.operation = "roll"
			tablename := util.BuildFullName(wp.table.Definition.Name, extMatches[1])
			tableRef, err := wp.repo.TableForName(tablename)
			if err != nil {
				tr.AddLog(fmt.Sprintf("%v", err))
				return
			}
			nextWp.table = tableRef
			safeAndSane = true
		}
		if table.PickCalledPattern.MatchString(bufParts[1]) {
			extMatches := table.PickCalledPattern.FindStringSubmatch(bufParts[1])
			nextWp.count, _ = strconv.Atoi(extMatches[1]) //no err per regex
			nextWp.operation = "pick"
			tableRef, err := wp.repo.TableForName(extMatches[2])
			if err != nil {
				tr.AddLog(fmt.Sprintf("%v", err))
				return
			}
			nextWp.table = tableRef
			safeAndSane = true
		}

		//should never happen but check anyway - if we fail here, tests and
		//table parsing logic has gone wrong
		if !safeAndSane {
			tr.AddLog(fmt.Sprintf("Unexpected table ref. This should NEVER happen: %s", bufParts[1]))
			return
		}

		//recurse to expand the first ref found in bufParts
		nextTr := res.NewTableResult()
		ee.executeInternal(nextWp, nextTr)

		//capture the results of the recursion - add expanded string and log(s)
		//only one result is expected unless there was a multi-pick call so
		//handle all cases via append. We expect many logs to be present in nextTr tho
		tr.Result = append(tr.Result, nextTr.Result...)
		tr.Log = append(tr.Log, nextTr.Log...)

		//tack on that part of the buffer that hasnt been checked for references
		sb.WriteString(bufParts[2])

		bufParts, exists = util.FindNextTableRef(sb.String()) //keep expanding
	}
	tr.AddResult(sb.String())
}

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

func (ee *executionEngine) rollDice(dpr []*dice.ParseResult) int {
	//this algo could be made faster. It is written this way for debugging and
	//clarity. Optimizing for speed is not a requirement of this lib

	//store the result of each die expression in a slice
	values := make([]int, len(dpr), len(dpr))
	ops := make([]string, len(dpr), len(dpr))
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

	//utilize the operators between each die expression to calc the total
	total := 0
	for i := 0; i < len(values); i++ {
		if i == 0 {
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
