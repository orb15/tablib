package table

import (
	"fmt"
	"tablib/util"
)

//InlinePart holds info about an inline table
type InlinePart struct {
	ID                 string   `yaml:"id"`
	Content            []string `yaml:"content"`
	fullyQualifiedName string
}

const (
	inlineSection = "Inline"
)

func (t *Table) validateInline(vr *util.ValidationResult) *util.ValidationResult {
	for _, il := range t.Inline {
		checkEmpty(il.ID, "ID", inlineSection, vr)
		if len(il.Content) <= 0 {
			vr.Fail(inlineSection, fmt.Sprintf("Inline table with id: %s is empty", il.ID))
		}
		il.fullyQualifiedName = util.BuildFullName(t.Definition.Namespace, t.Definition.Name, il.ID)
	}
	return vr
}
