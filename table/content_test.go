package table

import (
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
