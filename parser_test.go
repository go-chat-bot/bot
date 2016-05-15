package bot

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestPaser(t *testing.T) {
	channel := "#go-bot"
	command := "command"
	rawArgs := "   arg1  arg2   "
	user := &User{Nick: "user123"}
	args := []string{
		"arg1",
		"arg2",
	}
	validCmd := fmt.Sprintf("%v%v", CmdPrefix, command)
	cmdWithArgs := fmt.Sprintf("%v%v %v", CmdPrefix, command, rawArgs)

	tests := []struct {
		msg      string
		expected *Cmd
	}{
		{"", nil},
		{"!", nil},
		{"regular message", nil},
		{validCmd, &Cmd{
			Raw:     validCmd,
			Command: command,
			Channel: channel,
			User:    user,
			Message: command,
		}},
		{cmdWithArgs, &Cmd{
			Raw:     cmdWithArgs,
			Command: command,
			Channel: channel,
			User:    user,
			Message: command + " " + strings.TrimRight(rawArgs, " "),
			RawArgs: strings.TrimSpace(rawArgs),
			Args:    args,
		}},
	}

	for _, test := range tests {
		cmd := parse(test.msg, channel, user)
		if !reflect.DeepEqual(test.expected, cmd) {
			t.Errorf("Expected:\n%#v\ngot:\n%#v", test.expected, cmd)
		}
	}
}
