package main

import (
	"strings"
)

func StrAfter(s, sep string) string {
	return s[strings.Index(s, sep)+len(sep):]
}
