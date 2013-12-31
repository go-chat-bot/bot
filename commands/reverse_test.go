package commands

import (
	"testing"
)

func TestReverse(t *testing.T) {
	got := []string{
		"Hello",
		"World",
	}
	want := "dlroW olleH"
	s := Reverse(got)
	if s != want {
		t.Errorf("Expected %v got %v", want, s)
	}
}
