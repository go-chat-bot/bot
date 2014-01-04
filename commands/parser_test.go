package commands

import (
	"fmt"
	"testing"
)

func TestEmptyCommand(t *testing.T) {
	cmd := Parse("")
	if cmd.Command != "" {
		t.Fail()
	}
}

func TestTrimArgs(t *testing.T) {
	command := "command"
	arg := "arg"
	c := Parse(fmt.Sprintf(" %v %v ", command, arg))

	if c.Command != command {
		t.Errorf("Expected %v got %v", command, c.Command)
	}

	if c.Args[0] != arg {
		t.Errorf("Expected %v got %v", arg, c.Args[0])
	}
}

func TestRawCommand(t *testing.T) {
	raw := "raw command"
	c := Parse(raw)

	if c.Raw != raw {
		t.Errorf("Expected %v, got %v", raw, c.Command)
	}
}

func TestCommandWithNoArgs(t *testing.T) {
	raw := "command"
	c := Parse(raw)

	if c.Command != raw {
		t.Errorf("Expected %v, got %v", raw, c.Command)
	}
}

func TestArgs(t *testing.T) {
	command := "command"
	expectedArgs := []string{
		"all",
		"your",
		"base",
	}

	c := Parse(fmt.Sprintf("%v %v %v %v", command, expectedArgs[0], expectedArgs[1], expectedArgs[2]))

	if len(c.Args) != len(expectedArgs) {
		t.Errorf("Expected %v args got %v", len(expectedArgs), len(c.Args))
	}

	for i := 0; i < len(expectedArgs); i++ {
		if c.Args[i] != expectedArgs[i] {
			t.Errorf("Expected %v got %v", expectedArgs[i], c.Args[i])
		}
	}
}

func TestCommandWithArgs(t *testing.T) {
	command := "command"
	args := " with 2 args"
	c := Parse(command + args)

	if c.Command != "command" {
		t.Errorf("Expected %v, got %v", command, c.Command)
	}
}
