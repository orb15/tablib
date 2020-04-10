package tablib

import (
	"fmt"
	"strings"
	"sync"
	"tablib/dice"
	"tablib/table"
	"tablib/tableresult"
	"tablib/util"
	"tablib/validate"

	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"

	yaml "gopkg.in/yaml.v2"
)

type tableData struct {
	yamlSource  string
	parsedTable *table.Table
}

type scriptData struct {
	scriptSource string
	parsedScript *lua.FunctionProto
}

type concreteTableRepo struct {
	tableStore  map[string]*tableData
	scriptStore map[string]*scriptData
	lock        *sync.RWMutex
}

func (cr *concreteTableRepo) AddTable(yamlBytes []byte) (*validate.ValidationResult, error) {

	//Note: not locking repo here so parse + validate can be multithreaded if caller desires

	//is this even valid YAML?
	var tbl *table.Table
	err := yaml.Unmarshal(yamlBytes, &tbl)
	if err != nil {
		return nil, err
	}

	//validate the table and parse portions of it since we are tearing the table
	//apart to do the validation anyway
	validationResults := tbl.Validate()

	//by definition, tables that arrive here are not inline tables
	tbl.IsInlineTable = false

	//do not proceed if the table is invalid (but its ok if there are warnings)
	if !validationResults.Valid() {
		return validationResults, nil
	}

	//add dice information to flat tables since we need to roll on them
	if tbl.Definition.TableType == table.TypeFlat {
		addDiceParseResultForFlatAndInlineTables(tbl)
	}

	//build out inline tables in this table as first-class tables, then validate
	//the inline content for proper tablerefs
	var inlines []*table.Table
	if len(tbl.Inline) > 0 {
		inlines = extractInlineTables(tbl)
		for _, ilt := range inlines {
			ilt.ValidateContent(validationResults)
		}
	}

	//final validity check to prevent storing a bad table in the repo
	if !validationResults.Valid() {
		return validationResults, nil
	}

	//lock the repo now since we will write to it
	cr.lock.Lock()
	defer cr.lock.Unlock()

	//put the valid table in the repo
	fullName := util.BuildFullName(tbl.Definition.Name, "")
	cr.tableStore[fullName] = &tableData{
		yamlSource:  string(yamlBytes),
		parsedTable: tbl,
	}

	//store the inline tables for this table as first-class tables
	for _, ilt := range inlines {
		cr.tableStore[ilt.Definition.Name] = &tableData{
			yamlSource:  "",
			parsedTable: ilt,
		}
	}

	return validationResults, nil
}

func (cr *concreteTableRepo) AddLuaScript(scriptName, luaScript string) error {

	//not locking repo here so compilation can be multithreaded if caller desires

	//read and compile the lua script
	reader := strings.NewReader(luaScript)
	astStatements, err := parse.Parse(reader, luaScript)
	if err != nil {
		return err
	}
	proto, err := lua.Compile(astStatements, luaScript)
	if err != nil {
		return err
	}

	//lock the repo now since we will write to it
	cr.lock.Lock()
	defer cr.lock.Unlock()

	//store the Lua script bytecode in the repo
	cr.scriptStore[scriptName] = &scriptData{
		scriptSource: luaScript,
		parsedScript: proto,
	}
	return nil
}

func (cr *concreteTableRepo) List(objectName string) (string, error) {
	cr.lock.RLock()
	defer cr.lock.RUnlock()
	return "", nil
}
func (cr *concreteTableRepo) Search(namePredicate string, tags []string) []*SearchResult {
	cr.lock.RLock()
	defer cr.lock.RUnlock()
	return make([]*SearchResult, 0)
}
func (cr *concreteTableRepo) Roll(tableName string, execsDesired int) *tableresult.TableResult {
	cr.lock.RLock()
	defer cr.lock.RUnlock()

	tr := tableresult.NewTableResult()
	tbl, found := cr.tableStore[tableName]
	if !found {
		tr.AddLog(fmt.Sprintf("Table: %s does not exist", tableName))
		return tr
	}

	wp := &workPackage{
		repo:      cr,
		table:     tbl.parsedTable,
		script:    nil,
		operation: table.OpRoll,
		count:     execsDesired,
	}
	exeng := newExecutionEngine()
	exeng.execute(wp, tr)
	return tr
}

func (cr *concreteTableRepo) Pick(tableName string, count int) *tableresult.TableResult {
	cr.lock.RLock()
	defer cr.lock.RUnlock()

	tr := tableresult.NewTableResult()
	tbl, found := cr.tableStore[tableName]
	if !found {
		tr.AddLog(fmt.Sprintf("Table: %s does not exist", tableName))
		return tr
	}

	wp := &workPackage{
		repo:      cr,
		table:     tbl.parsedTable,
		script:    nil,
		operation: table.OpRoll,
		count:     1,
		pickCount: count,
	}
	exeng := newExecutionEngine()
	exeng.execute(wp, tr)
	return tr
}

func (cr *concreteTableRepo) Execute(scriptName string) (map[string]string, error) {
	cr.lock.RLock()
	defer cr.lock.RUnlock()

	//TODO: move all this to the execution engine, including helper methods

	//set up a new lua VM
	//TODO: limit call stack and repository sizes, prevent use of lua modules
	//that allow access to OS, filesys or other dangerous crap
	lState := lua.NewState()
	defer lState.Close()

	//tell the lua VM about the go code we are exposing to it
	luaMod := newLuaModule(cr)
	lState.PreloadModule("tables", luaMod.luaModuleLoader)

	//fetch the precompiled lua script by name
	scriptData, found := cr.scriptStore[scriptName]
	if !found {
		return map[string]string{"Error": "Script not found!"}, nil
	}

	//prep the lua script - all this does is store the precompiled code
	//in lState and await our call to lua functions it defines
	luafunc := lState.NewFunctionFromProto(scriptData.parsedScript)
	lState.Push(luafunc)
	err := lState.PCall(0, lua.MultRet, nil)
	if err != nil {
		return createErrorMap(scriptName,
			fmt.Sprintf("fail to load loading compiled script: %s", err)), nil
	}

	//TODO: here we need to call well-known lua function to get info about
	//the params the lua main() code needs to do its job. Once we get these, this
	//method (Execute) will need to utilize a callback function (needs passed in)
	//to request the param values from the caller of this lib.

	//For sanity sake, all lua functions should take and return a single well-known
	//type so we always know the size of the argument list being passed or
	//returned. A map[string]string is sufficent and simple to handle

	ldm := make(map[string]string) //hack: make up params to pass for now
	ldm["p1"] = "v1"
	ldm["p2"] = "v2"

	//call the well-known function "main" which is the 'main' for our lua script
	if err := lState.CallByParam(lua.P{
		Fn:      lState.GetGlobal("main"),
		NRet:    0,
		Protect: true,
	}, toLuaLTable(ldm)); err != nil {
		if err != nil {
			return createErrorMap(scriptName, fmt.Sprintf("executing main(): %s", err)), nil
		}
	}

	//retrieve the well-known return value from lua
	retval := lState.GetGlobal("rettbl")
	retmap := fromLuaTable(scriptName, lState, retval)

	return retmap, nil
}

//TableForName returns the underlying table for give table name
func (cr *concreteTableRepo) TableForName(name string) (*table.Table, error) {
	cr.lock.RLock()
	defer cr.lock.RUnlock()

	tbl, found := cr.tableStore[name]
	if !found {
		return nil, fmt.Errorf("Table does not exist: %s", name)
	}
	return tbl.parsedTable, nil
}

//for each inline table in a table, create a full-featured table. this
//enables inline lookups to be handled like any other table ref during execution
func extractInlineTables(mainTable *table.Table) []*table.Table {
	inlinesAsTables := make([]*table.Table, 0, len(mainTable.Inline))
	for _, ilt := range mainTable.Inline {

		def := &table.DefinitionPart{
			Name:      ilt.FullyQualifiedName,
			TableType: table.TypeFlat,
		}

		content := make([]string, 0, len(ilt.Content))
		for _, c := range ilt.Content {
			content = append(content, c)
		}

		tbl := &table.Table{
			Definition:    def,
			RawContent:    content,
			IsValid:       true,
			IsInlineTable: true,
		}

		//add dice info to this inline table since we need to roll on it
		addDiceParseResultForFlatAndInlineTables(tbl)
		inlinesAsTables = append(inlinesAsTables, tbl)
	}
	return inlinesAsTables
}

//creates parsed dice info to make rolling on the flat or inline
//table easier at exection - now these table types have dice info
//just as if they were ranged tables
func addDiceParseResultForFlatAndInlineTables(tbl *table.Table) {
	dpr := &dice.ParseResult{
		Count:   1,
		DieType: len(tbl.RawContent),
	}
	dp := make([]*dice.ParseResult, 1, 1)
	dp[0] = dpr
	tbl.Definition.DiceParsed = dp
}

//converts a go map to a lua LTable
func toLuaLTable(goMap map[string]string) *lua.LTable {
	ltbl := &lua.LTable{}
	for k, v := range goMap {
		ltbl.RawSetString(k, lua.LString(v))
	}
	return ltbl
}

//converts a lua LTable to a go map
func fromLuaTable(scriptName string, lState *lua.LState, lVal lua.LValue) map[string]string {

	//do we really have an LTable in the passed LValue?
	if lVal.Type() != lua.LTTable {
		return createErrorMap(scriptName,
			"script does not contain required return table variable 'rettbl'")
	}
	luaTable := lVal.(*lua.LTable)

	mp := make(map[string]string)

	luaTable.ForEach(func(k lua.LValue, v lua.LValue) {
		key := k.String()
		val := v.String()
		mp[key] = val
	})
	return mp
}

//helper to uniformly return errors during script execution
func createErrorMap(scriptName, details string) map[string]string {
	errMap := make(map[string]string)
	errMap["Script-Error"] = details
	return errMap
}
