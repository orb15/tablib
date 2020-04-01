package table

import (
	"fmt"
	"tablib/util"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestTable_shouldAcceptWellformedFlatTable(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:
    - item 1
    - item 2
    - item 3`

	vr := validateFromYaml(yml, t)
	failOnErrors(vr, t)
}

func TestTable_shouldRejectTableType1(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type:
  content:
    - item 1
    - item 2
    - item 3`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 2, t)

}

func TestTable_shouldRejectTableType2(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: Flat
  content:
    - item 1
    - item 2
    - item 3`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectTableType3(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: falt
  content:
    - item 1
    - item 2
    - item 3`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldWarnOnRollFlatTable(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    roll: 1d6
  content:
    - item 1
    - item 2
    - item 3`

	vr := validateFromYaml(yml, t)
	failOnErrors(vr, t)
	equals(vr.HasWarnings, true, t)
	equals(vr.WarnCount(), 1, t)
}

func TestTable_shouldAcceptDefinitionName1(t *testing.T) {
	yml := `
  definition:
    name: Table_1_2020_kobe
    type: flat
  content:
    - item 1
    - item 2
    - item 3`

	vr := validateFromYaml(yml, t)
	failOnErrors(vr, t)
}

func TestTable_shouldRejectDefinitionNameInvalid1(t *testing.T) {
	yml := `
  definition:
    name:
    type: flat
  content:
    - item 1
    - item 2
    - item 3`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectDefinitionNameInvalid2(t *testing.T) {
	yml := `
  definition:
    name: 1Table
    type: flat
  content:
    - item 1
    - item 2
    - item 3`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectDefinitionNameInvalid3(t *testing.T) {
	yml := `
  definition:
    name: Test1-Table
    type: flat
  content:
    - item 1
    - item 2
    - item 3`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectDefinitionNameInvalid4(t *testing.T) {
	yml := `
  definition:
    name: Test Table
    type: flat
  content:
    - item 1
    - item 2
    - item 3`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldAcceptWellformedRangeTable(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 1d6
  content:
    - '{1-2}item 1'
    - '{3-4}item 2'
    - '{5}item 3'
    - '{6}item 4'`

	vr := validateFromYaml(yml, t)
	failOnErrors(vr, t)
}

func TestTable_shouldAcceptDice1(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 3d6
  content:
    - '{1-2}item 1'
    - '{3-4}item 2'
    - '{5-6}item 3'`

	vr := validateFromYaml(yml, t)
	failOnErrors(vr, t)
}

func TestTable_shouldAcceptDice2(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 3d6 + 1d12
  content:
    - '{1-2}item 1'
    - '{3-4}item 2'
    - '{5-6}item 3'`

	vr := validateFromYaml(yml, t)
	failOnErrors(vr, t)
}

func TestTable_shouldAcceptDice3(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 3d6 - 1d12
  content:
    - '{1-2}item 1'
    - '{3-4}item 2'
    - '{5-6}item 3'`

	vr := validateFromYaml(yml, t)
	failOnErrors(vr, t)
}

func TestTable_shouldRejectDice1(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 0d6
  content:
    - '{1-2}item 1'
    - '{3-4}item 2'
    - '{5-6}item 3'`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectDice2(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 1 d6
  content:
    - '{1-2}item 1'
    - '{3-4}item 2'
    - '{5-6}item 3'`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectDice3(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 1d 6
  content:
    - '{1-2}item 1'
    - '{3-4}item 2'
    - '{5-6}item 3'`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectDice4(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 1d6+1d12
  content:
    - '{1-2}item 1'
    - '{3-4}item 2'
    - '{5-6}item 3'`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectDice5(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: d6
  content:
    - '{1-2}item 1'
    - '{3-4}item 2'
    - '{5-6}item 3'`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectMissingDice(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll:
  content:
    - '{1-2}item 1'
    - '{3-4}item 2'
    - '{5-6}item 3'`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectInvertedRange(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 1d6
  content:
    - '{1-2}item 1'
    - '{4-3}item 2'
    - '{5-6}item 3'`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectOverlappingRange(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 1d6
  content:
    - '{1-2}item 1'
    - '{3-4}item 2'
    - '{4-6}item 3'`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectStrangeRange1(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 1d6
  content:
    - '{1-2}item 1'
    - '{3-4}item 2'
    - '{1-2}item 3'`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectStrangeRange2(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 1d6
  content:
    - '{1-2}item 1'
    - '{1-2}item 2'
    - '{3-6}item 3'`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectStrangeRange3(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 1d6
  content:
    - '{1}item 1'
    - '{2}item 2'
    - '{1}item 3'`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectStrangeRange4(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 1d6
  content:
    - '{1}item 1'
    - '{1-2}item 2'
    - '{2}item 3'`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectStrangeRange5(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 1d6
  content:
    - '{-1}item 1'
    - '{1-2}item 2'
    - '{3}item 3'`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldAcceptWellformedInlineTable(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {@1}
    - item 2
    - item 3
  inline:
    - id: 1
      content:
        - Rare
        - Extremely Rare`

	vr := validateFromYaml(yml, t)
	failOnErrors(vr, t)
}

func TestTable_shouldRejectMissinglineID(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {@1}
    - item 2
    - item 3
  inline:
    - id:
      content:
        - Rare
        - Extremely Rare`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 3, t)
}

func TestTable_shouldRejectNonnumericlineID(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {@foo}
    - item 2
    - item 3
  inline:
    - id: foo
      content:
        - Rare
        - Extremely Rare`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectBadReflineID1(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {@1}
    - item 2
    - item 3
  inline:
    - id: 2
      content:
        - Rare
        - Extremely Rare`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
	equals(vr.WarnCount(), 1, t)
}

func TestTable_shouldRejectBadReflineID2(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {@2}
    - item 2
    - item 3
  inline:
    - id: 1
      content:
        - Rare
        - Extremely Rare`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
	equals(vr.WarnCount(), 1, t)
}

func TestTable_shouldRejectBadInlineID1(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {@0}
    - item 2
    - item 3
  inline:
    - id: 0
      content:
        - Rare
        - Extremely Rare`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectBadInlineID2(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {@-1}
    - item 2
    - item 3
  inline:
    - id: -1
      content:
        - Rare
        - Extremely Rare`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectDuplicateInlineIDs(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {@1}
    - item 2 - {@1}
    - item 3
  inline:
    - id: 1
      content:
        - Rare
        - Extremely Rare
    - id: 1
      content:
        - Rare
        - Extremely Rare`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldRejectMissingInlineContent(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {@2}
    - item 2
    - item 3
  inline:
    - id: 2
      content:`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTable_shouldConvertRangedTableToInternalFormat(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 1d6
  content:
    - '{1-2}item 1'
    - '{3-4}item 2 is {@SomeTable}'
    - '{5}item 3'
    - '{6}item 4'`

	tb := tableFromYaml(yml, t)
	vr := tb.Validate()
	failOnErrors(vr, t)

	equals(len(tb.RangeContent), 4, t)
	equals(tb.RangeContent[1].Low, 3, t)
	equals(tb.RangeContent[1].High, 4, t)
	equals(tb.RangeContent[1].Content, "item 2 is {@SomeTable}", t)
	equals(tb.RangeContent[3].Low, 6, t)
	equals(tb.RangeContent[3].High, 6, t)
	equals(tb.RangeContent[3].Content, "item 4", t)
}

func TestTable_shouldConvertInlineTableToInternalFormat(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {@1}
    - item 2 - {@2}
    - item 3
  inline:
    - id: 1
      content:
        - Rare
        - Extremely Rare
    - id: 2
      content:
        - foo
        - bar`
	tb := tableFromYaml(yml, t)
	vr := tb.Validate()
	failOnErrors(vr, t)

	equals(len(tb.Inline), 2, t)
	equals(tb.Inline[1].ID, "2", t)
	equals(tb.Inline[1].FullyQualifiedName, "TestTable_Flat.2", t)
	equals(len(tb.Inline[1].Content), 2, t)
	equals(tb.Inline[1].Content[0], "foo", t)
	equals(tb.Inline[1].Content[1], "bar", t)
}

/* ***********************************************
* Test Helpers
* ***********************************************/

func validateFromYaml(rawYaml string, t *testing.T) *util.ValidationResult {
	table := tableFromYaml(rawYaml, t)
	return table.Validate()
}

func tableFromYaml(rawYaml string, t *testing.T) *Table {
	var table Table
	err := yaml.Unmarshal([]byte(rawYaml), &table)
	if err != nil {
		fmt.Println(err)
		t.Fatal("Bad test YAML, FIX ME!")
	}

	return &table
}

func failOnErrors(vr *util.ValidationResult, t *testing.T) {
	if !vr.Valid() {
		t.Error("Expected no validation errors")
	}
}

func failOnNoErrors(vr *util.ValidationResult, t *testing.T) {
	if vr.Valid() {
		t.Error("Expected validation errors")
	}
}

func equals(have interface{}, want interface{}, t *testing.T) {
	if have == want {
		return
	}
	t.Fatalf("Equals Failed: have: %v want: %v", have, want)
}
