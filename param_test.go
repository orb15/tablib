package tablib

import (
	"testing"
)

func TestParamSpecificationsFromMap_shouldWorkWithWellformedInput(t *testing.T) {
	mp := map[string]string{"key1": "val1-4|val1-1|val1-2|val1-3", "key2": "val2"}
	result := paramSpecificationsFromMap(mp)

	if len(result) != 2 {
		t.Errorf("Map data added or missing")
	}
	if result[0].Name == "key1" {
		if result[0].Default != "val1-4" {
			t.Errorf("Expect: %s Got: %s", "val1-4", result[0].Default)
		}
		if len(result[0].Options) != 4 {
			t.Errorf("Map Options added or missing")
		}
		if result[1].Name != "key2" {
			t.Errorf("Map data missing")
		}
		if result[1].Default != "val2" {
			t.Errorf("Expect: %s Got: %s", "val1-4", result[0].Default)
		}
		if len(result[1].Options) != 1 {
			t.Errorf("Map Options added or missing")
		}
	} else if result[0].Name == "key2" {
		if result[0].Default != "val2" {
			t.Errorf("Expect: %s Got: %s", "val1-4", result[0].Default)
		}
		if len(result[0].Options) != 1 {
			t.Errorf("Map Options added or missing")
		}
		if result[1].Name != "key1" {
			t.Errorf("Map data missing")
		}
		if result[1].Default != "val1-4" {
			t.Errorf("Expect: %s Got: %s", "val1-4", result[0].Default)
		}
		if len(result[1].Options) != 4 {
			t.Errorf("Map Options added or missing")
		}
	} else {
		t.Errorf("Map data not properly translated")
	}
}

func TestParamSpecificationsFromMap_shouldHandleBadInput1(t *testing.T) {
	mp := map[string]string{"key1": ""}
	result := paramSpecificationsFromMap(mp)

	if len(result) != 1 {
		t.Errorf("Map data added or missing")
	}
	if result[0].Name != "key1" {
		t.Errorf("Expected key1 but got: %s", result[0].Name)
	}
	if result[0].Default != "" {
		t.Errorf("Default not handled as expected")
	}
	if len(result[0].Options) != 1 {
		t.Errorf("Unexpected Options")
	}
	if result[0].Options[0] != "" {
		t.Errorf("Unexpected Option value")
	}
}

func TestParamSpecificationsFromMap_shouldHandleBadInput2(t *testing.T) {
	mp := map[string]string{"key1": "||"}
	result := paramSpecificationsFromMap(mp)

	if len(result) != 1 {
		t.Errorf("Map data added or missing")
	}
	if result[0].Name != "key1" {
		t.Errorf("Expected key1 but got: %s", result[0].Name)
	}
	if result[0].Default != "" {
		t.Errorf("Default not handled as expected")
	}
	if len(result[0].Options) != 3 {
		t.Errorf("Unexpected Options")
	}
	if result[0].Options[0] != "" {
		t.Errorf("Unexpected Option value")
	}
}

func TestDefaultParamSpecificationCallback_shouldSelectDefaults(t *testing.T) {
	pList := make([]*ParamSpecification, 2, 2)
	pList[0] = &ParamSpecification{Name: "foo", Default: "foo-1", Options: []string{"foo-1|foo-2"}}
	pList[1] = &ParamSpecification{Name: "bar", Default: "bar-1", Options: []string{"bar-1|bar2-2"}}
	result := DefaultParamSpecificationCallback(pList)

	if len(result) != 2 {
		t.Errorf("Map data added or missing")
	}
	r, found := result["foo"]
	if !found {
		t.Errorf("Map data added or missing")
	}
	if r != "foo-1" {
		t.Errorf("Expected foo-1 but got %s", r)
	}
	r, found = result["bar"]
	if !found {
		t.Errorf("Map data added or missing")
	}
	if r != "bar-1" {
		t.Errorf("Expected bar-1 but got %s", r)
	}
}
