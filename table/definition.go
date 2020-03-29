package table

import (
	"fmt"
	"tablib/util"
)

const (
	definitionSection = "Definition"
)

//DefinitionPart holds the table definition or header
type DefinitionPart struct {
	Name       string `yaml:"name"`
	TableType  string `yaml:"type"`
	Roll       string `yaml:"roll"`
	diceParsed []*diceParseResult
}

func (t *Table) validateDefinition(vr *util.ValidationResult) {
	checkEmpty(t.Definition.Name, "Name", definitionSection, vr)
	validIdentifier(t.Definition.Name, "Name", definitionSection, vr)
	checkEmpty(t.Definition.TableType, "TableType", definitionSection, vr)

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

	if t.Definition.Roll != "" {
		parseResults := checkDice(t.Definition.Roll, definitionSection, vr)
		if parseResults != nil {
			t.Definition.diceParsed = parseResults
		}
	}
}
