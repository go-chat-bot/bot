package commands

import (
	"testing"
)

func TestHelloworld(t *testing.T) {
	want := "Hello world!"
	cmd := &Command{
		Command: "helloworld",
		FullArg: want,
	}
	got, error := Helloworld(cmd)

	if got != want {
		t.Errorf("Expected '%v' got '%v'", want, got)
	}

	if error != nil {
		t.Errorf("Expected '%v' got '%v'", nil, error)
	}
}
