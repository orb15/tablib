package validate

import (
	"fmt"
	"strings"
)

//ValidationResult holds table validation info
type ValidationResult struct {
	IsValid     bool
	HasWarnings bool
	Errors      []string
}

//NewValidationResult does what is says on the tin
func NewValidationResult() *ValidationResult {
	vr := &ValidationResult{
		IsValid:     true,
		HasWarnings: false,
		Errors:      make([]string, 0),
	}
	return vr
}

//Fail indicates a validation failure
func (vr *ValidationResult) Fail(section, reason string) {
	vr.IsValid = false
	vr.Errors = append(vr.Errors, fmt.Sprintf("ERROR: %s - %s", section, reason))
}

//Warn indicates a validation warning
func (vr *ValidationResult) Warn(section, reason string) {
	vr.HasWarnings = true
	vr.Errors = append(vr.Errors, fmt.Sprintf("WARN: %s - %s", section, reason))
}

//Valid returns true if table is valid (no errors)
func (vr *ValidationResult) Valid() bool {
	return vr.IsValid
}

//IssueCount provides the number of errors or warnings
func (vr *ValidationResult) IssueCount() int {
	return len(vr.Errors)
}

//ErrorCount provides the number of errors
func (vr *ValidationResult) ErrorCount() int {
	count := 0
	for _, e := range vr.Errors {
		if strings.HasPrefix(e, "ERROR") {
			count++
		}
	}
	return count
}

//WarnCount provides the number of warnings
func (vr *ValidationResult) WarnCount() int {
	count := 0
	for _, e := range vr.Errors {
		if strings.HasPrefix(e, "WARN") {
			count++
		}
	}
	return count
}
