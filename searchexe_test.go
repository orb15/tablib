package tablib

import (
	"testing"
)

func TestExecuteSearch_shouldReturnFullRepoOnNoParams(t *testing.T) {
	yml1 := `
  definition:
    name: test
    type: flat
    note: this is an optional note
    tags:
      - tag1
  content:
    - item 1`

	yml2 := `
  definition:
    name: test2
    type: flat
    note: this is an optional note
    tags:
      - tag2
  content:
    - item 1`

	lua := `
	--TAGS: tag3
  print("dlrow olleh")
  `
	sryaml1 := &SearchResult{Name: "test", Type: itemTypeTable, Tags: []string{"tag1"}}
	sryaml2 := &SearchResult{Name: "test2", Type: itemTypeTable, Tags: []string{"tag2"}}
	srlua1 := &SearchResult{Name: "foo", Type: itemTypeScript, Tags: []string{"tag3"}}

	cr := newConcreteRepo()
	cr.AddTable([]byte(yml1))
	cr.AddTable([]byte(yml2))
	cr.AddLuaScript("foo", lua)
	sr, _ := cr.Search("", nil)
	if len(sr) != 3 {
		t.Errorf("Unexpected search result length")
	}
	res := compareSearchResults(sr, sryaml1, sryaml2, srlua1)
	for n, r := range res {
		if !r {
			t.Errorf("Item: %s failed", n)
		}
	}
}

func TestExecuteSearch_shouldReturnTagSearchWithoutDuplicates(t *testing.T) {
	yml1 := `
  definition:
    name: test
    type: flat
    note: this is an optional note
    tags:
      - tag1
      - tag2
  content:
    - item 1`

	yml2 := `
  definition:
    name: test2
    type: flat
    note: this is an optional note
    tags:
      - tag2
  content:
    - item 1`

	lua := `
	--TAGS: tag3
  print("dlrow olleh")
  `
	sryaml1 := &SearchResult{Name: "test", Type: itemTypeTable, Tags: []string{"tag1", "tag2"}}
	sryaml2 := &SearchResult{Name: "test2", Type: itemTypeTable, Tags: []string{"tag2"}}

	cr := newConcreteRepo()
	cr.AddTable([]byte(yml1))
	cr.AddTable([]byte(yml2))
	cr.AddLuaScript("foo", lua)
	sr, _ := cr.Search("", []string{"tag1", "tag2"})
	if len(sr) != 2 {
		t.Errorf("Unexpected search result length")
	}
	res := compareSearchResults(sr, sryaml1, sryaml2)
	for n, r := range res {
		if !r {
			t.Errorf("Item: %s failed", n)
		}
	}
}

func TestExecuteSearch_shouldReturnTagSearchWithoutNameConflicts(t *testing.T) {
	yml1 := `
  definition:
    name: test
    type: flat
    note: this is an optional note
    tags:
      - tag1
      - tag2
  content:
    - item 1`

	yml2 := `
  definition:
    name: test2
    type: flat
    note: this is an optional note
    tags:
      - tag2
  content:
    - item 1`

	lua := `
	--TAGS: tag1
  print("dlrow olleh")
  `
	sryaml1 := &SearchResult{Name: "test", Type: itemTypeTable, Tags: []string{"tag1", "tag2"}}
	srlua1 := &SearchResult{Name: "test", Type: itemTypeScript, Tags: []string{"tag1"}}

	cr := newConcreteRepo()
	cr.AddTable([]byte(yml1))
	cr.AddTable([]byte(yml2))
	cr.AddLuaScript("test", lua)
	sr, _ := cr.Search("", []string{"tag1"})
	if len(sr) != 2 {
		t.Errorf("Unexpected search result length")
	}
	res := compareSearchResults(sr, sryaml1, srlua1)
	for n, r := range res {
		if !r {
			t.Errorf("Item: %s failed", n)
		}
	}
}

func TestExecuteSearch_shouldSearchByNameOnly(t *testing.T) {
	yml1 := `
  definition:
    name: test
    type: flat
    note: this is an optional note
    tags:
      - tag1
  content:
    - item 1`

	yml2 := `
  definition:
    name: test2
    type: flat
    note: this is an optional note
    tags:
      - tag2
  content:
    - item 1`

	lua := `
	--TAGS: tag3
  print("dlrow olleh")
  `
	sryaml1 := &SearchResult{Name: "test", Type: itemTypeTable, Tags: []string{"tag1"}}
	sryaml2 := &SearchResult{Name: "test2", Type: itemTypeTable, Tags: []string{"tag2"}}

	cr := newConcreteRepo()
	cr.AddTable([]byte(yml1))
	cr.AddTable([]byte(yml2))
	cr.AddLuaScript("foo", lua)
	sr, _ := cr.Search("test.*", nil)
	if len(sr) != 2 {
		t.Errorf("Unexpected search result length")
	}
	res := compareSearchResults(sr, sryaml1, sryaml2)
	for n, r := range res {
		if !r {
			t.Errorf("Item: %s failed", n)
		}
	}
}

func TestExecuteSearch_shouldSearchByNameFailsOnBadRegex(t *testing.T) {
	yml1 := `
  definition:
    name: test
    type: flat
    note: this is an optional note
    tags:
      - tag1
  content:
    - item 1`

	yml2 := `
  definition:
    name: test2
    type: flat
    note: this is an optional note
    tags:
      - tag2
  content:
    - item 1`

	lua := `
	--TAGS: tag3
  print("dlrow olleh")
  `

	cr := newConcreteRepo()
	cr.AddTable([]byte(yml1))
	cr.AddTable([]byte(yml2))
	cr.AddLuaScript("foo", lua)
	_, err := cr.Search("(", nil)
	if err == nil {
		t.Error("Did not fail on bad regex as expected")
	}
}

func TestExecuteSearch_shouldSearchByTagsThenName(t *testing.T) {
	yml1 := `
  definition:
    name: test
    type: flat
    note: this is an optional note
    tags:
      - tag1
  content:
    - item 1`

	yml2 := `
  definition:
    name: test2
    type: flat
    note: this is an optional note
    tags:
      - tag2
  content:
    - item 1`

	yml3 := `
  definition:
    name: test3
    type: flat
    note: this is an optional note
    tags:
      - tag1
  content:
    - item 1`

	sryaml3 := &SearchResult{Name: "test3", Type: itemTypeTable, Tags: []string{"tag1"}}

	cr := newConcreteRepo()
	cr.AddTable([]byte(yml1))
	cr.AddTable([]byte(yml2))
	cr.AddTable([]byte(yml3))

	sr, _ := cr.Search("test[0-9]+", []string{"tag1"})
	if len(sr) != 1 {
		t.Errorf("Unexpected search result length")
	}
	res := compareSearchResults(sr, sryaml3)
	for n, r := range res {
		if !r {
			t.Errorf("Item: %s failed", n)
		}
	}
}

/* ***********************************************
* Test Helpers
* ***********************************************/

func compareSearchResults(sr []*SearchResult, srs ...*SearchResult) map[string]bool {
	sortSearchResults(srs)
	results := make(map[string]bool)
	for i, r := range srs {
		if srs[i].toFullComparable() == sr[i].toFullComparable() {
			results[r.toComparable()] = true
		} else {
			results[r.toComparable()] = false
		}
	}
	return results
}
