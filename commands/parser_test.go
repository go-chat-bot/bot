package commands

import (
	"fmt"
	"testing"
)

var (
	DefaultPrefix = "!"
	DefaultCommand = "command"
	DefaultFullArg = "arg1 arg2"
	DefaultArgs = []string{
		"arg1",
		"arg2",
	}
)

func TestEmptyCommand(t *testing.T) {
	cmd := Parse("", DefaultPrefix)
	if cmd.Command != "" {
		t.Fail()
	}
}

func TestWithoutPrefix(t *testing.T) {
	IsCommand := false
	Message := "regular message"

	res := Parse(Message, DefaultPrefix)

	if res.IsCommand != IsCommand {
		t.Errorf("Expected %v got %v", IsCommand, res.IsCommand)
	}
	if res.Message != Message {
		t.Errorf("Expected %v got %v", Message, res.Message)
	}
}

func TestOnlyPrefix(t *testing.T) {
	IsCommand := false

	res := Parse(DefaultPrefix, DefaultPrefix)

	if res.IsCommand != IsCommand {
		t.Errorf("Expected %v got %v", IsCommand, res.IsCommand)
	}
}

func TestWithPrefixAndCommand(t *testing.T) {
	IsCommand := true

	res := Parse(fmt.Sprintf("%v%v", DefaultPrefix, DefaultCommand), DefaultPrefix)

	if res.IsCommand != IsCommand {
		t.Errorf("Expected %v got %v", IsCommand, res.IsCommand)
	}
	if res.Command != DefaultCommand {
		t.Errorf("Expected %v got %v", DefaultCommand, res.Command)
	}
}

func TestWithPrefixAndCommandAndArgs(t *testing.T) {
	IsCommand := true

	res := Parse(fmt.Sprintf("%v%v %v", DefaultPrefix, DefaultCommand, DefaultFullArg), DefaultPrefix)

	if res.IsCommand != IsCommand {
		t.Errorf("Expected %v got %v", IsCommand, res.IsCommand)
	}
	if res.Command != DefaultCommand {
		t.Errorf("Expected %v got %v", DefaultCommand, res.Command)
	}
	if res.Args[0] != DefaultArgs[0] {
		t.Errorf("Expected %v got %v", DefaultArgs[0], res.Args[0])
	}
	if res.FullArg != DefaultFullArg {
		t.Errorf("Expected %v got %v", DefaultFullArg, res.FullArg)
	}
}

func TestWithExtraSpaces(t *testing.T) {
	IsCommand := true

	res := Parse(fmt.Sprintf(" %v %v %v  %v  ", DefaultPrefix, DefaultCommand, DefaultArgs[0], DefaultArgs[1]), DefaultPrefix)

	if res.IsCommand != IsCommand {
		t.Errorf("Expected %v got %v", IsCommand, res.IsCommand)
	}
	if res.Command != DefaultCommand {
		t.Errorf("Expected %v got %v", DefaultCommand, res.Command)
	}

	for i := 0; i < len(DefaultArgs); i++ {
		if res.Args[i] != DefaultArgs[i] {
			t.Errorf("Expected %v got %v", DefaultArgs[i], res.Args[i])
		}
	}

	if res.FullArg != DefaultFullArg {
		t.Errorf("Expected %v got %v", DefaultFullArg, res.FullArg)
	}
}