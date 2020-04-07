package table

import (
	"fmt"
	"strings"
	"tablib/validate"
)

//Table is a table
type Table struct {
	Definition *DefinitionPart `yaml:"definition"`
	RawContent []string        `yaml:"content"`
	Inline     []*InlinePart   `yaml:"inline"`

	IsValid       bool
	IsInlineTable bool
	RangeContent  []*rangedContent
}

const (
	definitionSection = "Definition"
	inlineSection     = "Inline"
	contentSection    = "Content"

	//TableOpRoll represents a roll operation on a table
	TableOpRoll = "roll"

	//TableOpPick represents a pick operation on a table
	TableOpPick = "pick"
)

//Validate ensures the table is valid and parses some aspects if it makes
//sense to do so at validation
func (t *Table) Validate() *validate.ValidationResult {
	vr := validate.NewValidationResult()

	//validate and parse defintion
	t.validateDefinition(vr)

	//validate and parse Inline table(s) if needed
	if len(t.Inline) > 0 {
		t.validateInline(vr)
	}

	//check table internal consistency for inlineSection
	t.validateInternalInlineConsistency(vr)

	//the table is reasonably valid at this point, check content
	if vr.Valid() {
		t.validateContent(vr)
	}

	//have there been any actual validation errors? If so, mark table as Invalid
	t.IsValid = true
	if !vr.IsValid {
		t.IsValid = false
	}
	return vr
}

func (t *Table) validateInternalInlineConsistency(vr *validate.ValidationResult) {

	idsUsed := make(map[string]struct{})
	idsDefined := make([]string, 0, 1)
	for _, rc := range t.RawContent {

		if InlineCalledPattern.MatchString(rc) {
			allMatches := InlineCalledPattern.FindAllStringSubmatch(rc, -1)

			//for each inline table reference, add it to a set of Ids for later comparison
			for i := 0; i < len(allMatches); i++ {
				aMatch := allMatches[i][1]
				left := strings.TrimPrefix(aMatch, "{#")
				idAsString := strings.TrimSuffix(left, "}")
				idsUsed[idAsString] = struct{}{}
			}
		}
	}

	//collect all the defined inline tables
	for _, il := range t.Inline {
		idsDefined = append(idsDefined, il.ID)
	}

	//ensure each used id has a coorisponding inline def
	for uid := range idsUsed {
		found := false
		for _, did := range idsDefined {
			if uid == did {
				found = true
				break
			}
		}
		if !found {
			vr.Fail(contentSection, fmt.Sprintf("Inline table ID: %s is referenced but not defined", uid))
		}
	}

	//warn if an inline table is defined but not used
	for _, did := range idsDefined {
		found := false
		for uid := range idsUsed {
			if did == uid {
				found = true
				break
			}
		}
		if !found {
			vr.Warn(inlineSection, fmt.Sprintf("Inline table ID: %s is defined but not referenced", did))
		}
	}
}
