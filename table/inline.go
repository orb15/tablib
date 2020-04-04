package table

import (
	"fmt"
	"sort"
	"strconv"
	"tablib/util"
	"tablib/validate"
)

//InlinePart holds info about an inline table
type InlinePart struct {
	ID      string   `yaml:"id"`
	Content []string `yaml:"content"`

	FullyQualifiedName string
}

func (t *Table) validateInline(vr *validate.ValidationResult) {

	idsDefined := make([]string, 0, len(t.Inline))
	//ensure ID and content are both defined for each inline table
	for _, il := range t.Inline {
		idVal, err := strconv.Atoi(il.ID)
		if err != nil {
			vr.Fail(inlineSection, fmt.Sprintf("Invalid ID for Inline table: %s", il.ID))
		} else {
			if idVal <= 0 {
				vr.Fail(inlineSection, fmt.Sprintf("Invalid ID for Inline table: %s", il.ID))
			}
			idsDefined = append(idsDefined, il.ID)
		}
		if len(il.Content) <= 0 {
			vr.Fail(inlineSection, fmt.Sprintf("Inline table with id: %s is empty", il.ID))
		}
		il.FullyQualifiedName = util.BuildFullName(t.Definition.Name, il.ID)
	}

	//ensure uniqueness of inline ids
	sort.Strings(idsDefined)
	for i := 0; i < len(idsDefined)-1; i++ {
		if idsDefined[i] == idsDefined[i+1] {
			vr.Fail(inlineSection, fmt.Sprintf("Inline table ID: %s defined twice", idsDefined[i]))
		}
	}
}
