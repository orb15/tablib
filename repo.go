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

	yaml "gopkg.in/yaml.v3"
)

type tableData struct {
	yamlSource  string
	parsedTable *table.Table
	tags        []string
}

type scriptData struct {
	scriptSource string
	parsedScript *lua.FunctionProto
	tags         []string
}

type concreteTableRepo struct {
	tableStore      map[string]*tableData
	scriptStore     map[string]*scriptData
	tagSearchCache  map[string][]*SearchResult
	nameSearchCache []*SearchResult
	lock            *sync.RWMutex
}

type nameResolver interface {
	tableForName(name string) (*table.Table, error)
	scriptForName(name string) (*lua.FunctionProto, error)
}

const (
	itemTypeTable  = "table"
	itemTypeScript = "script"
)

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

	//update caches - need to do this before table is stored so we
	//can compare how the table has changed if it is being updated

	fullName := util.BuildFullName(tbl.Definition.Name, "")
	//update the tag cache with the new table info
	cr.updateTagCache(fullName, itemTypeTable, tbl.Definition.Tags)

	//update the name cache
	cr.updateNameCache(fullName, itemTypeTable, tbl.Definition.Tags)

	//put the valid table in the repo
	cr.tableStore[fullName] = &tableData{
		yamlSource:  string(yamlBytes),
		parsedTable: tbl,
		tags:        tbl.Definition.Tags,
	}

	//store the inline tables for this table as first-class tables
	//note that inline tables are not returned in searches and cant be Listed
	//seperately from their mater table
	for _, ilt := range inlines {
		cr.tableStore[ilt.Definition.Name] = &tableData{
			yamlSource:  "",
			parsedTable: ilt,
			tags:        nil,
		}
	}

	return validationResults, nil
}

func (cr *concreteTableRepo) AddLuaScript(scriptName, luaScript string, tags []string) error {

	//not locking repo here so compilation can be multithreaded if caller desires

	//read and compile the lua script
	reader := strings.NewReader(luaScript)
	astStatements, err := parse.Parse(reader, luaScript)
	if err != nil {
		return err
	}

	//from a read of the source this is pretty unlikely but catching it anyway
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
		tags:         tags,
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

	if execsDesired <= 0 {
		tr.AddLog(fmt.Sprintf("Attempt to roll 0 or fewer times on table: %s", tableName))
	}

	tbl, found := cr.tableStore[tableName]
	if !found {
		tr.AddLog(fmt.Sprintf("Table: %s does not exist", tableName))
		return tr
	}

	wp := &workPackage{
		nameSvc:   cr,
		table:     tbl.parsedTable,
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
		nameSvc:   cr,
		table:     tbl.parsedTable,
		operation: table.OpPick,
		count:     1,
		pickCount: count,
	}
	exeng := newExecutionEngine()
	exeng.execute(wp, tr)
	return tr
}

func (cr *concreteTableRepo) Execute(scriptName string,
	callback ParamSpecificationRequestCallback) map[string]string {
	cr.lock.RLock()
	defer cr.lock.RUnlock()

	return executeScript(scriptName, cr, cr, callback)
}

//tableForName returns the underlying table for give table name
func (cr *concreteTableRepo) tableForName(name string) (*table.Table, error) {
	cr.lock.RLock()
	defer cr.lock.RUnlock()

	tblData, found := cr.tableStore[name]
	if !found {
		return nil, fmt.Errorf("Table does not exist: %s", name)
	}
	return tblData.parsedTable, nil
}

//scriptForName returns the underlying compiled script for give table name
func (cr *concreteTableRepo) scriptForName(name string) (*lua.FunctionProto, error) {
	cr.lock.RLock()
	defer cr.lock.RUnlock()

	scriptData, found := cr.scriptStore[name]
	if !found {
		return nil, fmt.Errorf("Script does not exist: %s", name)
	}
	return scriptData.parsedScript, nil
}

func (cr *concreteTableRepo) updateTagCache(fullName string, itemType string, tags []string) {
	//first, see if this object exists - if it does, note the previously cached tags
	var prevCachedTags []string
	switch itemType {
	case itemTypeTable:
		tableData, found := cr.tableStore[fullName]
		if found {
			prevCachedTags = tableData.tags
		}
	case itemTypeScript:
		scriptData, found := cr.scriptStore[fullName]
		if found {
			prevCachedTags = scriptData.tags
		}
	}

	//case 1: nothing to add to the tag cache
	if len(prevCachedTags) == 0 && len(tags) == 0 {
		return
	}

	sr := &SearchResult{
		Name: fullName,
		Type: itemType,
		Tags: tags,
	}

	//case 2: tags have been added but there were none before
	if len(prevCachedTags) == 0 && len(tags) > 0 {
		cr.addToTagCache(sr)
		return
	}

	//case 3: we need to remove previously cached tags
	if len(prevCachedTags) > 0 && len(tags) == 0 {
		cr.removeFromTagCache(sr, prevCachedTags)
		return
	}

	//case 4: we need to remove old tags and add new ones
	cr.removeFromTagCache(sr, prevCachedTags)
	cr.addToTagCache(sr)
}

func (cr *concreteTableRepo) removeFromTagCache(sr *SearchResult, oldTags []string) {
	for _, oldTag := range oldTags {
		cachedItems, found := cr.tagSearchCache[oldTag]
		if found { //checking this in case tag no longer exists
			idx := -1
			for _, item := range cachedItems {
				idx++
				if item.Name == sr.Name && item.Type == sr.Type {
					break
				}
			}
			//drop the matched item if one exists
			if idx > -1 && idx <= len(cachedItems)-1 { //make sure idx is valid

				//if this is the last cached item for this tag, delete the tag from the map
				if len(cachedItems) == 1 {
					delete(cr.tagSearchCache, oldTag)
				} else { //pull only the element out we no longer want (this screws with list ordering!)
					cachedItems[idx] = cachedItems[len(cachedItems)-1]
					cachedItems = cachedItems[:len(cachedItems)-1]
					cr.tagSearchCache[oldTag] = cachedItems
				}
			}
		}
	}
}

func (cr *concreteTableRepo) addToTagCache(sr *SearchResult) {
	for _, tag := range sr.Tags {
		list, found := cr.tagSearchCache[tag]
		if !found {
			list = make([]*SearchResult, 0, 1)
		}
		cr.tagSearchCache[tag] = append(list, sr)
	}
}

func (cr *concreteTableRepo) updateNameCache(fullName string, itemType string, tags []string) {

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
		Count:    1,
		DieType:  len(tbl.RawContent),
		Operator: "none",
	}
	dp := make([]*dice.ParseResult, 1, 1)
	dp[0] = dpr
	tbl.Definition.DiceParsed = dp
}
