package util

import (
	"fmt"
)

//BuildFullName builds the full name of the table or inline table
func BuildFullName(name, idnum string) string {
	if idnum == "" {
		return fmt.Sprintf("%s", name)
	}
	return fmt.Sprintf("%s.%s", name, idnum)
}
