package util

import (
	"testing"
)

func TestValidateDiceExpr_shouldAcceptValidDice(t *testing.T) {

	var ids = []string{"1d8", "3d6", "4d2001", "3d6 + 4d8", "1d6 - 7d3",
		"1d6 + 2", "3d6 * 100", "1d4 * 2d6 + 3"}

	for _, id := range ids {
		vr := NewValidationResult()
		ValidateDiceExpr(id, "testval", vr)
		if !vr.IsValid {
			t.Errorf("This should be a valid die expression: %s", id)
		}
	}
}

func TestValidateDiceExpr_shouldRejectInvalidDice(t *testing.T) {

	var ids = []string{"", "0d6", "1d0", "3d6 / 4d8", "d6", "7d", "3d6 +",
		"3d6 2d6", "3d6 + 3 + 2d8", "2 + 1d6", "1d6+8", "3d0 + 3"}

	for _, id := range ids {
		vr := NewValidationResult()
		ValidateDiceExpr(id, "testval", vr)
		if vr.IsValid {
			t.Errorf("This should be an invalid die expression: %s", id)
		}
	}
}

func TestValidateDiceExpr_shouldParseDiceExpr1(t *testing.T) {
	vr := NewValidationResult()
	pde := ValidateDiceExpr("4d6", "testval", vr)

	if len(pde) != 1 {
		t.Error("Bad parse length")
	}
	if pde[0].Count != 4 {
		t.Error("Bad count")
	}
	if pde[0].DieType != 6 {
		t.Error("Bad die type")
	}
	if pde[0].Operator != "none" {
		t.Error("Bad operator")
	}
}

func TestValidateDiceExpr_shouldParseDiceExpr2(t *testing.T) {
	vr := NewValidationResult()
	pde := ValidateDiceExpr("4d6 + 3d7 - 1d3 * 21", "testval", vr)

	if len(pde) != 4 {
		t.Error("Bad parse length")
	}
	if pde[0].Count != 4 {
		t.Error("Bad count")
	}
	if pde[0].DieType != 6 {
		t.Error("Bad die type")
	}
	if pde[0].Operator != "+" {
		t.Error("Bad operator")
	}
	if pde[1].Count != 3 {
		t.Error("Bad count")
	}
	if pde[1].DieType != 7 {
		t.Error("Bad die type")
	}
	if pde[1].Operator != "-" {
		t.Error("Bad operator")
	}
	if pde[2].Count != 1 {
		t.Error("Bad count")
	}
	if pde[2].DieType != 3 {
		t.Error("Bad die type")
	}
	if pde[2].Operator != "*" {
		t.Error("Bad operator")
	}
	if pde[3].Count != 21 {
		t.Error("Bad count")
	}
	if pde[3].DieType != 0 {
		t.Error("Bad die type")
	}
	if pde[3].Operator != "none" {
		t.Error("Bad operator")
	}
}
