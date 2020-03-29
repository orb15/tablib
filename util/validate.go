package util

import "fmt"

type ValidationResult struct {
	IsValid     bool
	HasWarnings bool
	Errors      []string
}

func (vr *ValidationResult) Fail(section, reason string) {
	vr.IsValid = false
	vr.Errors = append(vr.Errors, fmt.Sprintf("ERROR: %s - %s", section, reason))
}
func (vr *ValidationResult) Warn(section, reason string) {
	vr.IsValid = false
	vr.Errors = append(vr.Errors, fmt.Sprintf("Warn: %s - %s", section, reason))
}
