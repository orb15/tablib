package util

import (
	"fmt"
)

func BuildFullName(namespace, name, idnum string) string {
	if idnum == "" {
		return fmt.Sprintf("%s.%s", namespace, name)
	} else {
		return fmt.Sprintf("%s.%s.%s", namespace, name, idnum)
	}
}
