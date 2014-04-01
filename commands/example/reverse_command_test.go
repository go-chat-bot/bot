package example

import (
	"github.com/fabioxgn/go-bot/cmd"
	"testing"
)

func TestReverseString(t *testing.T) {
	arg := "Hello world"
	want := "dlrow olleH"
	cmd := &cmd.Cmd{
		Command: "reverse",
		FullArg: arg,
	}

	got, error := reverse(cmd)

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
	cmd := &cmd.Cmd{
		Command: "reverse",
		FullArg: arg,
	}
	got, error := reverse(cmd)

	if got != want {
		t.Errorf("Expected '%v' got '%v'", want, got)
	}

	if error != nil {
		t.Errorf("Expected '%v' got '%v'", nil, error)
	}
}
