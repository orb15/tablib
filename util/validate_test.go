package util

import (
	"testing"
)

func TestFail_shouldBuildFail(t *testing.T) {
	vr := NewValidationResult()
	vr.Fail("TestSection", "TestReason")

	if vr.Valid() {
		t.Error("IsValid")
	}
	if vr.IssueCount() != 1 {
		t.Error("Wrong Issue Count")
	}
	if vr.ErrorCount() != 1 {
		t.Error("Wrong Errors Count")
	}
	if vr.WarnCount() != 0 {
		t.Error("Wrong Warn Count")
	}
	if vr.HasWarnings {
		t.Error("HasWarnings")
	}
	if vr.Errors[0] != "ERROR: TestSection - TestReason" {
		t.Error("Bad string")
	}
}

func TestFail_shouldBuildWarn(t *testing.T) {
	vr := NewValidationResult()
	vr.Warn("TestSection", "TestReason")

	if !vr.Valid() {
		t.Error("IsValid")
	}
	if vr.IssueCount() != 1 {
		t.Error("Wrong Issue Count")
	}
	if vr.ErrorCount() != 0 {
		t.Error("Wrong Errors Count")
	}
	if vr.WarnCount() != 1 {
		t.Error("Wrong Warn Count")
	}
	if !vr.HasWarnings {
		t.Error("HasWarnings")
	}
	if vr.Errors[0] != "WARN: TestSection - TestReason" {
		t.Error("Bad string")
	}
}
