package tablib

/*
These tests focus on the non-executable portions of the repo - adding tables
and scripts, searching and listing etc. Tests focused on table or script
execution are located in other test files.
*/
import (
	"sync"
	"tablib/validate"
	"testing"
)

func TestAddTable_shouldAddValidFlatTableToRepo(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {#1}
  inline:
    - id: 1
      content:
        - Inline 1`

	//attempt to add valid table
	cr := newConcreteRepo()
	vr, err := cr.AddTable([]byte(yml))
	failOnErr("Bad YAML", err, t)
	failOnInvalid("Invalid table", vr, t)

	//ensure table is actually stored
	if len(cr.tableStore) != 2 { //2 here b/c of inline
		t.Error("Did not store valid table")
	}

	//fetch the table directly
	tblData, found := cr.tableStore["TestTable_Flat"]
	if !found {
		t.Error("Unable to find expected table in repo")
	}
	if tblData == nil {
		t.Error("Did not retrieve table by name")
	}

	//do we have the table YAML captured?
	if tblData.yamlSource != yml {
		t.Error("Repo does not contain table source")
	}

	//do we have a parsed version of the table?
	tbl := tblData.parsedTable

	if len(tbl.RawContent) != 1 {
		t.Errorf("Did not capture raw content properly")
	}

	//check implicit dice on flat table
	if len(tbl.Definition.DiceParsed) != 1 {
		t.Error("Missing or invalid parsed dice")
	}
	if tbl.Definition.DiceParsed[0].Count != 1 {
		t.Error("Failed to set implicit Count")
	}
	if tbl.Definition.DiceParsed[0].DieType != len(tbl.RawContent) {
		t.Error("Failed to set implicit Dice Type based on content length")
	}
	if tbl.Definition.DiceParsed[0].Operator != "none" {
		t.Error("Failed to set implicit operator")
	}

	//ensure that the inline table has been properly set up as a first-class table
	tblData, found = cr.tableStore["TestTable_Flat.1"]
	if !found {
		t.Error("Unable to find inline table in repo")
	}
	if tblData == nil {
		t.Error("Did not retrieve table by name")
	}

	tbl = tblData.parsedTable

	if len(tbl.RawContent) != 1 {
		t.Error("Did not capture raw content properly")
	}
	if !tbl.IsInlineTable {
		t.Error("Inline table not flagged properly")
	}
	if len(tbl.Definition.DiceParsed) != 1 {
		t.Error("Missing or invalid parsed dice for inline table")
	}
	if tbl.Definition.DiceParsed[0].Count != 1 {
		t.Error("Failed to set implicit Count for inline table")
	}
	if tbl.Definition.DiceParsed[0].DieType != len(tbl.RawContent) {
		t.Error("Failed to set implicit Dice Type based on content length for inline table")
	}
	if tbl.Definition.DiceParsed[0].Operator != "none" {
		t.Error("Failed to set implicit operator for inline table")
	}
}

func TestAddTable_shouldAddValidFlatRangeToRepo(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Range
    type: range
    roll: 1d4
  content:
    - "{1-2}item 1 - {#1}"
    - "{3-4}item 2"
  inline:
    - id: 1
      content:
        - Inline 1`

	//attempt to add valid table
	cr := newConcreteRepo()
	vr, err := cr.AddTable([]byte(yml))
	failOnErr("Bad YAML", err, t)
	failOnInvalid("Invalid table", vr, t)

	//ensure table is actually stored
	if len(cr.tableStore) != 2 { //2 here b/c of ioline
		t.Error("Did not store valid table")
	}

	//fetch the table directly
	tblData, found := cr.tableStore["TestTable_Range"]
	if !found {
		t.Error("Unable to find expected table in repo")
	}
	if tblData == nil {
		t.Error("Did not retrieve table by name")
	}

	//do we have the table YAML captured?
	if tblData.yamlSource != yml {
		t.Error("Repo does not contain table source")
	}

	//do we have a parsed version of the table?
	tbl := tblData.parsedTable

	if len(tbl.RangeContent) != 2 {
		t.Error("Did not capture raw content properly")
	}

	//check implicit dice on flat table
	if len(tbl.Definition.DiceParsed) != 1 {
		t.Error("Missing or invalid parsed dice")
	}
	if tbl.Definition.DiceParsed[0].Count != 1 {
		t.Error("Failed to set die Count")
	}
	if tbl.Definition.DiceParsed[0].DieType != 4 {
		t.Error("Failed to set Dice Type based on roll description")
	}
	if tbl.Definition.DiceParsed[0].Operator != "none" {
		t.Error("Failed to set die operator")
	}

	//ensure that the inline table has been properly set up as a first-class table
	tblData, found = cr.tableStore["TestTable_Range.1"]
	if !found {
		t.Error("Unable to find inline table in repo")
	}
	if tblData == nil {
		t.Error("Did not retrieve table by name")
	}

	tbl = tblData.parsedTable

	if len(tbl.RawContent) != 1 {
		t.Error("Did not capture raw content properly")
	}
	if !tbl.IsInlineTable {
		t.Error("Inline table not flagged properly")
	}
	if len(tbl.Definition.DiceParsed) != 1 {
		t.Error("Missing or invalid parsed dice for inline table")
	}
	if tbl.Definition.DiceParsed[0].Count != 1 {
		t.Error("Failed to set implicit Count for inline table")
	}
	if tbl.Definition.DiceParsed[0].DieType != len(tbl.RawContent) {
		t.Error("Failed to set implicit Dice Type based on content length for inline table")
	}
	if tbl.Definition.DiceParsed[0].Operator != "none" {
		t.Error("Failed to set implicit operator for inline table")
	}
}

func TestAddTable_shouldNotAddInvalidTableToRepo(t *testing.T) {
	yml := `
  definition:
    name:
    type: flat
    note: this is an optional note
  content:
    - item 1`

	cr := newConcreteRepo()
	cr.AddTable([]byte(yml))
	if len(cr.tableStore) != 0 {
		t.Error("Added invalid table to repo")
	}
}

func TestAddTable_shouldNotAddTablesTpoRepoOnInvalidInlineTable(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {#1}
  inline:
    - id: 1
      content:
        - "Inline 1{$$}"`

	cr := newConcreteRepo()
	cr.AddTable([]byte(yml))
	if len(cr.tableStore) != 0 {
		t.Error("Added invalid table to repo")
	}
}

func TestAddTable_shouldNotAddInvalidTableYAMLToRepo(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: "flat
    note: this is an optional note
  content:
    - item 1`

	cr := newConcreteRepo()
	_, err := cr.AddTable([]byte(yml))
	if err == nil {
		t.Error("Did not get expected YAML parser error")
	}
	if len(cr.tableStore) != 0 {
		t.Error("Added invalid table to repo")
	}
}

func TestAddLuaScript_shouldAddValidLuaToRepo(t *testing.T) {
	lua := `
  print("dlrow olleh")
  `

	cr := newConcreteRepo()
	err := cr.AddLuaScript("test", lua)
	failOnErr("Not not accept Lua script", err, t)

	//access the script directly
	if len(cr.scriptStore) != 1 {
		t.Error("Script is either not in the store or too many entries created")
	}
	scriptData, found := cr.scriptStore["test"]
	if !found {
		t.Error("Did not store lua script")
	}
	if scriptData == nil {
		t.Error("Did not fetch lua script")
	}
	if scriptData.scriptSource != lua {
		t.Error("Did not store raw script source as expected")
	}
	if scriptData.parsedScript == nil {
		t.Error("Did not precompile lua script")
	}
}

func TestAddLuaScript_shouldNotAddUnparsableLuaToRepo(t *testing.T) {
	lua := `
  print("dlrow olleh)
  `

	cr := newConcreteRepo()
	err := cr.AddLuaScript("test", lua)
	if err == nil {
		t.Error("Failed to error on unparsable lua")
	}
	if len(cr.scriptStore) != 0 {
		t.Error("Added bad lua script to store")
	}
}

func TestTableForName_shouldFetchTable(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {#1}
  inline:
    - id: 1
      content:
        - Inline 1`

	//attempt to add valid table
	cr := newConcreteRepo()
	vr, err := cr.AddTable([]byte(yml))
	failOnErr("Bad YAML", err, t)
	failOnInvalid("Invalid table", vr, t)

	tbl, err := cr.tableForName("TestTable_Flat")
	failOnErr("Did not fetch table as expected", err, t)
	if tbl == nil {
		t.Error("Did not fetch table as expected")
	}
}

func TestTableForName_shouldNotFetchMissingTable(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {#1}
  inline:
    - id: 1
      content:
        - Inline 1`

	//attempt to add valid table
	cr := newConcreteRepo()
	vr, err := cr.AddTable([]byte(yml))
	failOnErr("Bad YAML", err, t)
	failOnInvalid("Invalid table", vr, t)

	tbl, err := cr.tableForName("TestTable_Foo")
	if err == nil {
		t.Error("Did not error on bad table name")
	}
	if tbl != nil {
		t.Error("Fetched a non-existant table!?")
	}
}

func TestScriptForName_shouldFetchScript(t *testing.T) {
	lua := `
  print("dlrow olleh")
  `
	cr := newConcreteRepo()
	err := cr.AddLuaScript("test", lua)
	failOnErr("Not not accept Lua script", err, t)

	script, err := cr.scriptForName("test")
	failOnErr("Did not fetch script as expected", err, t)
	if script == nil {
		t.Error("Did not fetch script as expected")
	}
}

func TestScriptForName_shouldNotFetchMissingScript(t *testing.T) {
	lua := `
  print("dlrow olleh")
  `
	cr := newConcreteRepo()
	err := cr.AddLuaScript("test", lua)
	failOnErr("Not not accept Lua script", err, t)

	script, err := cr.scriptForName("tset")
	if err == nil {
		t.Error("Did not error on missing script")
	}
	if script != nil {
		t.Error("Fetched a non-existant script!?S")
	}
}

/* ***********************************************
* Test Helpers
* ***********************************************/

func failOnErr(msg string, err error, t *testing.T) {
	if err != nil {
		t.Errorf("Unexpected err: %s: %s", msg, err)
	}
}

func failOnInvalid(msg string, vr *validate.ValidationResult, t *testing.T) {
	if !vr.Valid() {
		t.Errorf("%s: %s", msg, vr.Errors[0])
	}
}

func newConcreteRepo() *concreteTableRepo {
	return &concreteTableRepo{
		tableStore:  make(map[string]*tableData),
		scriptStore: make(map[string]*scriptData),
		lock:        &sync.RWMutex{},
	}
}