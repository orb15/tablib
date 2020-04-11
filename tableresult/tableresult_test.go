package tableresult

import (
	"testing"
)

func TestAddLog_shouldAddToLog(t *testing.T) {
	tr := NewTableResult()
	tr.AddLog("1")
	tr.AddLog("two")

	if len(tr.Log) != 2 {
		t.Fail()
	}
	if tr.Log[0] != "1" {
		t.Fail()
	}
	if tr.Log[1] != "two" {
		t.Fail()
	}
}

func TestAddLog_shouldAddToResult(t *testing.T) {
	tr := NewTableResult()
	tr.AddResult("1")
	tr.AddResult("two")

	if len(tr.Result) != 2 {
		t.Fail()
	}
	if tr.Result[0] != "1" {
		t.Fail()
	}
	if tr.Result[1] != "two" {
		t.Fail()
	}
}
