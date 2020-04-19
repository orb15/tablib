package table

import (
	"fmt"
	"regexp"
	"strings"
	"tablib/util"
	"tablib/validate"
)

var (
	//InlineCalledPattern represents syntax for an inline table call
	InlineCalledPattern = regexp.MustCompile("\\{#([0-9]+)\\}")
	//ExternalCalledPattern represents syntax for an external table call
	ExternalCalledPattern = regexp.MustCompile("\\{@(.*)\\}")
	//PickCalledPattern represents syntax for a pick table call
	PickCalledPattern = regexp.MustCompile("\\{([0-9]+)!(.*)\\}")
)

//ValidateContent ensures the content portion of the table is well-formed
func (t *Table) ValidateContent(vr *validate.ValidationResult) {

	//content section must exist
	if len(t.RawContent) == 0 {
		vr.Fail(contentSection, "A table must have content")
	}

	var allContent []string
	switch t.Definition.TableType {
	case "range":
		//validate the ranges expressed in a ranged table and parse those valid ranges
		//this must be done before the content section of a table can be examined for
		//proper references since the range expressions (eg {2-3}) appear to be
		//invalid table references
		t.validateRanges(vr)
		//now that ranges are validated and parsed, store the actual ranged
		//content for firther Validation
		allContent = make([]string, 0, len(t.RangeContent))
		for _, rc := range t.RangeContent {
			allContent = append(allContent, rc.Content)
		}
	case "flat":
		allContent = t.RawContent
	}

	//allContent contains the actual content of the table. Ensure all tablerefs
	//are valid
	for _, c := range allContent {
		t.validateContentTableRefPairs(c, vr) //do we have closed {}?
	}

	//at this point we can check for valid table refs - if no failures so far
	if vr.Valid() {
		for _, c := range allContent {
			t.validateContentTableRefs(c, vr) //do we have valid tableref syntax?
		}
	}
}

func (t *Table) validateContentTableRefs(entry string, vr *validate.ValidationResult) {
	parts, found := util.FindNextTableRef(entry)
	for found {
		if matches := ExternalCalledPattern.FindStringSubmatch(parts[1]); matches != nil {
			util.IsValidIdentifier(matches[1], parts[1], contentSection, vr)
			parts, found = util.FindNextTableRef(parts[2])
			continue
		}
		if InlineCalledPattern.MatchString(parts[1]) {
			//inline eferences are validated elsewhere
			parts, found = util.FindNextTableRef(parts[2])
			continue
		}
		if matches := PickCalledPattern.FindStringSubmatch(parts[1]); matches != nil {
			util.IsValidIdentifier(matches[2], parts[1], contentSection, vr)
			parts, found = util.FindNextTableRef(parts[2])
			continue
		}
		vr.Fail(contentSection, fmt.Sprintf("Invalid table ref: %s", parts[1]))
		parts, found = util.FindNextTableRef(parts[2])
	}
}

func (t *Table) validateContentTableRefPairs(entry string, vr *validate.ValidationResult) {

	//loop over string, ensuring {} occur in closed pairs
	reader := strings.NewReader(entry)
	state := 0 //a state = 0 means no { or } encountered
	for reader.Len() > 0 {
		c, _ := reader.ReadByte()
		if c == '{' {
			if state == 0 {
				state = 1 //we have an open
			} else {
				vr.Fail(contentSection, "Unexpected open brace {")
				return //this contnt line is likely to be a mess, stop checking
			}
		}
		if c == '}' {
			if state == 1 {
				state = 0
			} else {
				vr.Fail(contentSection, "Unexpected close brace }")
				return //this contnt line is likely to be a mess, stop checking
			}
		}
	}
	//if we reach the end of the string, state should be 0
	if state != 0 {
		vr.Fail(contentSection, "Unclosed open brace {")
	}
}
