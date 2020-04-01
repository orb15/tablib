package table

import (
	"fmt"
	"tablib/dice"
	"tablib/util"
	"tablib/validate"
)

const (
	definitionSection = "Definition"
)

//DefinitionPart holds the table definition or header
type DefinitionPart struct {
	Name      string `yaml:"name"`
	Note      string `yaml:"note"`
	TableType string `yaml:"type"`
	Roll      string `yaml:"roll"`

	DiceParsed []*dice.ParseResult
}

func (t *Table) validateDefinition(vr *validate.ValidationResult) {
	util.IsValidIdentifier(t.Definition.Name, "Name", definitionSection, vr)
	util.IsNotEmpty(t.Definition.TableType, "TableType", definitionSection, vr)

	//ensure valid table type, ensure alignment between table type and roll
	//information
	switch t.Definition.TableType {
	case "flat":
		if t.Definition.Roll != "" {
			vr.Warn(definitionSection, fmt.Sprintf("Roll defined but not used for this table type"))
		}
	case "range":
		if t.Definition.Roll == "" {
			vr.Fail(definitionSection, fmt.Sprintf("Roll must be defined for this table type"))
		}
	default:
		vr.Fail(definitionSection, fmt.Sprintf("Unknown TableType: %s", t.Definition.TableType))
	}

	//if a roll is provided, make sure it is valid
	if t.Definition.Roll != "" {
		parseResults := dice.ValidateDiceExpr(t.Definition.Roll, definitionSection, vr)
		if parseResults != nil {
			t.Definition.DiceParsed = parseResults
		}
	}
}
