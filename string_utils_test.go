package main

import (
	"testing"
)

func TestString(t *testing.T) {
	sep := "!go-bot"
	want := " reverse"
	got := StrAfter(sep+want, sep)
	if got != want {
		t.Errorf("Want %v got %v", want, got)
	}
}
