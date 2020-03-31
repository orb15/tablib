package util

import (
	"testing"
)

func TestTable_shouldBuildSimpleFullName(t *testing.T) {
	if BuildFullName("foo", "") != "foo" {
		t.Fail()
	}
}

func TestTable_shouldBuildSimpleFullNameCaseMatters(t *testing.T) {
	if BuildFullName("Foo", "") != "Foo" {
		t.Fail()
	}
}

func TestTable_shouldBuildCompoundFullName(t *testing.T) {
	if BuildFullName("foo", "7") != "foo.7" {
		t.Fail()
	}
}
