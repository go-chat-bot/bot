package bot

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type connMock struct {
	Channel  string
	Messages []string
	Nick     string
	channels []string
}

func (m *connMock) MessageReceived(target, message, sender string) {
	m.Channel = target
	m.Nick = sender
	m.Messages = append(m.Messages, message)
}

func (m *connMock) Channels() []string {
	return m.channels
}

func TestMessageReceived(t *testing.T) {
	Convey("Given a new message in the channel", t, func() {
		commands = make(map[string]*customCommand)
		conn := &connMock{}
		b := NewBot(conn.MessageReceived, conn.Channels())

		Convey("When the command is not registered", func() {
			Convey("It should not post to the channel", func() {

				b.MessageReceived("#go-bot", "!not_a_cmd", "user")

				So(conn.Messages, ShouldBeEmpty)
			})
		})

		Convey("The command can return an error", func() {
			Convey("it sould send the message with the error to the channel", func() {
				cmdError := errors.New("error")
				RegisterCommand("cmd", "", "",
					func(c *Cmd) (string, error) {
						return "", cmdError
					})

				b.MessageReceived("#go-bot", "!cmd", "user")

				So(conn.Channel, ShouldEqual, "#go-bot")
				So(conn.Messages, ShouldResemble,
					[]string{fmt.Sprintf(errorExecutingCommand, "cmd", cmdError.Error())})
			})
		})

		Convey("When the command is valid and registered", func() {
			commands = make(map[string]*customCommand)
			expectedMsg := "msg"
			cmd := "cmd"
			cmdDescription := "Command description"
			cmdExampleArgs := "arg1 arg2"

			RegisterCommand(cmd, cmdDescription, cmdExampleArgs,
				func(c *Cmd) (string, error) {
					return expectedMsg, nil
				})

			Convey("If it is called in the channel, reply on the channel", func() {
				b.MessageReceived("#go-bot", "!cmd", "user")

				So(conn.Channel, ShouldEqual, "#go-bot")
				So(conn.Messages, ShouldResemble, []string{expectedMsg})
			})

			Convey("If it is a private message, reply to the user", func() {
				conn.Nick = "go-bot"
				b.MessageReceived("go-bot", "!cmd", "sender-nick")

				So(conn.Nick, ShouldEqual, "sender-nick")
			})

			Convey("When the command is help", func() {
				Convey("Display the available commands in the channel", func() {
					b.MessageReceived("#go-bot", "!help", "user")

					So(conn.Channel, ShouldEqual, "#go-bot")
					So(conn.Messages, ShouldResemble, []string{
						fmt.Sprintf(helpAboutCommand, CmdPrefix),
						fmt.Sprintf(availableCommands, "cmd"),
					})
				})

				Convey("If the command exists send a message to the channel", func() {
					b.MessageReceived("#go-bot", "!help cmd", "user")

					So(conn.Channel, ShouldEqual, "#go-bot")
					So(conn.Messages, ShouldResemble, []string{
						fmt.Sprintf(helpDescripton, cmdDescription),
						fmt.Sprintf(helpUsage, CmdPrefix, cmd, cmdExampleArgs),
					})
				})

				Convey("If the command does not exists, display the generic help", func() {
					b.MessageReceived("#go-bot", "!help not_a_command", "user")

					So(conn.Channel, ShouldEqual, "#go-bot")
					So(conn.Messages, ShouldResemble, []string{
						fmt.Sprintf(helpAboutCommand, CmdPrefix),
						fmt.Sprintf(availableCommands, "cmd"),
					})
				})
			})
		})

		Convey("When the command is V2", func() {
			Convey("it should send the message with the error to the channel", func() {
				RegisterCommandV2("cmd", "", "",
					func(c *Cmd) (CmdResult, error) {
						return CmdResult{
							Channel: "#channel",
							Message: "message"}, nil
					})

				b.MessageReceived("#go-bot", "!cmd", "user")

				So(conn.Channel, ShouldEqual, "#channel")
				So(conn.Messages, ShouldResemble, []string{"message"})
			})

			Convey("it should reply to the current channel if the command does not specify one", func() {
				RegisterCommandV2("cmd", "", "",
					func(c *Cmd) (CmdResult, error) {
						return CmdResult{
							Message: "message"}, nil
					})

				b.MessageReceived("#go-bot", "!cmd", "user")

				So(conn.Channel, ShouldEqual, "#go-bot")
				So(conn.Messages, ShouldResemble, []string{"message"})
			})
		})

		Convey("When the command is passive", func() {
			passiveCommands = make(map[string]passiveCmdFunc)

			echo := func(cmd *PassiveCmd) (string, error) {
				return cmd.Raw, nil
			}
			ping := func(cmd *PassiveCmd) (string, error) {
				return "pong", nil
			}
			errored := func(cmd *PassiveCmd) (string, error) {
				return "", errors.New("error")
			}

			RegisterPassiveCommand("echo", echo)
			RegisterPassiveCommand("ping", ping)
			RegisterPassiveCommand("errored", errored)

			Convey("If it is called in the channel, reply on the channel", func() {
				b.MessageReceived("#go-bot", "test", "user")

				So(conn.Channel, ShouldEqual, "#go-bot")
				So(len(conn.Messages), ShouldEqual, 2)
				So(conn.Messages, ShouldContain, "test")
				So(conn.Messages, ShouldContain, "pong")
			})

			Convey("If it is a private message, reply to the user", func() {
				conn.Nick = "go-bot"
				b.MessageReceived("go-bot", "test", "sender-nick")

				So(conn.Nick, ShouldEqual, "sender-nick")
				So(len(conn.Messages), ShouldEqual, 2)
				So(conn.Messages, ShouldContain, "test")
				So(conn.Messages, ShouldContain, "pong")
			})
		})
	})
}
