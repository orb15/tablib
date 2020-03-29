package table

import "tablib/util"

//Table is a table
type Table struct {
	Definition *DefinitionPart `yaml:"definition"`
	RawContent []string        `yaml:"content"`
	Inline     []*InlinePart   `yaml:"inline"`

	IsValid      bool
	rangeContent []*rangedContent
}

//Validate ensures the table is valid and parses some aspects if it makes
//sense to do so at validation
func (t *Table) Validate() *util.ValidationResult {
	vr := &util.ValidationResult{
		IsValid:     true,
		HasWarnings: false,
		Errors:      make([]string, 0),
	}

	//validate and parse defintion
	vr = t.validateDefinition(vr)

	//validate and parse content section
	switch t.Definition.TableType {
	case "range":
		vr = t.validateRangeContent(vr)
	}

	//validate and parse Inline table(s) if needed
	if len(t.Inline) > 0 {
		vr = t.validateInline(vr)
	}

	//have there been any actual validation errors? If so, mark table as Invalid
	t.IsValid = true
	if !vr.IsValid {
		t.IsValid = false
	}
	return vr
}
