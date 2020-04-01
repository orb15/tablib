package table

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"tablib/validate"
)

type rangedContent struct {
	Low     int
	High    int
	Content string
}

const (
	contentSection = "Content"
)

var (
	rangedContentPattern = regexp.MustCompile("^\\{([0-9]+)-([0-9]+)\\}.*$")
	fixedContentPattern  = regexp.MustCompile("^\\{([0-9]+)\\}.*$")
)

func (t *Table) validateRangeContent(vr *validate.ValidationResult) {

	//set up to store parsed ranged content
	allContent := make([]*rangedContent, 0, 1)

	for _, rc := range t.RawContent {
		if rangedContentPattern.MatchString(rc) { //{x-y}
			matches := rangedContentPattern.FindStringSubmatch(rc)
			lowVal, _ := strconv.Atoi(matches[1])  //no err, regex protects this
			highVal, _ := strconv.Atoi(matches[2]) //no err, regex protects this
			if lowVal >= highVal {
				vr.Fail(contentSection, fmt.Sprintf("Invalid range: %d greater or equal to %d", lowVal, highVal))
			}
			splitStrings := strings.SplitAfterN(rc, "}", 2)
			rgCont := &rangedContent{
				Low:     lowVal,
				High:    highVal,
				Content: splitStrings[1],
			}
			allContent = append(allContent, rgCont)
		} else if fixedContentPattern.MatchString(rc) { // range of a single value eg {x}}
			matches := fixedContentPattern.FindStringSubmatch(rc)
			onlyVal, _ := strconv.Atoi(matches[1]) //no err, regex protects this
			splitStrings := strings.SplitAfterN(rc, "}", 2)
			rgCont := &rangedContent{
				Low:     onlyVal,
				High:    onlyVal,
				Content: splitStrings[1],
			}
			allContent = append(allContent, rgCont)
		} else {
			vr.Fail(contentSection, fmt.Sprintf("Invalid ranged content: %s", rc))
		}

	}
	t.RangeContent = allContent

	//final validation - make sure no ranges are out of order or overlap
	overlap := false
	for i := 0; i < len(t.RangeContent)-1; i++ {
		if t.RangeContent[i].High >= t.RangeContent[i+1].Low && !overlap {
			vr.Fail(contentSection, "this table has a range overlap or ordering issue")
			overlap = true //supress similar failures - out of order will cause a mess
		}
	}
}
