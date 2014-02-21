package commands

import (
	"testing"
)

func TestReverseString(t *testing.T) {
	arg := "Hello world"
	want := "dlrow olleH"
	cmd := &Command{
		Command: "reverse",
		FullArg: arg,
	}

	got, error := Reverse(cmd)

	if got != want {
		t.Errorf("Expected '%v' got '%v'", want, got)
	}

	if error != nil {
		t.Errorf("Expected '%v' got '%v'", nil, error)
	}
}

func TestReverseEmptyString(t *testing.T) {
	arg := ""
	want := ""
	cmd := &Command{
		Command: "reverse",
		FullArg: arg,
	}
	got, error := Reverse(cmd)

	if got != want {
		t.Errorf("Expected '%v' got '%v'", want, got)
	}

	if error != nil {
		t.Errorf("Expected '%v' got '%v'", nil, error)
	}
}
