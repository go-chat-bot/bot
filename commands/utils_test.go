package commands

import (
	"testing"
)

func TestArgsToString(t *testing.T) {
	want := "all your base are belong to us"

	got := ArgsToString([]string{"all", "your", "base", "are", "belong", "to", "us"})

	if got != want {
		t.Errorf("Expected '%v' got '%v'", want, got)
	}
}
