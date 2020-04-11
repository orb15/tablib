package tablib

/*
These tests focus on the table executable portions of the repo like
rolling on and picking from tables. Because of my lack of desire to write mocks
(and ensure I am using interfaces everywhere I could be to make this doable),
some of these tests are more like integration tests than unit tests.
*/

import (
	"strings"
	"testing"

	"tablib/dice"
	"tablib/validate"
)

const (
	diceCycleCount = 100 //number of times to run the dice roll test
)

func TestRoll_shouldRollAsExpectedFlat(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	tr := repo.Roll("TestTable_Flat", 1)
	if len(tr.Result) != 1 {
		t.Error("Unexpected result count")
	}
	if tr.Result[0] != "item 1" {
		t.Error("Wrong result from table")
	}
	if len(tr.Log) != 2 {
		t.Error("Wrong Log info captured")
	}
}

func TestRoll_shouldMultiRollAsExpectedFlat(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	tr := repo.Roll("TestTable_Flat", 3)
	if len(tr.Result) != 3 {
		t.Error("Unexpected result count")
	}
	if tr.Result[0] != "item 1" {
		t.Error("Wrong result from table")
	}
	if tr.Result[1] != "item 1" {
		t.Error("Wrong result from table")
	}
	if tr.Result[2] != "item 1" {
		t.Error("Wrong result from table")
	}
	if len(tr.Log) != 6 {
		t.Error("Wrong Log info captured")
	}
}

func TestRoll_shouldRollAsExpectedRange(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: range
    roll: 1d4
    note: this is an optional note
  content:
    - "{1-4}item 1"`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	tr := repo.Roll("TestTable_Flat", 1)
	if len(tr.Result) != 1 {
		t.Error("Unexpected result count")
	}
	if tr.Result[0] != "item 1" {
		t.Error("Wrong result from table")
	}
	if len(tr.Log) != 2 {
		t.Error("Wrong Log info captured")
	}
}

func TestRoll_shouldMultiRollAsExpectedRange(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: range
    roll: 1d4
    note: this is an optional note
  content:
    - "{1-4}item 1"`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	tr := repo.Roll("TestTable_Flat", 3)
	if len(tr.Result) != 3 {
		t.Error("Unexpected result count")
	}
	if tr.Result[0] != "item 1" {
		t.Error("Wrong result from table")
	}
	if tr.Result[1] != "item 1" {
		t.Error("Wrong result from table")
	}
	if tr.Result[2] != "item 1" {
		t.Error("Wrong result from table")
	}
	if len(tr.Log) != 6 {
		t.Error("Wrong Log info captured")
	}
}

func TestRoll_shouldRollAsExpectedInlineRef(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - "item 1 - {#1}"
  inline:
    - id: 1
      content:
        - inline 1`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	tr := repo.Roll("TestTable_Flat", 1)
	if len(tr.Result) != 1 {
		t.Error("Unexpected result count")
	}
	if tr.Result[0] != "item 1 - inline 1" {
		t.Error("Wrong result from table")
	}
	if len(tr.Log) != 4 {
		t.Error("Wrong Log info captured")
	}
}

func TestRoll_shouldRollAsExpectedFlatToRangeRef(t *testing.T) {
	ymlf := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - "flat item 1 - {@TestTable_Range}"`

	ymlr := `
  definition:
    name: TestTable_Range
    type: range
    roll: 1d4
    note: this is an optional note
  content:
    - "{1-4}range item 1"`

	repo := NewTableRepository()
	repo.AddTable([]byte(ymlf))
	repo.AddTable([]byte(ymlr))
	tr := repo.Roll("TestTable_Flat", 1)
	if len(tr.Result) != 1 {
		t.Error("Unexpected result count")
	}
	if tr.Result[0] != "flat item 1 - range item 1" {
		t.Error("Wrong result from table")
	}
	if len(tr.Log) != 4 {
		t.Error("Wrong Log info captured")
	}
}

func TestRoll_shouldRollAsExpectedRangeToFlatRef(t *testing.T) {
	ymlf := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - "flat item 1"`

	ymlr := `
  definition:
    name: TestTable_Range
    type: range
    roll: 1d4
    note: this is an optional note
  content:
    - "{1-4}range item 1 - {@TestTable_Flat}"`

	repo := NewTableRepository()
	repo.AddTable([]byte(ymlf))
	repo.AddTable([]byte(ymlr))
	tr := repo.Roll("TestTable_Range", 1)
	if len(tr.Result) != 1 {
		t.Error("Unexpected result count")
	}
	if tr.Result[0] != "range item 1 - flat item 1" {
		t.Error("Wrong result from table")
	}
	if len(tr.Log) != 4 {
		t.Error("Wrong Log info captured")
	}
}

func TestRoll_shouldFailtoExpandBadTableRef(t *testing.T) {
	ymlf := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - "flat item 1"`

	ymlr := `
  definition:
    name: TestTable_Range
    type: range
    roll: 1d4
    note: this is an optional note
  content:
    - "{1-4}range item 1 - {@TestTable_Flat1}"`

	repo := NewTableRepository()
	repo.AddTable([]byte(ymlf))
	repo.AddTable([]byte(ymlr))
	tr := repo.Roll("TestTable_Range", 1)
	if len(tr.Result) != 1 {
		t.Error("Unexpected result count")
	}
	if tr.Result[0] != "range item 1 -  --BADREF: {@TestTable_Flat1}--" {
		t.Error("Wrong result from table")
	}
	if len(tr.Log) != 3 {
		t.Error("Wrong Log info captured")
	}
}

func TestRoll_shouldFailOnMissingTable(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	tr := repo.Roll("TestTable_Flat1", 1)
	if len(tr.Result) != 0 {
		t.Error("Unexpected result count")
	}
	if len(tr.Log) != 1 {
		t.Error("Wrong Log info captured")
	}
}

//if this test hangs, the depth counter code is borked
func TestRoll_shouldPreventInfiniteSelfRecursion(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - "item 1: {@TestTable_Flat}"`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.Roll("TestTable_Flat", 1)
	t.Log("Self recursion test PASS")
}

//if this test hangs, the depth counter code is borked
func TestRoll_shouldPreventInfiniteReferenceRecursion(t *testing.T) {
	ymlf := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - "flat item 1 - {@TestTable_Range}"`

	ymlr := `
  definition:
    name: TestTable_Range
    type: range
    roll: 1d4
    note: this is an optional note
  content:
    - "{1-4}range item 1 - {@TestTable_Flat1}"`

	repo := NewTableRepository()
	repo.AddTable([]byte(ymlf))
	repo.AddTable([]byte(ymlr))
	repo.Roll("TestTable_Range", 1)
	t.Log("External reference recursion test PASS")
}

//if this test hangs, the depth counter code is borked
func TestRoll_shouldPreventInfinteInlineRecursion(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - "item 1 - {#1}"
  inline:
    - id: 1
      content:
        - "inline 1 - {@TestTableFlat}"`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.Roll("TestTable_Flat", 1)
	t.Log("Inline reference recursion test PASS")
}

func TestRoll_shouldFailFastOnTooHighRollCount(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - "item 1: {@TestTable_Flat}"`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	tr := repo.Roll("TestTable_Flat", 1000)
	if len(tr.Result) != 0 {
		t.Error("Unexpected results present")
	}
	if len(tr.Log) != 1 {
		t.Error("Missing or extra Log information")
	}
	if !strings.HasPrefix(tr.Log[0], "Too many rolls requested, max is") {
		t.Error("Log message missing")
	}
}

func TestRollDice_shouldCalcProperly(t *testing.T) {

	//this is not the Worlds Greatest Test but it does stress the code a bit
	//and will find glaring issues (Idx oo Bounds, obvious algo errors)
	//it is hard to test randomizers...
	data := []*rollTestData{toRTD("1d6", 1, 6), toRTD("3d6", 3, 18),
		toRTD("3d6 - 3", 0, 15), toRTD("1d6 * 100", 100, 600), toRTD("3d1", 3, 3),
		toRTD("3d1 + 3", 6, 6), toRTD("1d1 - 7", -6, -6), toRTD("1d1 - 1d1 * 2", 0, 0)}
	ee := newExecutionEngine()

	for i := 1; i <= diceCycleCount; i++ {
		for _, d := range data {
			vr := validate.NewValidationResult()
			dpr := dice.ValidateDiceExpr(d.expr, "TestSection", vr)
			if !vr.Valid() {
				t.Errorf("Bad test case: %s has error: %s", d.expr, vr.Errors[0])
			}
			total := ee.rollDice(dpr)
			if total < d.low || total > d.high {
				t.Errorf("Roll of: %s generated an unexpected result: %d", d.expr, total)
			}
		}
	}
}

type rollTestData struct {
	expr string
	low  int
	high int
}

func toRTD(expr string, low, high int) *rollTestData {
	return &rollTestData{
		expr: expr,
		low:  low,
		high: high,
	}
}
