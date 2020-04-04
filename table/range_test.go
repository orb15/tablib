package table

import (
	"tablib/validate"
	"testing"
)

func TestRangeValidation_shouldRejectInvertedRange(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 1d6
  content:
    - '{1-2}item 1'
    - '{4-3}item 2'
    - '{5-6}item 3'`

	tb := tableFromYaml(yml, t)
	vr := validate.NewValidationResult()
	tb.validateRanges(vr)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestRangeValidation_shouldRejectOverlappingRange(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 1d6
  content:
    - '{1-2}item 1'
    - '{3-4}item 2'
    - '{4-6}item 3'`

	tb := tableFromYaml(yml, t)
	vr := validate.NewValidationResult()
	tb.validateRanges(vr)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestRangeValidation_shouldRejectStrangeRange1(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 1d6
  content:
    - '{1-2}item 1'
    - '{3-4}item 2'
    - '{1-2}item 3'`

	tb := tableFromYaml(yml, t)
	vr := validate.NewValidationResult()
	tb.validateRanges(vr)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestRangeValidation_shouldRejectStrangeRange2(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 1d6
  content:
    - '{1-2}item 1'
    - '{1-2}item 2'
    - '{3-6}item 3'`

	tb := tableFromYaml(yml, t)
	vr := validate.NewValidationResult()
	tb.validateRanges(vr)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestRangeValidation_shouldRejectStrangeRange3(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 1d6
  content:
    - '{1}item 1'
    - '{2}item 2'
    - '{1}item 3'`

	tb := tableFromYaml(yml, t)
	vr := validate.NewValidationResult()
	tb.validateRanges(vr)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestRangeValidation_shouldRejectStrangeRange4(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 1d6
  content:
    - '{1}item 1'
    - '{1-2}item 2'
    - '{2}item 3'`

	tb := tableFromYaml(yml, t)
	vr := validate.NewValidationResult()
	tb.validateRanges(vr)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestRangeValidation_shouldRejectStrangeRange5(t *testing.T) {
	yml := `
  definition:
    name: TestTable
    type: range
    roll: 1d6
  content:
    - '{-1}item 1'
    - '{1-2}item 2'
    - '{3}item 3'`

	tb := tableFromYaml(yml, t)
	vr := validate.NewValidationResult()
	tb.validateRanges(vr)
	failOnNoErrors(vr, t)
	equals(vr.ErrorCount(), 1, t)
}

func TestRangeValidation_shouldConvertRangedTableToInternalFormat(t *testing.T) {
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
	vr := validate.NewValidationResult()
	tb.validateRanges(vr)
	failOnErrors(vr, t)

	equals(len(tb.RangeContent), 4, t)
	equals(tb.RangeContent[1].Low, 3, t)
	equals(tb.RangeContent[1].High, 4, t)
	equals(tb.RangeContent[1].Content, "item 2 is {@SomeTable}", t)
	equals(tb.RangeContent[3].Low, 6, t)
	equals(tb.RangeContent[3].High, 6, t)
	equals(tb.RangeContent[3].Content, "item 4", t)
}
