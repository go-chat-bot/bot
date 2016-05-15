package bot

import (
	"reflect"
	"testing"
)

func TestPaser(t *testing.T) {
	channel := "#go-bot"
	user := &User{Nick: "user123"}
	cmdWithoutArgs := CmdPrefix + "cmd"
	cmdWithArgs := CmdPrefix + "cmd    arg1  arg2   "

	tests := []struct {
		msg      string
		expected *Cmd
	}{
		{"", nil},
		{"!", nil},
		{"regular message", nil},
		{cmdWithoutArgs, &Cmd{
			Raw:     cmdWithoutArgs,
			Command: "cmd",
			Channel: channel,
			User:    user,
			Message: "cmd",
		}},
		{cmdWithArgs, &Cmd{
			Raw:     cmdWithArgs,
			Command: "cmd",
			Channel: channel,
			User:    user,
			Message: "cmd    arg1  arg2",
			RawArgs: "arg1  arg2",
			Args:    []string{"arg1", "arg2"},
		}},
	}

	for _, test := range tests {
		cmd := parse(test.msg, channel, user)
		if !reflect.DeepEqual(test.expected, cmd) {
			t.Errorf("Expected:\n%#v\ngot:\n%#v", test.expected, cmd)
		}
	}
}
