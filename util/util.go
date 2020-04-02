package util

import (
	"fmt"
	"regexp"
	"strings"
	"tablib/validate"
)

var (
	validIdentifierPattern = regexp.MustCompile("^[A-Za-z][a-zA-Z0-9_]+$")
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
	retval[1] = string(asByteSlice[startIdx+1 : stopIdx])
	retval[2] = string(asByteSlice[stopIdx+1:])
	return retval, true
}
