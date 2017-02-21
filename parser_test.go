package bot

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	channel := &ChannelData{Channel: "#go-bot"}
	user := &User{Nick: "user123"}
	cmdWithoutArgs := CmdPrefix + "cmd"
	cmdWithArgs := CmdPrefix + "cmd    arg1  arg2   "
	cmdWithQuotes := CmdPrefix + "cmd    \"arg1  arg2\""

	tests := []struct {
		msg      string
		expected *Cmd
	}{
		{"", nil},
		{"!", nil},
		{"regular message", nil},
		{cmdWithoutArgs, &Cmd{
			Raw:         cmdWithoutArgs,
			Command:     "cmd",
			Channel:     channel.Channel,
			ChannelData: channel,
			User:        user,
			Message:     "cmd",
		}},
		{cmdWithArgs, &Cmd{
			Raw:         cmdWithArgs,
			Command:     "cmd",
			Channel:     channel.Channel,
			ChannelData: channel,
			User:        user,
			Message:     "cmd    arg1  arg2",
			RawArgs:     "arg1  arg2",
			Args:        []string{"arg1", "arg2"},
		}},
		{cmdWithQuotes, &Cmd{
			Raw:         cmdWithQuotes,
			Command:     "cmd",
			Channel:     channel.Channel,
			ChannelData: channel,
			User:        user,
			Message:     "cmd    \"arg1  arg2\"",
			RawArgs:     "\"arg1  arg2\"",
			Args:        []string{"arg1  arg2"},
		}},
	}

	for _, test := range tests {
		cmd, _ := parse(test.msg, channel, user)
		if !reflect.DeepEqual(test.expected, cmd) {
			t.Errorf("Expected:\n%#v\ngot:\n%#v", test.expected, cmd)
		}
	}
}

func TestInvalidArguments(t *testing.T) {
	cmd, err := parse("!cmd Invalid \"arg", &ChannelData{Channel: "#go-bot"}, &User{Nick: "user123"})
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if cmd != nil {
		t.Errorf("Expected nil, got %#v", cmd)
	}
}
