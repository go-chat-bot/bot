package commands

import (
	"testing"
)

func TestQuitNoArgs(t *testing.T) {
	cmd := &Command{
		Command: "quit",
		Args: []string{},
	}
	got, error := Quit(cmd)

	if error != nil {
		t.Errorf("Expected '%v' got '%v'", nil, error)
	}

	if got != "" {
		t.Errorf("Expected '%v' got '%v'", "", got)
	}
}

func TestQuitWithArgs(t *testing.T) {
	cmd := &Command{
		Command: "quit",
		Args: []string{"arg1", "arg2"},
	}
	got, error := Quit(cmd)

	if error != nil {
		t.Errorf("Expected '%v' got '%v'", nil, error)
	}

	if got != "" {
		t.Errorf("Expected '%v' got '%v'", "", got)
	}
}
