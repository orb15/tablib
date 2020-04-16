package tablib

/*
These tests focus on the script executable portions of the repo like
rolling on and picking from tables. Because of my lack of desire to write mocks
(and ensure I am using interfaces everywhere I could be to make this doable),
some of these tests are more like integration tests than unit tests.
*/

import (
	"strings"
	"testing"
)

func TestExecute_shouldErrorOnBadScriptName(t *testing.T) {
	lua := `
  print("hello world")
  `
	repo := NewTableRepository()
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test1", DefaultParamSpecificationCallback)

	if len(mp) != 1 {
		t.Error("Invalid script results")
	}
	_, found := mp["Script-Error"]
	if !found {
		t.Error("Missing expected error message")
	}
}

func TestExecute_shouldExecuteBasicRollNoParams(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	lua := `
  local t = require("tables")
  results = {}
  function main(goData)
  val = t.roll("TestTable_Flat")
  results["iRolled"] = val
  end
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test", DefaultParamSpecificationCallback)

	if len(mp) != 1 {
		t.Error("Invalid script results")
	}
	val, found := mp["iRolled"]
	if !found {
		t.Error("Missing expected value")
	}
	if val != "item 1" {
		t.Error("Unexpected or missing table value")
	}
}

func TestExecute_shouldExecuteSeveralRollsNoParams(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	lua := `
  local t = require("tables")
  results = {}
  function main(goData)
  val1 = t.roll("TestTable_Flat")
  val2 = t.roll("TestTable_Flat")
  val3 = t.roll("TestTable_Flat")
  results["val1"] = val1
  results["val2"] = val2
  results["val3"] = val3
  end
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test", DefaultParamSpecificationCallback)

	if len(mp) != 3 {
		t.Error("Invalid script results")
	}
	_, found := mp["val1"]
	if !found {
		t.Error("Missing expected val1")
	}
	_, found = mp["val2"]
	if !found {
		t.Error("Missing expected val2")
	}
	_, found = mp["val3"]
	if !found {
		t.Error("Missing expected val3")
	}
}

func TestExecute_shouldFailIfRollCalledBadTable(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	lua := `
  local t = require("tables")
  results = {}
  function main(goData)
  val1 = t.roll("foo")
  results["val1"] = val1
  end
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test", DefaultParamSpecificationCallback)

	if len(mp) != 1 {
		t.Error("Invalid script results")
	}
	val, found := mp["val1"]
	if !found {
		t.Error("Missing expected error message")
	}
	if !strings.HasPrefix(val, "ERROR: The roll failed. Does the table:") {
		t.Error("Did not receive expected script error")
	}
}

func TestExecute_shouldFailIfRollCalledNoTable(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	lua := `
  local t = require("tables")
  results = {}
  function main(goData)
  val1 = t.roll()
  results["val1"] = val1
  end
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test", DefaultParamSpecificationCallback)

	if len(mp) != 1 {
		t.Error("Invalid script results")
	}
	val, found := mp["val1"]
	if !found {
		t.Error("Missing expected error message")
	}
	if val != "ERROR: roll(tableName) requires 1 argument, received: 0" {
		t.Error("Did not receive expected script error")
	}
}

func TestExecute_shouldFailIfRollCalledTooManyArgs(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	lua := `
  local t = require("tables")
  results = {}
  function main(goData)
  val1 = t.roll("TestTable_Flat", 2)
  results["val1"] = val1
  end
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test", DefaultParamSpecificationCallback)

	if len(mp) != 1 {
		t.Error("Invalid script results")
	}
	val, found := mp["val1"]
	if !found {
		t.Error("Missing expected error message")
	}
	if val != "ERROR: roll(tableName) requires 1 argument, received: 2" {
		t.Error("Did not receive expected script error")
	}
}

func TestExecute_shouldFailIfRollCalledWrongArgType(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	lua := `
  local t = require("tables")
  results = {}
  function main(goData)
  val1 = t.roll(2)
  results["val1"] = val1
  end
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test", DefaultParamSpecificationCallback)

	if len(mp) != 1 {
		t.Error("Invalid script results")
	}
	val, found := mp["val1"]
	if !found {
		t.Error("Missing expected error message")
	}
	if !strings.HasPrefix(val, "ERROR: roll(tableName) requires string argument") {
		t.Error("Did not receive expected script error")
	}
}

func TestExecute_shouldUsePassedParams(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	lua := `
  params = {}
  params["p1"] = "opt1-1|opt1-2"
  params["p2"] = "opt2-1|opt2-2"

  results = {}
  function main(goData)
  results["p1"] = goData["p1"]
  results["p2"] = goData["p2"]
  end
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test", func([]*ParamSpecification) map[string]string {
		retmap := make(map[string]string)
		retmap["p1"] = "opt1-2"
		retmap["p2"] = "opt2-2"
		return retmap
	})

	if len(mp) != 2 {
		t.Error("Invalid script results")
	}
	p1, found := mp["p1"]
	if !found {
		t.Error("Missing p1 param value")
	}
	if p1 != "opt1-2" {
		t.Error("p1 parameter value not transmitted properly")
	}
	p2, found := mp["p2"]
	if !found {
		t.Error("Missing p2 param value")
	}
	if p2 != "opt2-2" {
		t.Error("p2 parameter value not transmitted properly")
	}
}

func TestExecute_shouldUseDefaultParamsOnNoCallbackProvided(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	lua := `
  params = {}
  params["p1"] = "opt1-1|opt1-2"
  params["p2"] = "opt2-1|opt2-2"

  results = {}
  function main(goData)
  results["p1"] = goData["p1"]
  results["p2"] = goData["p2"]
  end
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test", nil)

	if len(mp) != 2 {
		t.Error("Invalid script results")
	}
	p1, found := mp["p1"]
	if !found {
		t.Error("Missing p1 param value")
	}
	if p1 != "opt1-1" {
		t.Error("p1 parameter value not transmitted properly")
	}
	p2, found := mp["p2"]
	if !found {
		t.Error("Missing p2 param value")
	}
	if p2 != "opt2-1" {
		t.Error("p2 parameter value not transmitted properly")
	}
}

func TestExecute_shouldHandleMissingMainNoParams(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	lua := `
  results = {}
  results["p1"] = "foo"
  results["p2"] = "bar"
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test", DefaultParamSpecificationCallback)
	if len(mp) != 1 {
		t.Error("Missing expected error message")
	}
	errMsg, found := mp["Script-Error"]
	if !found {
		t.Error("Missing error key value")
	}
	if !strings.HasPrefix(errMsg, "executing main():  attempt to call a non-function ") {
		t.Error("improper error message")
	}
}

func TestExecute_shouldHandleMissingMainWithParams(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	lua := `
  params = {}
  results = {}
  results["p1"] = "foo"
  results["p2"] = "bar"
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test", DefaultParamSpecificationCallback)
	if len(mp) != 1 {
		t.Error("Missing expected error message")
	}
	errMsg, found := mp["Script-Error"]
	if !found {
		t.Error("Missing error key value")
	}
	if !strings.HasPrefix(errMsg, "executing main():  attempt to call a non-function") {
		t.Error("improper error message")
	}
}

func TestExecute_shouldHandleMissingResultsTable(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	lua := `
  function main(goData)
  end
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test", DefaultParamSpecificationCallback)
	if len(mp) != 1 {
		t.Error("Missing expected error message")
	}
	errMsg, found := mp["Script-Error"]
	if !found {
		t.Error("Missing error key value")
	}
	if !strings.HasPrefix(errMsg, "missing the required execution results table:") {
		t.Error("improper error message")
	}
}

func TestExecute_shouldHandleMissingGoData(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	lua := `
  params = {}
  params["p1"] = "7even|foo|bar"
  results = {}
  function main()
  results["p1"] = "bar"
  end
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test", DefaultParamSpecificationCallback)
	if len(mp) != 1 {
		t.Error("Missing expected data")
	}
	p1, found := mp["p1"]
	if !found {
		t.Error("Missing returned map key")
	}
	if p1 != "bar" {
		t.Error("missing or bad returned data value")
	}
}

func TestExecute_shouldHandleBasicPick(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	lua := `
  local t = require("tables")
  results = {}
  function main()
  results["pick"] = t.pick("TestTable_Flat", 1)
  end
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test", DefaultParamSpecificationCallback)
	if len(mp) != 1 {
		t.Error("Missing expected data")
	}
	p1, found := mp["pick"]
	if !found {
		t.Error("Missing returned map key")
	}
	if p1 != "item 1" {
		t.Error("missing or bad returned data value")
	}
}

func TestExecute_shouldHandleBasicMultiPick(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1
    - item 2`

	lua := `
  local t = require("tables")
  results = {}
  function main()
  results["pick"] = t.pick("TestTable_Flat", 2)
  end
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test", DefaultParamSpecificationCallback)
	if len(mp) != 1 {
		t.Error("Missing expected data")
	}
	p1, found := mp["pick"]
	if !found {
		t.Error("Missing returned map key")
	}
	if p1 != "item 1|item 2" && p1 != "item 2|item 1" {
		t.Error("missing or bad returned data value")
	}
}

func TestExecute_shouldFailPickCalledTooManyArgs(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	lua := `
  local t = require("tables")
  results = {}
  function main()
  results["pick"] = t.pick("TestTable_Flat", 1, "|")
  end
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test", DefaultParamSpecificationCallback)
	if len(mp) != 1 {
		t.Error("Missing expected data")
	}
	p1, found := mp["pick"]
	if !found {
		t.Error("Missing returned map key")
	}
	if p1 != "ERROR: pick(tableName, count) requires 2 arguments received: 3" {
		t.Error("missing or bad returned data value")
	}
}

func TestExecute_shouldFailPickCalledWrongArgType1(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	lua := `
  local t = require("tables")
  results = {}
  function main()
  results["pick"] = t.pick(1, "TestTable_Flat")
  end
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test", DefaultParamSpecificationCallback)
	if len(mp) != 1 {
		t.Error("Missing expected data")
	}
	p1, found := mp["pick"]
	if !found {
		t.Error("Missing returned map key")
	}
	if !strings.HasPrefix(p1, "ERROR: pick(tableName, count), tablename must be a string, received type:") {
		t.Error("missing or bad returned data value")
	}
}

func TestExecute_shouldFailPickCalledWrongArgType2(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	lua := `
  local t = require("tables")
  results = {}
  function main()
  results["pick"] = t.pick("TestTable_Flat", "TestTable_Flat")
  end
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test", DefaultParamSpecificationCallback)
	if len(mp) != 1 {
		t.Error("Missing expected data")
	}
	p1, found := mp["pick"]
	if !found {
		t.Error("Missing returned map key")
	}
	if !strings.HasPrefix(p1, "ERROR: pick(tableName, count), count must be an integer, received type:") {
		t.Error("missing or bad returned data value")
	}
}

func TestExecute_shouldFailIfPickCalledBadTable(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	lua := `
  local t = require("tables")
  results = {}
  function main()
  results["pick"] = t.pick("foo", 1)
  end
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test", DefaultParamSpecificationCallback)
	if len(mp) != 1 {
		t.Error("Missing expected data")
	}
	p1, found := mp["pick"]
	if !found {
		t.Error("Missing returned map key")
	}
	if !strings.HasPrefix(p1, "ERROR: The pick failed. Does the table:") {
		t.Error("Missing expected error message")
	}
}

//the prupose of this test is to perform operations typically found in the
//expected lua code since the lua VM here is intended to have limited
//functionality as a general purpose software VM. Specifically, certain lua
//modeules from the core library are excluded  and this test is to assure
//lua can still fulfill its role in this tighted security environment. If this
//test fails, it means that the lua modules allowed in the VM are too restrictive
//or alternately, the thing attempted in lua is currently deemed too unsafe
func TestExecute_moduleCompatibilityTest(t *testing.T) {
	yml := `
  definition:
    name: Icecream_Flavors
    type: flat
    note: this is an optional note
  content:
    - chocolate
    - vanilla
    - strawberry`

	lua := `
  local t = require("tables")

	params = {}
  params["scoops"] = "2|1|3"
	params["favorite"] = "random|chocolate|vanilla|strawberry"

	results = {}
	local scoop1 = ""
	local scoop2 = ""
	local scoop3 = ""

  function main(goData)

		if goData["favorite"] == "random" then
		  scoop1 = t.roll("Icecream_Flavors")
		  if goData["scoops"] == "2" then
				scoop2 = t.roll("Icecream_Flavors")
			end
			if goData["scoops"] == "3" then
				scoop2 = t.roll("Icecream_Flavors")
				scoop3 = t.roll("Icecream_Flavors")
			end
    else
		  scoop1 = goData["favorite"]
			if goData["scoops"] == "2" then
				scoop2 = goData["favorite"]
			end
			if goData["scoops"] == "3" then
				scoop2 = goData["favorite"]
				scoop3 = goData["favorite"]
			end
		end

		results["scoop1"] = scoop1
		results["scoop2"] = scoop2
		results["scoop3"] = scoop3
  end
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.AddLuaScript("test", lua, nil)
	mp := repo.Execute("test", DefaultParamSpecificationCallback)
	t.Log(mp)

	if len(mp) != 3 {
		t.Error("Invalid script results")
	}
	val1, found := mp["scoop1"]
	if !found {
		t.Error("Missing expected val1")
	}
	if val1 != "chocolate" && val1 != "vanilla" && val1 != "strawberry" {
		t.Error("Val1 bad data")
	}
	val2, found := mp["scoop2"]
	if !found {
		t.Error("Missing expected val1")
	}
	if val2 != "chocolate" && val2 != "vanilla" && val2 != "strawberry" && val2 != "" {
		t.Error("Val2 bad data")
	}
	val3, found := mp["scoop3"]
	if !found {
		t.Error("Missing expected val1")
	}
	if val3 != "chocolate" && val3 != "vanilla" && val3 != "strawberry" && val3 != "" {
		t.Error("Val3 bad data")
	}
}
