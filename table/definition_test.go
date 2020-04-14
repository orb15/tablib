package table

import (
	"tablib/validate"
	"testing"
)

func TestDefinitionValidation_shouldRejectTableType1(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type:
  content:
    - item 1
    - item 2
    - item 3`

	tb := tableFromYaml(yml, t)
	vr := validate.NewValidationResult()
	tb.validateDefinition(vr)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 2, t)
}

func TestDefinitionValidation_shouldRejectTableType2(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: Flat
  content:
    - item 1
    - item 2
    - item 3`

	tb := tableFromYaml(yml, t)
	vr := validate.NewValidationResult()
	tb.validateDefinition(vr)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestDefinitionValidation_shouldRejectTableType3(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: falt
  content:
    - item 1
    - item 2
    - item 3`

	tb := tableFromYaml(yml, t)
	vr := validate.NewValidationResult()
	tb.validateDefinition(vr)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestDefinitionValidation_shouldWarnOnRollFlatTable(t *testing.T) {
	yml := `
  definition:
    name: TestTable_Flat
    type: flat
    roll: 1d6
  content:
    - item 1
    - item 2
    - item 3`

	tb := tableFromYaml(yml, t)
	vr := validate.NewValidationResult()
	tb.validateDefinition(vr)
	failOnErrors(vr, t)
	equals(vr.HasWarnings, true, t)
	equals(vr.WarnCount(), 1, t)
}

func TestDefinitionValidation_shouldAcceptValidDefinitionName(t *testing.T) {
	yml := `
  definition:
    name: Table_1_2020_kobe
    type: flat
  content:
    - item 1
    - item 2
    - item 3`

	tb := tableFromYaml(yml, t)
	vr := validate.NewValidationResult()
	tb.validateDefinition(vr)
	failOnErrors(vr, t)
}

func TestDefinitionValidation_shouldRejectInvalidDefinitionName(t *testing.T) {
	yml := `
  definition:
    name:
    type: flat
  content:
    - item 1
    - item 2
    - item 3`

	tb := tableFromYaml(yml, t)
	vr := validate.NewValidationResult()
	tb.validateDefinition(vr)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestDefinitionValidation_shouldAcceptValidDice(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 2d8
  content:
    - '{1-2}item 1'
    - '{3-4}item 2'
    - '{5-6}item 3'`

	vr := validateFromYaml(yml, t)
	failOnErrors(vr, t)
}

func TestDefinitionValidation_shouldRejectMissingDice(t *testing.T) {
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

func TestDefinitionValidation_shouldUniquifyAndLowercaseTags(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 2d8
    tags:
      - test
      - Test
      - foo
      - foO
      - bar
  content:
    - '{1-2}item 1'
    - '{3-4}item 2'
    - '{5-6}item 3'`

	tb := tableFromYaml(yml, t)
	vr := tb.Validate()
	failOnErrors(vr, t)
	equals(len(tb.Definition.Tags), 3, t)
}
