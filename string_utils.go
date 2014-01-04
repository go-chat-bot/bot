package main

import (
	"strings"
)

// Returns the string after the separator
func StrAfter(s, sep string) string {
	if !strings.Contains(s, sep) {
		return ""
	}
	return s[strings.Index(s, sep)+len(sep):]
}
