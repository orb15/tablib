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
    name: TestTable_Range
    type: range
    roll: 1d4
    note: this is an optional note
  content:
    - "{1-4}item 1"`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	tr := repo.Roll("TestTable_Range", 1)
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

func TestRoll_shouldhandleRangeDiceMismatch(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Range
    type: range
    roll: 3d4
    note: this is an optional note
  content:
    - "{1} item 1"
    - "{2} item 2"`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	tr := repo.Roll("TestTable_Range", 1)
	if len(tr.Result) != 1 {
		t.Error("Unexpected result count")
	}
	if !strings.HasPrefix(tr.Result[0], "ERROR: roll of") {
		t.Error("Should have received error but did not")
	}
}

func TestRoll_shouldMultiRollAsExpectedRange(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Range
    type: range
    roll: 1d4
    note: this is an optional note
  content:
    - "{1-4}item 1"`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	tr := repo.Roll("TestTable_Range", 3)
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

func TestDiceRef_shouldResolveDiceAsExpected(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - "item 1: {$1d1} {$1d8 * 0}"`

	repo := NewTableRepository()
	vr, err := repo.AddTable([]byte(yml))
	if err != nil {
		t.Error(err)
	}
	if !vr.Valid() {
		t.Error(vr.Errors[0])
	}
	tr := repo.Roll("TestTable_Flat", 1)
	if len(tr.Result) != 1 {
		t.Error("Unexpected result count")
	}
	if tr.Result[0] != "item 1: 1 0" {
		t.Error("Wrong result from table")
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

func TestRoll_shouldFailtoExpandBadRollTableRef(t *testing.T) {
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

func TestRoll_shouldFailtoExpandBadPickTableRef(t *testing.T) {
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
    - "{1-4}range item 1 - {1!TestTable_Flat1}"`

	repo := NewTableRepository()
	repo.AddTable([]byte(ymlf))
	repo.AddTable([]byte(ymlr))
	tr := repo.Roll("TestTable_Range", 1)
	if len(tr.Result) != 1 {
		t.Error("Unexpected result count")
	}
	if tr.Result[0] != "range item 1 -  --BADREF: {1!TestTable_Flat1}--" {
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
func TestRoll_shouldPreventInfiniteSelfRecursionOnRoll(t *testing.T) {
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
	t.Log("Roll self recursion test PASS")
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

func TestPick_shouldPickAsExpected(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1
    - item 2`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	tr := repo.Pick("TestTable_Flat", 1)
	if len(tr.Result) != 1 {
		t.Error("Pick not returing proper amount of data")
	}
	if len(tr.Log) != 1 {
		t.Error("Pick not logging properly")
	}
	if tr.Result[0] != "item 1" && tr.Result[0] != "item 2" {
		t.Error("Pick returned invalid result")
	}
}

func TestPick_shouldMultiPickAsExpected(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1
    - item 2
    - item 3`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	tr := repo.Pick("TestTable_Flat", 2)
	if len(tr.Result) != 1 {
		t.Error("Pick not returing proper amount of data")
	}
	if len(tr.Log) != 1 {
		t.Error("Pick not logging properly")
	}

	//ensure results are true pick
	parts := strings.Split(tr.Result[0], "|")
	if len(parts) != 2 {
		t.Error("Did not pick the desired number")
	}
	if parts[0] == parts[1] {
		t.Error("Pick was not unique")
	}
}

func TestPick_shouldPickAllAndWarn1(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1
    - item 2`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	tr := repo.Pick("TestTable_Flat", 2)
	if len(tr.Result) != 1 {
		t.Error("Pick not returing proper amount of data")
	}
	if len(tr.Log) != 2 {
		t.Error("Pick not logging properly")
	}
	if tr.Result[0] != "item 1|item 2" {
		t.Error("Pick returned invalid result")
	}
}

func TestPick_shouldPickAllAndWarn2(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1
    - item 2`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	tr := repo.Pick("TestTable_Flat", 3)
	if len(tr.Result) != 1 {
		t.Error("Pick not returing proper amount of data")
	}
	if len(tr.Log) != 2 {
		t.Error("Pick not logging properly")
	}
	if tr.Result[0] != "item 1|item 2" {
		t.Error("Pick returned invalid result")
	}
}

func TestPick_shouldErrorOnRangeTable(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Range
    type: range
    roll: 1d4
    note: this is an optional note
  content:
    - "{1-4}item 1"`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	tr := repo.Pick("TestTable_Range", 1)
	if len(tr.Result) != 1 {
		t.Error("Unexpected result count")
	}
	if tr.Result[0] != "Pick on range table not allowed" {
		t.Errorf("Result missing error message")
	}
	if len(tr.Log) != 2 {
		t.Error("Wrong Log info captured")
	}
	if tr.Log[1] != "Pick requested on ranged table: TestTable_Range" {
		t.Error("Bad log message")
	}
}

func TestPick_shouldFailOnMissingTable(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	tr := repo.Pick("TestTable_Flat1", 1)
	if len(tr.Result) != 0 {
		t.Error("Unexpected result count")
	}
	if len(tr.Log) != 1 {
		t.Error("Wrong Log info captured")
	}
}

//if this test hangs, the depth counter code is borked
func TestPick_shouldPreventInfiniteSelfRecursionOnPick(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - "item 1: {1!TestTable_Flat}"`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml))
	repo.Roll("TestTable_Flat", 1)
	t.Log("Pick self recursion test PASS")
}

func TestRecurse_shouldGenerallyRecurse1(t *testing.T) {
	yml1 := `
  definition:
    name: Flat1
    type: flat
    note: this is an optional note
  content:
    - "{#1}"
  inline:
    - id: 1
      content:
        - "Flat1-inline: {@Range1}"`

	yml2 := `
  definition:
    name: Range1
    type: range
    roll: 1d4
    note: this is an optional note
  content:
    - "{1-4}Range1|{2!Flat2}"`

	yml3 := `
  definition:
    name: Flat2
    type: flat
    note: this is an optional note
  content:
    - Flat2-1
    - Flat2-2`

	repo := NewTableRepository()
	repo.AddTable([]byte(yml1))
	repo.AddTable([]byte(yml2))
	repo.AddTable([]byte(yml3))
	tr := repo.Roll("Flat1", 1)

	if len(tr.Result) != 1 {
		t.Error("Unexpected Result count")
	}
	if len(tr.Log) != 8 {
		t.Error("Unexpected Log count")
	}
	if !strings.HasPrefix(tr.Result[0], "Flat1-inline: Range1|") {
		t.Error("First portion of Result is wrong")
	}

	firstPipe := strings.Index(tr.Result[0], "|")
	pickString := tr.Result[0][firstPipe+1 : len(tr.Result[0])]
	if pickString != "Flat2-1|Flat2-2" && pickString != "Flat2-2|Flat2-1" {
		t.Error("Invalid pick string")
	}
}

func TestRecurse_shouldGenerallyMultiRecurse1(t *testing.T) {
	yml1 := `
  definition:
    name: Flat1
    type: flat
    note: this is an optional note
  content:
    - "Pick from Flat2 by inline: |{#1}| but this is direct: |{1!Flat2}| and this is a dice roll: |{$1d1}"
  inline:
    - id: 1
      content:
        - "Flat1-inline: {1!Flat2}"`

	yml2 := `
  definition:
    name: Flat2
    type: flat
    note: this is an optional note
  content:
    - Flat2-1
    - Flat2-2`

	repo := NewTableRepository()
	vr, _ := repo.AddTable([]byte(yml1))
	if !vr.Valid() {
		t.Error("Invalid")
	}
	repo.AddTable([]byte(yml2))
	tr := repo.Roll("Flat1", 1)

	if len(tr.Result) != 1 {
		t.Error("Unexpected Result count")
	}
	if len(tr.Log) != 7 {
		t.Error("Unexpected Log count")
	}
	if !strings.HasPrefix(tr.Result[0], "Pick from Flat2 by inline: |") {
		t.Error("First portion of Result is wrong")
	}
	pickParts := strings.Split(tr.Result[0], "|")
	if len(pickParts) != 6 {
		t.Error("Unexpected output")
	}
	if pickParts[1] != "Flat1-inline: Flat2-1" && pickParts[1] != "Flat1-inline: Flat2-2" {
		t.Error("Bad first pick")
	}
	if pickParts[3] != "Flat2-1" && pickParts[3] != "Flat2-2" {
		t.Error("Bad second pick")
	}
	if pickParts[5] != "1" {
		t.Error("Bad dice eval")
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
