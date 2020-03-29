package util

import "fmt"

//ValidationResult holds table validation info
type ValidationResult struct {
	IsValid     bool
	HasWarnings bool
	Errors      []string
}

//Fail indicates a validation failure
func (vr *ValidationResult) Fail(section, reason string) {
	vr.IsValid = false
	vr.Errors = append(vr.Errors, fmt.Sprintf("ERROR: %s - %s", section, reason))
}

//Warn indicates a validation warning
func (vr *ValidationResult) Warn(section, reason string) {
	vr.IsValid = false
	vr.Errors = append(vr.Errors, fmt.Sprintf("WARN: %s - %s", section, reason))
}
