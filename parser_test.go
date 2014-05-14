package bot

import (
	"fmt"
	"testing"
)

var (
	DefaultChannel = "#go-bot"
	DefaultNick    = "user123"
	DefaultCommand = "command"
	DefaultFullArg = "arg1 arg2"
	DefaultArgs    = []string{
		"arg1",
		"arg2",
	}
)

func TestEmptyCommand(t *testing.T) {
	cmd := Parse("", DefaultChannel, DefaultNick)
	if cmd.Command != "" {
		t.Fail()
	}
}

func TestWithoutPrefix(t *testing.T) {
	IsCommand := false
	Message := "regular message"

	res := Parse(Message, DefaultChannel, DefaultNick)

	if res.IsCommand != IsCommand {
		t.Errorf("Expected %v got %v", IsCommand, res.IsCommand)
	}
	if res.Message != Message {
		t.Errorf("Expected %v got %v", Message, res.Message)
	}
	if res.Channel != DefaultChannel {
		t.Errorf("Expected %v got %v", DefaultChannel, res.Channel)
	}
}

func TestOnlyPrefix(t *testing.T) {
	IsCommand := false

	res := Parse(CmdPrefix, DefaultChannel, DefaultNick)

	if res.IsCommand != IsCommand {
		t.Errorf("Expected %v got %v", IsCommand, res.IsCommand)
	}
	if res.Channel != DefaultChannel {
		t.Errorf("Expected %v got %v", DefaultChannel, res.Channel)
	}
}

func TestWithPrefixAndCommand(t *testing.T) {
	IsCommand := true
	cmd := fmt.Sprintf("%v%v", CmdPrefix, DefaultCommand)
	res := Parse(cmd, DefaultChannel, DefaultNick)

	if res.IsCommand != IsCommand {
		t.Errorf("Expected %v got %v", IsCommand, res.IsCommand)
	}
	if res.Command != DefaultCommand {
		t.Errorf("Expected %v got %v", DefaultCommand, res.Command)
	}
	if res.Channel != DefaultChannel {
		t.Errorf("Expected %v got %v", DefaultChannel, res.Channel)
	}
}

func TestWithPrefixAndCommandAndArgs(t *testing.T) {
	IsCommand := true
	cmd := fmt.Sprintf("%v%v %v", CmdPrefix, DefaultCommand, DefaultFullArg)
	res := Parse(cmd, DefaultChannel, DefaultNick)

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
	if res.Channel != DefaultChannel {
		t.Errorf("Expected %v got %v", DefaultChannel, res.Channel)
	}
}

func TestWithExtraSpaces(t *testing.T) {
	IsCommand := true
	cmd := fmt.Sprintf(" %v %v %v  %v  ", CmdPrefix, DefaultCommand, DefaultArgs[0], DefaultArgs[1])
	res := Parse(cmd, DefaultChannel, DefaultNick)

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
	if res.Channel != DefaultChannel {
		t.Errorf("Expected %v got %v", DefaultChannel, res.Channel)
	}
}
