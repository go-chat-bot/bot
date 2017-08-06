package bot

import (
	"strings"
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
			MessageData: &Message{Text: strings.TrimLeft("cmd", CmdPrefix)},
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
			MessageData: &Message{Text: strings.TrimLeft("cmd    arg1  arg2", CmdPrefix)},
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
			MessageData: &Message{Text: strings.TrimLeft("cmd    \"arg1  arg2\"", CmdPrefix)},
		}},
	}

	for _, test := range tests {
		cmd, _ := parse(test.msg, channel, user)
		if test.expected != nil && cmd != nil {
			if test.expected.Raw != cmd.Raw {
				t.Errorf("Expected Raw:\n%#v\ngot:\n%#v", test.expected.Raw, cmd.Raw)
			}
			if test.expected.Channel != cmd.Channel {
				t.Errorf("Expected Channel:\n%#v\ngot:\n%#v", test.expected.Channel, cmd.Channel)
			}
			if test.expected.Message != cmd.Message {
				t.Errorf("Expected Message:\n%#v\ngot:\n%#v", test.expected.Message, cmd.Message)
			}
			if test.expected.Command != cmd.Command {
				t.Errorf("Expected Command:\n%#v\ngot:\n%#v", test.expected.Command, cmd.Command)
			}
			if test.expected.RawArgs != cmd.RawArgs {
				t.Errorf("Expected RawArgs:\n%#v\ngot:\n%#v", test.expected.RawArgs, cmd.RawArgs)
			}
			if test.expected.ChannelData.Protocol != cmd.ChannelData.Protocol {
				t.Errorf("Expected ChannelData.Protocol:\n%#v\ngot:\n%#v", test.expected.ChannelData.Protocol, cmd.ChannelData.Protocol)
			}
			if test.expected.ChannelData.Server != cmd.ChannelData.Server {
				t.Errorf("Expected ChannelData.Server:\n%#v\ngot:\n%#v", test.expected.ChannelData.Server, cmd.ChannelData.Server)
			}
			if test.expected.ChannelData.Channel != cmd.ChannelData.Channel {
				t.Errorf("Expected ChannelData.Channel:\n%#v\ngot:\n%#v", test.expected.ChannelData.Channel, cmd.ChannelData.Channel)
			}
			if test.expected.ChannelData.IsPrivate != cmd.ChannelData.IsPrivate {
				t.Errorf("Expected ChannelData.IsPrivate:\n%#v\ngot:\n%#v", test.expected.ChannelData.IsPrivate, cmd.ChannelData.IsPrivate)
			}
			if test.expected.User.ID != cmd.User.ID {
				t.Errorf("Expected User.ID:\n%#v\ngot:\n%#v", test.expected.User.ID, cmd.User.ID)
			}
			if test.expected.User.Nick != cmd.User.Nick {
				t.Errorf("Expected User.Nick:\n%#v\ngot:\n%#v", test.expected.User.Nick, cmd.User.Nick)
			}
			if test.expected.User.RealName != cmd.User.RealName {
				t.Errorf("Expected User.RealName:\n%#v\ngot:\n%#v", test.expected.User.RealName, cmd.User.RealName)
			}
			if test.expected.User.IsBot != cmd.User.IsBot {
				t.Errorf("Expected User.IsBot:\n%#v\ngot:\n%#v", test.expected.User.IsBot, cmd.User.IsBot)
			}
			if test.expected.MessageData.Text != cmd.MessageData.Text {
				t.Errorf("Expected MessageData.Text:\n%#v\ngot:\n%#v", test.expected.MessageData.Text, cmd.MessageData.Text)
			}
			if test.expected.MessageData.IsAction != cmd.MessageData.IsAction {
				t.Errorf("Expected MessageData.IsAction:\n%#v\ngot:\n%#v", test.expected.MessageData.IsAction, cmd.MessageData.IsAction)
			}
			for i, arg := range test.expected.Args {
				if arg != cmd.Args[i] {
					t.Errorf("Expected cmd.Args[]:\n%#v\ngot:\n%#v", arg, cmd.Args[i])
				}
			}
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
