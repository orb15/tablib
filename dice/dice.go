package dice

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"tablib/validate"
)

//ParseResult holds a result of parsing a Dice expression
type ParseResult struct {
	Count    int
	DieType  int
	Operator string // "+" or "-" or "*" or "none"
}

var (
	xdyPattern = regexp.MustCompile("^([1-9][0-9]*)d([1-9][0-9]*)$")
	intPattern = regexp.MustCompile("^[0-9]+$")
)

//ValidateDiceExpr validates and parses a dice expression
func ValidateDiceExpr(diceExpr, section string, vr *validate.ValidationResult) []*ParseResult {

	//safety - should have been checked before this was called
	if diceExpr == "" {
		vr.Fail(section, "Programming error: call to checkDice with empty dice string!")
		return nil
	}

	parseResults := make([]*ParseResult, 0, 1)
	components := strings.Split(diceExpr, " ")

	//check that the length of components is always odd eg xdy or xdy +- mdn
	//should always be an odd number. If this fails to be true, then we have been
	//provided something like xdy + (missing the mdn portion)
	if len(components)%2 == 0 { //even!
		vr.Fail(section, fmt.Sprintf("Malformed die expression: %s", diceExpr))
		return nil
	}

	//consuming parser to validate & parse expressions: xdy +- mdn +- ... +- N
	for len(components) > 0 {

		//first, check to see if we have a constant ending eg 2d6 + 3
		//if we do, we want to consume that specially
		if intPattern.MatchString(components[0]) {
			constant, _ := strconv.Atoi(components[0]) //no err here as regex protects
			pr := &ParseResult{
				Count:    constant,
				DieType:  0,
				Operator: "none",
			}
			parseResults = append(parseResults, pr)
			components = components[1:] //consume the constant we just parsed, []string now len = 0
			continue
		}

		//we are expecting an xdy pattern at this point or we fail to validate
		if xdyPattern.MatchString(components[0]) {
			matches := xdyPattern.FindStringSubmatch(components[0])
			count, _ := strconv.Atoi(matches[1])   //no err here as regex protects
			dieType, _ := strconv.Atoi(matches[2]) //no err here as regex protects
			pr := &ParseResult{
				Count:   count,
				DieType: dieType,
			}
			if len(components) > 1 {
				switch components[1] {
				case "+", "-", "*":
					pr.Operator = components[1]
					components = components[2:] //consume the die we just parsed and follow-on operator
				default:
					vr.Fail(section, fmt.Sprintf("Invalid die operator: %s in %s", components[1], diceExpr))
					return nil
				}
			} else {
				pr.Operator = "none"
				components = components[1:] //consume the die we just parsed, []string now len = 0
			}
			parseResults = append(parseResults, pr)
		} else {
			vr.Fail(section, fmt.Sprintf("Invalid dice expression: %s", diceExpr))
			return nil
		}
	}

	return parseResults
}
