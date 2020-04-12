package tablib

import (
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

	lua := `local t = require("tables")
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

func TestExecute_shouldFailIfRollCalledBadArgs1(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	lua := `local t = require("tables")
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

func TestExecute_shouldFailIfRollCalledBadArgs2(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	lua := `local t = require("tables")
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
