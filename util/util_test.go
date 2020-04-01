package util

import (
	"tablib/validate"
	"testing"
)

func TestBuildFullName_shouldBuildSimpleFullName(t *testing.T) {
	if BuildFullName("foo", "") != "foo" {
		t.Fail()
	}
}

func TestBuildFullName__shouldBuildSimpleFullNameCaseMatters(t *testing.T) {
	if BuildFullName("Foo", "") != "Foo" {
		t.Fail()
	}
}

func TestBuildFullName__shouldBuildCompoundFullName(t *testing.T) {
	if BuildFullName("foo", "7") != "foo.7" {
		t.Fail()
	}
}

func TestIsNotEmpty_shouldRejectOnEmptyString(t *testing.T) {
	vr := validate.NewValidationResult()
	IsNotEmpty("", "testval", "test", vr)
	if vr.IsValid {
		t.Fail()
	}
}

func TestIsNotEmpty_shouldAcceptOnNonEmptyString(t *testing.T) {
	vr := validate.NewValidationResult()
	IsNotEmpty("dlrow olleh", "testval", "test", vr)
	if !vr.IsValid {
		t.Fail()
	}
}

func TestIsValidIdentifier_shouldAcceptValidIds(t *testing.T) {

	var ids = []string{"Table_1_2020_kobe", "table_1_2020_kobe", "t1"}

	for _, id := range ids {
		vr := validate.NewValidationResult()
		IsValidIdentifier(id, "testval", "test", vr)
		if !vr.IsValid {
			t.Errorf("This should be a valid identifier: %s", id)
		}
	}
}

func TestIsValidIdentifier_shouldRejectInvalidIds(t *testing.T) {

	var ids = []string{"1Table", "", "_Table1", "a?Table", "A"}

	for _, id := range ids {
		vr := validate.NewValidationResult()
		IsValidIdentifier(id, "testval", "test", vr)
		if vr.IsValid {
			t.Errorf("This should be an invalid identifier: %s", id)
		}
	}
}
