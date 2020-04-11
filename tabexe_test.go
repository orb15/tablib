package tablib

/*
These tests focus on the table executable portions of the repo like
rolling on and picking from tables. Because of my lack of desire to write mocks
(and ensure I am using interfaces everywhere I could be to make this doable),
some of these tests are more like integration tests than unit tests.
*/

import (
	"testing"

	"tablib/dice"
	"tablib/validate"
)

const (
	diceCycleCount = 100 //number of times to run the dice roll test
)

func TestRollDice_shouldCalcProperly(t *testing.T) {

	//this is not the Worlds Greatest Test but it does stress the code a bit
	//and will find glaring issues (Idx oo Bounds, obvious algo errors)
	//it is hard to test randomizers...
	data := []*rollTestData{toRTD("1d6", 1, 6), toRTD("3d6", 3, 18),
		toRTD("3d6 - 3", 0, 15), toRTD("1d6 * 100", 100, 600), toRTD("3d1", 3, 3),
		toRTD("3d1 + 3", 6, 6), toRTD("1d1 - 7", -6, -6), toRTD("1d1 - 1d1 * 2", 0, 0)}
	ee := newExecutionEngine()

	for i := 1; i <= diceCycleCount; i++ {
		for _, d := range data {
			vr := validate.NewValidationResult()
			dpr := dice.ValidateDiceExpr(d.expr, "TestSection", vr)
			if !vr.Valid() {
				t.Errorf("Bad test case: %s has error: %s", d.expr, vr.Errors[0])
			}
			total := ee.rollDice(dpr)
			if total < d.low || total > d.high {
				t.Errorf("Roll of: %s generated an unexpected result: %d", d.expr, total)
			}
		}
	}
}

type rollTestData struct {
	expr string
	low  int
	high int
}

func toRTD(expr string, low, high int) *rollTestData {
	return &rollTestData{
		expr: expr,
		low:  low,
		high: high,
	}
}
