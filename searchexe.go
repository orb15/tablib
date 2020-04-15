package tablib

import (
	"regexp"
	"sort"
)

type searchResultSet map[string]*SearchResult

func (cr *concreteTableRepo) executeSearch(namePredicate string, tags []string) ([]*SearchResult, error) {

	//case 1: nothing specified, return the entire repo
	if namePredicate == "" && len(tags) == 0 {
		return cr.fetchFullRepo(), nil
	}

	//case 2: tags only
	if namePredicate == "" && len(tags) > 0 {
		return cr.byTags(tags), nil
	}

	//case 3: name only
	if len(namePredicate) > 0 && len(tags) == 0 {
		sr := cr.fetchFullRepo()
		return cr.byName(namePredicate, sr)
	}

	//case 4: filter first by tags, then reduce this list by name
	sr := cr.byTags(tags)
	return cr.byName(namePredicate, sr)
}

func (cr *concreteTableRepo) byName(namePredicate string, sr []*SearchResult) ([]*SearchResult, error) {
	namePattern, err := regexp.Compile(namePredicate)
	if err != nil {
		return nil, err
	}

	srFinal := newSearchResultList(0)
	for _, i := range sr {
		if namePattern.MatchString(i.Name) {
			srFinal = append(srFinal, i)
		}
	}
	sortSearchResults(srFinal)
	return srFinal, nil
}

func (cr *concreteTableRepo) byTags(tags []string) []*SearchResult {
	//using map here to prevent duplicates as the same item can appear in more
	//than one tag entry in the tag cache
	rs := newSearchResultSet()
	for _, t := range tags {
		itemsWithTag := cr.tagSearchCache[t]
		for _, i := range itemsWithTag {
			rs[i.toComparable()] = i
		}
	}

	//convert map to sorted list
	sr := newSearchResultList(len(cr.nameSearchCache))
	for _, v := range rs {
		sr = append(sr, v)
	}
	sortSearchResults(sr)
	return sr
}

func (cr *concreteTableRepo) fetchFullRepo() []*SearchResult {
	sr := newSearchResultList(len(cr.nameSearchCache))
	for _, val := range cr.nameSearchCache {
		sr = append(sr, val)
	}
	sortSearchResults(sr)
	return sr
}

func sortSearchResults(sr []*SearchResult) {
	sort.Slice(sr, func(i, j int) bool {
		//scripts first
		if sr[i].Type == itemTypeScript && sr[j].Type == itemTypeTable {
			return true
		}
		if sr[i].Type == itemTypeTable && sr[j].Type == itemTypeScript {
			return false
		}
		//sort by name if both are scripts or tables
		return sr[i].Name < sr[j].Name
	})
}

func newSearchResultSet() searchResultSet {
	return make(map[string]*SearchResult)
}

func newSearchResultList(cap int) []*SearchResult {
	if cap <= 0 {
		cap = 0
	}
	return make([]*SearchResult, 0, cap)
}
