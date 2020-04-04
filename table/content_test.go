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

func TestContent_shouldRejectMalformedTablePairs(t *testing.T) {
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

func TestContent_shouldAcceptWellformedTablePairs(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1
    - item 2
    - item 3`

	testContent := []string{"{}", "hello{world}", "{}{}{}", "{@Goo}",
		"{3!Foo}", "{#667}"}

	tb := tableFromYaml(yml, t)

	for _, tc := range testContent {
		vr := validate.NewValidationResult()
		tb.validateContentTableRefPairs(tc, vr)
		failOnErrors(vr, t)
	}
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

	testContent := []string{"{}", "{!2}", "{W@rld}", "good{@Ref} then {!Bad}",
		"Content was {$bad} but then good {#3}", "{3#}"}

	tb := tableFromYaml(yml, t)

	for _, tc := range testContent {
		vr := validate.NewValidationResult()
		tb.validateContentTableRefs(tc, vr)
		failOnNoErrors(vr, t)
		equals(vr.ErrorCount(), 1, t)
	}
}

func TestContent_shouldAcceptWellformedTableRefs(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1
    - item 2
    - item 3`

	testContent := []string{"{@Sloopy}", "{2!Goober}", "good{#3} then {3!Better}",
		"Content was {#2} but then got {@Better} and {@Better_yet}", "perfectly ok"}

	tb := tableFromYaml(yml, t)

	for _, tc := range testContent {
		vr := validate.NewValidationResult()
		tb.validateContentTableRefs(tc, vr)
		failOnErrors(vr, t)
	}
}
