package table

import (
	"tablib/validate"
	"testing"
)

func TestInlineValidation_shouldRejectMissingInlineID(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1
    - item 2
    - item 3
  inline:
    - id:
      content:
        - Rare
        - Extremely Rare`

	tb := tableFromYaml(yml, t)
	vr := validate.NewValidationResult()
	tb.validateInline(vr)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestInlineValidation_shouldRejectNonnumericInlineID(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {#foo} //this is not the value under test tho it is invalid
    - item 2
    - item 3
  inline:
    - id: foo
      content:
        - Rare
        - Extremely Rare`

	tb := tableFromYaml(yml, t)
	vr := validate.NewValidationResult()
	tb.validateInline(vr)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTestInlineValidation_shouldRejectBadInlineID1(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {#0} //this is not the value under test tho it is invalid
    - item 2
    - item 3
  inline:
    - id: 0
      content:
        - Rare
        - Extremely Rare`

	tb := tableFromYaml(yml, t)
	vr := validate.NewValidationResult()
	tb.validateInline(vr)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTestInlineValidation_shouldRejectBadInlineID2(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {#-1} //this is not the value under test tho it is invalid
    - item 2
    - item 3
  inline:
    - id: -1
      content:
        - Rare
        - Extremely Rare`

	tb := tableFromYaml(yml, t)
	vr := validate.NewValidationResult()
	tb.validateInline(vr)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTestInlineValidation_shouldRejectDuplicateInlineIDs(t *testing.T) {
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
        - Extremely Rare
    - id: 1
      content:
        - Rare
        - Extremely Rare`

	tb := tableFromYaml(yml, t)
	vr := validate.NewValidationResult()
	tb.validateInline(vr)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTestInlineValidation_shouldRejectMissingInlineContent(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {#2}
    - item 2
    - item 3
  inline:
    - id: 2
      content:`

	vr := validateFromYaml(yml, t)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestTestInlineValidation_shouldConvertInlineTableToInternalFormat(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
  content:
    - item 1 - {#1}
    - item 2 - {#2}
    - item 3
  inline:
    - id: 1
      content:
        - Rare
        - Extremely Rare
        - Hyper Rare
        - Pretty damned rare if you ask me
    - id: 2
      content:
        - foo
        - bar`

	tb := tableFromYaml(yml, t)
	vr := validate.NewValidationResult()
	tb.validateInline(vr)
	failOnErrors(vr, t)

	equals(len(tb.Inline), 2, t)
	equals(tb.Inline[0].ID, "1", t)
	equals(tb.Inline[0].FullyQualifiedName, "TestTable_Flat.1", t)
	equals(len(tb.Inline[0].Content), 4, t)
	equals(tb.Inline[0].Content[0], "Rare", t)
	equals(tb.Inline[0].Content[1], "Extremely Rare", t)
	equals(tb.Inline[0].Content[2], "Hyper Rare", t)
	equals(tb.Inline[0].Content[3], "Pretty damned rare if you ask me", t)
	equals(tb.Inline[1].ID, "2", t)
	equals(tb.Inline[1].FullyQualifiedName, "TestTable_Flat.2", t)
	equals(len(tb.Inline[1].Content), 2, t)
	equals(tb.Inline[1].Content[0], "foo", t)
	equals(tb.Inline[1].Content[1], "bar", t)
}
