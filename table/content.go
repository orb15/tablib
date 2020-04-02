package table

import (
	"regexp"
)

var (
	//InlineCalledPattern represents syntax for an inline table call
	InlineCalledPattern = regexp.MustCompile("(\\{#[0-9]+\\})")
	//ExternalCalledPattern represents syntax for an external table call
	ExternalCalledPattern = regexp.MustCompile("\\{@(.*)\\}")
	//PickCalledPattern represents syntax for a pick table call
	PickCalledPattern = regexp.MustCompile("\\{[0-9]+!(.*)\\}")
)
