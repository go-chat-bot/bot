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

func TestSeparatorNotFound(t *testing.T) {
	sep := "!go-bot"
	want := ""
	got := StrAfter("all your base are belong to us", sep)
	if got != want {
		t.Errorf("Expected '%v' got '%v'", want, got)
	}
}

func TestEmptyString(t *testing.T) {
	sep := "!go-bot"
	want := ""
	got := StrAfter("", sep)
	if got != want {
		t.Errorf("Expected '%v' got '%v'", want, got)
	}
}
