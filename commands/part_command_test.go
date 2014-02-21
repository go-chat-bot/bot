package commands

import (
	"testing"
)

func TestPartNoArgs(t *testing.T) {
	cmd := &Command{
		Command: "part",
		Args: []string{},
	}
	got, error := Part(cmd)

	if error != nil {
		t.Errorf("Expected '%v' got '%v'", nil, error)
	}

	if got != "" {
		t.Errorf("Expected '%v' got '%v'", "", got)
	}
}

func TestPartWithArgs(t *testing.T) {
	cmd := &Command{
		Command: "part",
		Args: []string{"#channel1", "#channel2"},
	}
	got, error := Part(cmd)

	if error != nil {
		t.Errorf("Expected '%v' got '%v'", nil, error)
	}

	if got != "" {
		t.Errorf("Expected '%v' got '%v'", "", got)
	}
}
