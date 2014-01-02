package main

import (
	"strings"
)

func StrAfter(s, sep string) string {
	if !strings.Contains(s, sep) {
		return ""
	}
	return s[strings.Index(s, sep)+len(sep):]
}
