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
  print("dlrow olleh")
  `
	sryaml1 := &SearchResult{Name: "test", Type: itemTypeTable, Tags: []string{"tag1"}}
	sryaml2 := &SearchResult{Name: "test2", Type: itemTypeTable, Tags: []string{"tag2"}}
	srlua1 := &SearchResult{Name: "foo", Type: itemTypeScript, Tags: []string{"tag3"}}

	cr := newConcreteRepo()
	cr.AddTable([]byte(yml1))
	cr.AddTable([]byte(yml2))
	cr.AddLuaScript("foo", lua, []string{"tag3"})
	sr, _ := cr.Search("", nil)
	res := compareSearchResults(sr, sryaml1, sryaml2, srlua1)
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
