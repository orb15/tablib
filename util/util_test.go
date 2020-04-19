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
	var ids = []string{"Table_1_2020_kobe", "table_1_2020_kobe", "t1", "table-1"}

	for _, id := range ids {
		vr := validate.NewValidationResult()
		IsValidIdentifier(id, "testval", "test", vr)
		if !vr.IsValid {
			t.Errorf("This should be a valid identifier: %s", id)
		}
	}
}

func TestIsValidIdentifier_shouldRejectInvalidIds(t *testing.T) {
	var ids = []string{"1Table", "", "_Table1", "a,Table", "A", "-table1", "tab le"}

	for _, id := range ids {
		vr := validate.NewValidationResult()
		IsValidIdentifier(id, "testval", "test", vr)
		if vr.IsValid {
			t.Errorf("This should be an invalid identifier: %s", id)
		}
	}
}
func TestFindNextTableRef_shouldWork(t *testing.T) {
	ttrs := []*TestTableRef{
		toTTR("", false, "", "", ""),
		toTTR("hello world", false, "", "", ""),
		toTTR("hello {@world}", true, "hello ", "{@world}", ""),
		toTTR("{@hello} {@world}", true, "", "{@hello}", " {@world}"),
		toTTR("{@hello}", true, "", "{@hello}", ""),
		toTTR("Foo {@Arg} not {@Error} ", true, "Foo ", "{@Arg}", " not {@Error} "),
	}

	for _, ttr := range ttrs {
		s, b := FindNextTableRef(ttr.tableString)
		if b != ttr.b {
			t.Errorf("Bool value mismatch wanted: %t got: %t", ttr.b, b)
		}
		if ttr.b {
			if s[0] != ttr.s0 {
				t.Errorf("Prefix mismatch wanted: %s got: %s", ttr.s0, s[0])
			}
			if s[1] != ttr.s1 {
				t.Errorf("Reference mismatch wanted: %s got: %s", ttr.s1, s[1])
			}
			if s[2] != ttr.s2 {
				t.Errorf("Suffix mismatch wanted: %s got: %s", ttr.s2, s[2])
			}
		}
	}
}

type TestTableRef struct {
	tableString string
	b           bool
	s0          string
	s1          string
	s2          string
}

func toTTR(tableString string, expectedBool bool, expected0, expected1, expected2 string) *TestTableRef {
	return &TestTableRef{
		tableString: tableString,
		b:           expectedBool,
		s0:          expected0,
		s1:          expected1,
		s2:          expected2,
	}
}
