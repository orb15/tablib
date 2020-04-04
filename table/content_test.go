package table

import (
	"tablib/validate"
	"testing"
)

func TestContent_shouldRejectEmptyContent1(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note
  content:`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
}

func TestContent_shouldRejectEmptyContent2(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    note: this is an optional note`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
}

func TestContent_shouldRejectMalformedTableRefs(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1
    - item 2
    - item 3`

	testContent := []string{"item 1 {", "item 1 {{", "item 1 }", "}{", "{",
		"hell{ worl{ d", "unexpected} is open{}", "mid}of string{}"}

	tb := tableFromYaml(yml, t)

	for _, tc := range testContent {
		vr := validate.NewValidationResult()
		tb.validateContentTableRefPairs(tc, vr)
		failOnNoErrors(vr, t)
		equals(vr.ErrorCount(), 1, t)
	}
}
