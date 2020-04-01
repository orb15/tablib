package table

import (
	"fmt"
	"strconv"
	"tablib/util"
)

//InlinePart holds info about an inline table
type InlinePart struct {
	ID      string   `yaml:"id"`
	Content []string `yaml:"content"`

	FullyQualifiedName string
}

const (
	inlineSection = "Inline"
)

func (t *Table) validateInline(vr *util.ValidationResult) {

	//ensure ID and content are both defined
	for _, il := range t.Inline {
		util.IsNotEmpty(il.ID, "ID", inlineSection, vr)
		idVal, err := strconv.Atoi(il.ID)
		if err != nil {
			vr.Fail(inlineSection, fmt.Sprintf("Invalid ID for Inline table: %s", il.ID))
		} else {
			if idVal <= 0 {
				vr.Fail(inlineSection, fmt.Sprintf("Invalid ID for Inline table: %s", il.ID))
			}
		}
		if len(il.Content) <= 0 {
			vr.Fail(inlineSection, fmt.Sprintf("Inline table with id: %s is empty", il.ID))
		}
		il.FullyQualifiedName = util.BuildFullName(t.Definition.Name, il.ID)
	}
}
