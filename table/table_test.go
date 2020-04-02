package table

import (
	"fmt"
	"tablib/validate"
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

func TestTable_shouldAcceptWellformedInlineTable(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {#1}
    - item 2
    - item 3
  inline:
    - id: 1
      content:
        - Rare
        - Extremely Rare`

	tb := tableFromYaml(yml, t)
	vr := tb.Validate()
	failOnErrors(vr, t)
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

func TestTable_shouldRejectEmptyContent1(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
}

func TestTable_shouldRejectEmptyContent2(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
}

func TestTable_shouldRejectBadReflineID1(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {#1}
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
    - item 1 - {#2}
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

/* ***********************************************
* Test Helpers
* ***********************************************/

func validateFromYaml(rawYaml string, t *testing.T) *validate.ValidationResult {
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

func failOnErrors(vr *validate.ValidationResult, t *testing.T) {
	if !vr.Valid() {
		t.Error("Expected no validation errors")
	}
}

func failOnNoErrors(vr *validate.ValidationResult, t *testing.T) {
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
