package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//DiceParseResult holds a result of parsing a Dice expression
type DiceParseResult struct {
	Count    int
	DieType  int
	Operator string // "+" or "-" or "none"
}

var (
	xdyPattern = regexp.MustCompile("^([1-9][0-9]*)d([1-9][0-9]*)$")
)

//ValidateDiceExpr validates and parses a dice expression
func ValidateDiceExpr(diceExpr, section string, vr *ValidationResult) []*DiceParseResult {

	//safety - should have been checked before this wass called
	if diceExpr == "" {
		vr.Fail(section, "Programming error: call to checkDice with empty dice string!")
		return nil
	}

	parseResults := make([]*DiceParseResult, 0, 1)
	components := strings.Split(diceExpr, " ")

	//check that the length of components is always odd eg xdy or xdy +- mdn
	//should always be an odd number. If this fails to be true, then we have been
	//provided something like xdy + (missing the mdn portion)
	if len(components)%2 == 0 { //even!
		vr.Fail(section, fmt.Sprintf("Malformed die expression: %s", diceExpr))
		return nil
	}

	//consuming parser to validate & parse expressions: xdy +- mdn +- ...
	for len(components) > 0 {
		if xdyPattern.MatchString(components[0]) {
			matches := xdyPattern.FindStringSubmatch(components[0])
			count, err := strconv.Atoi(matches[1])
			if err != nil {
				vr.Fail(section, fmt.Sprintf("Programming error: unable to parse count portion of dice that should be pre-validated: %s!", diceExpr))
				return nil
			}
			dieType, err := strconv.Atoi(matches[2])
			if err != nil {
				vr.Fail(section, fmt.Sprintf("Programming error: unable to parse dietype portion of dice that should be pre-validated: %s!", diceExpr))
				return nil
			}
			pr := &DiceParseResult{
				Count:   count,
				DieType: dieType,
			}
			if len(components) > 1 {
				switch components[1] {
				case "+", "-":
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
