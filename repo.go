package tablib

import (
	"fmt"
	"strings"
	"sync"
	"tablib/table"
	"tablib/tableresult"
	"tablib/util"
	"tablib/validate"

	lua "github.com/yuin/gopher-lua"
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

	//not locking repo here so parse + validate can be multithreaded if caller desires

	var table *table.Table

	//is this even valid YAML?
	err := yaml.Unmarshal(yamlBytes, &table)
	if err != nil {
		return nil, err
	}

	//validate the table and parse portions of it since we are tearing the table
	//apart to do the validation anyway
	validationResults := table.Validate()

	//by definition, tables that arrive here are not inline tables
	table.IsInlineTable = false

	//lock the repo now since we will write to it
	cr.lock.Lock()
	defer cr.lock.Unlock()

	//store the table in the repo
	fullName := util.BuildFullName(table.Definition.Name, "")
	cr.tableStore[fullName] = &tableData{
		yamlSource:  string(yamlBytes),
		parsedTable: table,
	}

	//if the table has any inline tables, add these as well - inline
	//tables will be first-class 'flat' tables now
	if len(table.Inline) > 0 {
		inlines := extractInlineTables(table)
		for _, ilt := range inlines {
			cr.tableStore[ilt.Definition.Name] = &tableData{
				yamlSource:  "",
				parsedTable: ilt,
			}
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
	defer cr.lock.Unlock()
	return "", nil
}
func (cr *concreteTableRepo) Search(namePredicate string, tags []string) []*SearchResult {
	cr.lock.RLock()
	defer cr.lock.Unlock()
	return make([]*SearchResult, 0)
}
func (cr *concreteTableRepo) Roll(tableName string, execsDesired int) *tableresult.TableResult {
	cr.lock.RLock()
	defer cr.lock.Unlock()

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
		operation: "roll",
		count:     execsDesired,
	}
	exeng := newExecutionEngine()
	exeng.execute(wp, tr)
	return tr
}
func (cr *concreteTableRepo) Pick(tableName string, count int) *tableresult.TableResult {
	cr.lock.RLock()
	defer cr.lock.Unlock()
	return nil
}
func (cr *concreteTableRepo) Execute(scriptName string) (map[string]string, error) {
	cr.lock.RLock()
	defer cr.lock.Unlock()
	return make(map[string]string), nil
}

//for each inline table in a table, create a full-featured table
func extractInlineTables(mainTable *table.Table) []*table.Table {
	inlinesAsTables := make([]*table.Table, len(mainTable.Inline))
	for _, ilt := range mainTable.Inline {

		def := &table.DefinitionPart{
			Name:      ilt.FullyQualifiedName,
			TableType: "flat",
		}

		content := make([]string, len(ilt.Content))
		for _, c := range ilt.Content {
			content = append(content, c)
		}

		tbl := &table.Table{
			Definition:    def,
			RawContent:    content,
			IsValid:       true,
			IsInlineTable: true,
		}
		inlinesAsTables = append(inlinesAsTables, tbl)
	}
	return inlinesAsTables
}
