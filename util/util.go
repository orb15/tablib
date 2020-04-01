package util

import (
	"fmt"
	"regexp"
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
func IsNotEmpty(stringVal, yamlName, section string, vr *ValidationResult) {
	if stringVal == "" {
		vr.Fail(section, fmt.Sprintf("Empty %s", yamlName))
	}
}

//IsValidIdentifier validates the supplied string against the valid id regex
func IsValidIdentifier(stringVal, yamlName, section string, vr *ValidationResult) {
	if !validIdentifierPattern.MatchString(stringVal) {
		vr.Fail(section, fmt.Sprintf("Invalid identifier for %s: %s", yamlName, stringVal))
	}
}
