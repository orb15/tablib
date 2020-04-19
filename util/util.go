package util

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"tablib/validate"

	"github.com/yuin/gopher-lua"
)

var (
	validIdentifierPattern = regexp.MustCompile("^[A-Za-z][a-zA-Z0-9_\\-]+$")
)

//BuildFullName builds the full name of the table or inline table
func BuildFullName(name, idnum string) string {
	if idnum == "" {
		return fmt.Sprintf("%s", name)
	}
	return fmt.Sprintf("%s.%s", name, idnum)
}

//IsNotEmpty validates that the incoming string is not ""
func IsNotEmpty(stringVal, yamlName, section string, vr *validate.ValidationResult) {
	if stringVal == "" {
		vr.Fail(section, fmt.Sprintf("Empty %s", yamlName))
	}
}

//IsValidIdentifier validates the supplied string against the valid id regex
func IsValidIdentifier(stringVal, yamlName, section string, vr *validate.ValidationResult) {
	if !validIdentifierPattern.MatchString(stringVal) {
		vr.Fail(section, fmt.Sprintf("Invalid identifier for %s: %s", yamlName, stringVal))
	}
}

//FindNextTableRef parses the string for '{.*}' and returns true if the string
//contains these characters and returns a slice that contains all characters prior
//to the first occurance in element 0, the actual reference in element 1 and all
//remaining characters in element 2
func FindNextTableRef(tableEntry string) ([]string, bool) {
	startIdx := strings.Index(tableEntry, "{")
	if startIdx == -1 {
		return nil, false
	}
	stopIdx := strings.Index(tableEntry, "}")
	asByteSlice := []byte(tableEntry)

	retval := make([]string, 3, 3)
	if startIdx == 0 {
		retval[0] = ""
	} else {
		retval[0] = string(asByteSlice[:startIdx])
	}
	retval[1] = string(asByteSlice[startIdx : stopIdx+1])
	retval[2] = string(asByteSlice[stopIdx+1:])
	return retval, true
}

//NewLuaState creates a new lua state/VM
//
//Set up a new lua VM. Limit the lua basic lib to essential functions in
//an attempt to reduce the scope of malicious scripts. This is actually really
//hard to do and the modules here are still condsidered dangerously unsafe but
//are neccessary if lua is to be used at all. Note that clever attackers
//can easily work around these limitations.
//
//See http://lua-users.org/wiki/SandBoxes for info on the relative futility of
//trying to make lua VMs both safe and functional
func NewLuaState() *lua.LState {
	//TODO: limit call stack and repository sizes - maybe?
	lState := lua.NewState(lua.Options{SkipOpenLibs: true})
	for _, pair := range []struct {
		n string
		f lua.LGFunction
	}{
		{lua.LoadLibName, lua.OpenPackage},
		{lua.BaseLibName, lua.OpenBase},
		{lua.TabLibName, lua.OpenTable},
		{lua.MathLibName, lua.OpenMath},
		{lua.StringLibName, lua.OpenString},
	} {
		if err := lState.CallByParam(lua.P{
			Fn:      lState.NewFunction(pair.f),
			NRet:    0,
			Protect: true,
		}, lua.LString(pair.n)); err != nil {
			fmt.Printf("Unable to fully establish Lua virtual machine: %s", err)
			os.Exit(-1)
		}
	}

	return lState
}
