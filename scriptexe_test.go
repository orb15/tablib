package tablib

/*
These tests focus on the script executable portions of the repo like
rolling on and picking from tables. Because of my lack of desire to write mocks
(and ensure I am using interfaces everywhere I could be to make this doable),
some of these tests are more like integration tests than unit tests.
*/

import (
	"strconv"
	"strings"
	"testing"
)

func TestExecute_shouldErrorOnBadScriptName(t *testing.T) {
	lua := `
  print("hello world")
  `
	repo := NewTableRepository()
	repo.AddLuaScript("test", lua)
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
	repo.AddLuaScript("test", lua)
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
	repo.AddLuaScript("test", lua)
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
	repo.AddLuaScript("test", lua)
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
	repo.AddLuaScript("test", lua)
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
	repo.AddLuaScript("test", lua)
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
	repo.AddLuaScript("test", lua)
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
	repo.AddLuaScript("test", lua)
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
	repo.AddLuaScript("test", lua)
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
	repo.AddLuaScript("test", lua)
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
	repo.AddLuaScript("test", lua)
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
	repo.AddLuaScript("test", lua)
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
	repo.AddLuaScript("test", lua)
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
	repo.AddLuaScript("test", lua)
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
	repo.AddLuaScript("test", lua)
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
	repo.AddLuaScript("test", lua)
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
	repo.AddLuaScript("test", lua)
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
	repo.AddLuaScript("test", lua)
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
	repo.AddLuaScript("test", lua)
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

func TestExecute_shouldExecuteValidDiceExpression(t *testing.T) {
	lua1 := `
  local t = require("tables")
  results = {}
  function main(goData)
  val = t.dice("1d1")
  results["dice"] = val
  end
  `

	lua2 := `
  local t = require("tables")
  results = {}
  function main(goData)
  val = t.dice("5d6 * 0")
  results["dice"] = val
  end
  `

	repo := NewTableRepository()
	repo.AddLuaScript("test1", lua1)
	repo.AddLuaScript("test2", lua2)
	mp1 := repo.Execute("test1", nil)
	mp2 := repo.Execute("test2", nil)

	if len(mp1) != 1 {
		t.Error("Invalid script results lua1")
	}
	val1, found := mp1["dice"]
	if !found {
		t.Error("Missing expected value lua1")
	}
	if val1 != "1" {
		t.Error("Unexpected or missing table value lua1")
	}

	if len(mp2) != 1 {
		t.Error("Invalid script results lua2")
	}
	val2, found := mp2["dice"]
	if !found {
		t.Error("Missing expected value lua2")
	}
	if val2 != "0" {
		t.Error("Unexpected or missing table value lua2")
	}
}

func TestExecute_shouldReturnErrorValOnInvalidDiceExpression(t *testing.T) {
	lua1 := `
  local t = require("tables")
  results = {}
  function main(goData)
  val = t.dice("1d")
  results["dice"] = val
  end
  `

	repo := NewTableRepository()
	repo.AddLuaScript("test1", lua1)
	mp1 := repo.Execute("test1", nil)

	if len(mp1) != 1 {
		t.Error("Invalid script results lua1")
	}
	val1, found := mp1["dice"]
	if !found {
		t.Error("Missing expected value lua1")
	}
	if val1 != strconv.Itoa(badDiceRollInteger) {
		t.Error("Unexpected or missing table value lua1")
	}
}

func TestExecute_shouldFailOnDiceExpressionWrongArgType(t *testing.T) {
	lua1 := `
  local t = require("tables")
  results = {}
  function main(goData)
  val = t.dice(6)
  results["dice"] = val
  end
  `

	repo := NewTableRepository()
	repo.AddLuaScript("test1", lua1)
	mp1 := repo.Execute("test1", nil)

	if len(mp1) != 1 {
		t.Error("Invalid script results lua1")
	}
	val1, found := mp1["dice"]
	if !found {
		t.Error("Missing expected value lua1")
	}
	if val1 != strconv.Itoa(badDiceRollInteger) {
		t.Error("Unexpected or missing table value lua1")
	}
}

func TestExecute_shouldFailOnDiceExpressionWrongArgCount(t *testing.T) {
	lua1 := `
  local t = require("tables")
  results = {}
  function main(goData)
  val = t.dice("1d6", 6)
  results["dice"] = val
  end
  `

	repo := NewTableRepository()
	repo.AddLuaScript("test1", lua1)
	mp1 := repo.Execute("test1", nil)

	if len(mp1) != 1 {
		t.Error("Invalid script results lua1")
	}
	val1, found := mp1["dice"]
	if !found {
		t.Error("Missing expected value lua1")
	}
	if val1 != strconv.Itoa(badDiceRollInteger) {
		t.Error("Unexpected or missing table value lua1")
	}
}

func TestExecute_shouldConcatenateProperly(t *testing.T) {
	lua1 := `
  local t = require("tables")
  results = {}
	local plsConcat = {}
	plsConcat[1] = "Col. Mustard"
	plsConcat[2] = " "
	plsConcat[3] = "did it in the "
	plsConcat[4] = "Library with the"
	plsConcat[5] = " Candlestick"
	plsConcat[6] = "!\n"
	plsConcat[7] = "Now is the time for all good men to come to the aid of their country"
	function main(goData)
		results["concat"] = t.concat(plsConcat)
  end
  `

	repo := NewTableRepository()
	repo.AddLuaScript("test1", lua1)
	mp1 := repo.Execute("test1", nil)

	if len(mp1) != 1 {
		t.Error("Invalid script results lua1")
	}
	val, found := mp1["concat"]
	if !found {
		t.Error("Missing expected value")
	}
	if val != "Col. Mustard did it in the Library with the Candlestick!\nNow is the time for all good men to come to the aid of their country" {
		t.Errorf("String not properly concatenated: %s", val)
	}
}

func TestExecute_shouldConcatenateWithLuaNumberTypes(t *testing.T) {
	lua1 := `
  local t = require("tables")
  results = {}
	local plsConcat = {}
	plsConcat[1] = "This is the number 1: "
	plsConcat[2] = t.dice("1d1")
	function main(goData)
		results["concat"] = t.concat(plsConcat)
  end
  `

	repo := NewTableRepository()
	repo.AddLuaScript("test1", lua1)
	mp1 := repo.Execute("test1", nil)

	if len(mp1) != 1 {
		t.Error("Invalid script results lua1")
	}
	val, found := mp1["concat"]
	if !found {
		t.Error("Missing expected value")
	}
	if val != "This is the number 1: 1" {
		t.Errorf("String not properly concatenated: %s", val)
	}
}

func TestExecute_shouldFailToConcatenateWithBadArgCount(t *testing.T) {
	lua1 := `
  local t = require("tables")
  results = {}
	local plsConcat = {}
	plsConcat[1] = "This is the number 1: "
	plsConcat[2] = t.dice("1d1")
	function main(goData)
		results["concat"] = t.concat(plsConcat, 7)
  end
  `

	repo := NewTableRepository()
	repo.AddLuaScript("test1", lua1)
	mp1 := repo.Execute("test1", nil)

	if len(mp1) != 1 {
		t.Error("Invalid script results lua1")
	}
	val, found := mp1["concat"]
	if !found {
		t.Error("Missing expected value")
	}
	if !strings.HasPrefix(val, "ERROR: ") {
		t.Errorf("Did not receive expected error msg: %s", val)
	}
}

func TestExecute_shouldFailToConcatenateWithBadArgType(t *testing.T) {
	lua1 := `
  local t = require("tables")
  results = {}
	plsConcat = "This should fail"
	function main(goData)
		results["concat"] = t.concat(plsConcat)
  end
  `

	repo := NewTableRepository()
	repo.AddLuaScript("test1", lua1)
	mp1 := repo.Execute("test1", nil)

	if len(mp1) != 1 {
		t.Error("Invalid script results lua1")
	}
	val, found := mp1["concat"]
	if !found {
		t.Error("Missing expected value")
	}
	if !strings.HasPrefix(val, "ERROR: ") {
		t.Errorf("Did not receive expected error msg: %s", val)
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

	local timesToRoll = t.dice("1d3")
	results["timesRolled"] = timesToRoll

  function main(goData)

		if goData["favorite"] == "random" then
		  scoop1 = t.roll("Icecream_Flavors")
		  if goData["scoops"] == "2" then
				scoop2 = t.roll("Icecream_Flavors")
			else
				scoop2 = "Error: Should have only 2 scoops present"
			end
		else
			scoop1 = "Error: favorite should always be 'random'"
		end

		results["scoop1"] = scoop1
		results["scoop2"] = scoop2

		for i=1,timesToRoll
		do
			results[tostring(i)] = t.roll("Icecream_Flavors")
		end
  end
  `

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	err := repo.AddLuaScript("test", lua)
	if err != nil {
		t.Error(err)
	}
	mp := repo.Execute("test", DefaultParamSpecificationCallback)
	t.Log(mp)

	if len(mp) < 4 || len(mp) > 6 { //scoop1 and scoop2 values, timesRolled, plus 1-3 others
		t.Error("Invalid script results")
	}
	val1, found := mp["scoop1"]
	if !found {
		t.Error("Missing expected scoop1 from map")
	}
	if val1 != "chocolate" && val1 != "vanilla" && val1 != "strawberry" {
		t.Error("scoop1 bad data")
	}
	val2, found := mp["scoop2"]
	if !found {
		t.Error("Missing expected scoop2 from map")
	}
	if val2 != "chocolate" && val2 != "vanilla" && val2 != "strawberry" && val2 != "" {
		t.Error("scoop2 bad data")
	}

	//variable part of test
	timesString, found := mp["timesRolled"]
	if !found {
		t.Error("Missing timesRolled from results map")
	}
	timesRolled, err := strconv.Atoi(timesString)
	if err != nil {
		t.Error(err)
	}
	if timesRolled < 1 || timesRolled > 3 {
		t.Errorf("TimesRolled outside of bounds: %d", timesRolled)
	}
	if len(mp) != 3+timesRolled {
		t.Error("Wrong length for return results")
	}
	for i := 1; i <= timesRolled; i++ {
		val, found := mp[strconv.Itoa(i)]
		if !found {
			t.Errorf("Missing %dth value from variable part of map", i)
		}
		if val != "chocolate" && val != "vanilla" && val != "strawberry" {
			t.Errorf("the %dth value has bad data: %s", i, val)
		}
	}

}
