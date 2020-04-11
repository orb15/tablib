package tablib

import "strings"

//ParamSpecification specifies a parameter a lua script requires
type ParamSpecification struct {
	Name    string   `json:"name"`
	Default string   `json:"default"`
	Options []string `json:"options"`
}

//NewParamSpecification does what it says on the tin
func NewParamSpecification() *ParamSpecification {
	return &ParamSpecification{
		Options: make([]string, 0),
	}
}

func paramSpecificationsFromMap(paramMap map[string]string) []*ParamSpecification {

	psList := make([]*ParamSpecification, 0, len(paramMap))
	for k, v := range paramMap {
		ps := NewParamSpecification()
		ps.Name = k
		parts := strings.Split(v, "|")
		defaultSet := false
		for _, p := range parts {
			if !defaultSet { //the first element is the default by (wait for it) default
				ps.Default = parts[0]
				defaultSet = true
			}
			ps.Options = append(ps.Options, p)
		}
		psList = append(psList, ps)
	}
	return psList
}

//ParamSpecificationRequestCallback is a function that will be called by tablib
//to allow the main program to supply params needed for a lua function
type ParamSpecificationRequestCallback = func([]*ParamSpecification) map[string]string

//DefaultParamSpecificationCallback is a convenience function that simply
//returns the default values in the provided map. This permits users of this
//lib from providing support for paramaterization of lua scripts at the cost of
//always having those scripts that do require params to always use their defaul
//value
func DefaultParamSpecificationCallback(specs []*ParamSpecification) map[string]string {
	mp := make(map[string]string)
	for _, s := range specs {
		mp[s.Name] = s.Default
	}
	return mp
}
