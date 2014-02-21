package commands

import (
	"testing"
)

func TestJoin(t *testing.T) {
	cmd := &Command{
		Command: "join",
		Args: []string{"#channel1", "#channel2"},
	}
	got, error := Join(cmd)

	if error != nil {
		t.Errorf("Expected '%v' got '%v'", nil, error)
	}

	if got != "" {
		t.Errorf("Expected '%v' got '%v'", "", got)
	}
}
