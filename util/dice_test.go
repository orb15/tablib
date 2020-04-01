package util

import (
	"testing"
)

func TestIsValidIdentifier_shouldAcceptValidDice(t *testing.T) {

	var ids = []string{"1d8", "3d6", "4d2001", "3d6 + 4d8", "1d6 - 7d3"}

	for _, id := range ids {
		vr := NewValidationResult()
		ValidateDiceExpr(id, "testval", vr)
		if !vr.IsValid {
			t.Errorf("This should be a valid die expression: %s", id)
		}
	}
}

func TestIsValidIdentifier_shouldRejectInvalidDice(t *testing.T) {

	var ids = []string{"", "0d6", "1d0", "3d6 * 4d8", "d6", "7d", "3d6 +"}

	for _, id := range ids {
		vr := NewValidationResult()
		ValidateDiceExpr(id, "testval", vr)
		if vr.IsValid {
			t.Errorf("This should be an invalid die expression: %s", id)
		}
	}
}
