package commands

import (
	"testing"
)

func TestReverseString(t *testing.T) {
	want := "dlrow olleH"
	got := Reverse([]string{"Hello", "world"})

	if got != want {
		t.Errorf("Expected '%v' got '%v'", want, got)
	}
}
