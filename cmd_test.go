package bot

import (
	"errors"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type ircConnectionMock struct {
	Channel  string
	Messages []string
	Nick     string
	Joined   string
	Parted   string
}

func (m *ircConnectionMock) Privmsg(target, message string) {
	m.Channel = target
	m.Messages = append(m.Messages, message)
}

func (m ircConnectionMock) GetNick() string {
	return m.Nick
}

func (m *ircConnectionMock) Join(target string) {
	m.Joined = target
}

func (m *ircConnectionMock) Part(target string) {
	m.Parted = target
}

func TestMessageReceived(t *testing.T) {
	Convey("Given a new message in the channel", t, func() {
		commands = make(map[string]*customCommand)
		conn := &ircConnectionMock{}

		Convey("When the command is join", func() {

			Reset(func() {
				conn.Messages = []string{}
			})

			Convey("if the channel is not specified", func() {
				messageReceived("#go-bot", "!join    ", "user", conn)

				So(conn.Messages, ShouldResemble, []string{joinUsage})
			})

			Convey("if the channel is specified", func() {
				messageReceived("#go-bot", "!join #channel pass", "user", conn)

				So(conn.Joined, ShouldEqual, "#channel pass")
				So(conn.Channel, ShouldEqual, "#channel")
				So(conn.Messages, ShouldResemble,
					[]string{fmt.Sprintf(joinMessage, "user")})
			})
		})

		Convey("When the command is part", func() {
			config = &Config{
				Channels: []string{"#go-bot", "#safechan passwd", ""},
			}

			Reset(func() {
				conn.Parted = ""
				conn.Messages = []string{}
			})

			Convey("it should part the channel", func() {
				messageReceived("#mychannel", "!part    ", "user", conn)

				So(conn.Parted, ShouldEqual, "#mychannel")
				So(conn.Messages, ShouldResemble, []string{partMessage})
			})

			Convey("if the channel is in the config", func() {
				messageReceived("#Go-Bot", "!part    ", "user", conn)
				So(conn.Parted, ShouldEqual, "")
				So(conn.Messages, ShouldResemble, []string{partNotAllowed})
			})

			Convey("if the channel is in the config and has a password", func() {
				messageReceived("#safechan", "!part", "user", conn)
				So(conn.Parted, ShouldEqual, "")
				So(conn.Messages, ShouldResemble, []string{partNotAllowed})
			})
		})

		Convey("When the command is not registered", func() {
			conn = &ircConnectionMock{}

			Convey("It should not post to the channel", func() {

				messageReceived("#go-bot", "!not_a_cmd", "user", conn)

				So(conn.Messages, ShouldBeEmpty)
			})

		})

		Convey("The command can return an error", func() {
			conn = &ircConnectionMock{}

			Convey("it sould send the message with the error to the channel", func() {
				cmdError := errors.New("error")
				RegisterCommand("cmd", "", "",
					func(c *Cmd) (string, error) {
						return "", cmdError
					})

				messageReceived("#go-bot", "!cmd", "user", conn)

				So(conn.Channel, ShouldEqual, "#go-bot")
				So(conn.Messages, ShouldResemble,
					[]string{fmt.Sprintf(errorExecutingCommand, "cmd", cmdError.Error())})
			})
		})

		Convey("When the command is valid and registered", func() {
			conn = &ircConnectionMock{}

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
				messageReceived("#go-bot", "!cmd", "user", conn)

				So(conn.Channel, ShouldEqual, "#go-bot")
				So(conn.Messages, ShouldResemble, []string{expectedMsg})
			})

			Convey("If it is a private message, reply to the user", func() {
				conn.Nick = "go-bot"
				messageReceived("go-bot", "!cmd", "sender-nick", conn)

				So(conn.Channel, ShouldEqual, "sender-nick")
			})

			Convey("When the command is help", func() {

				Convey("Display the available commands in the channel", func() {
					messageReceived("#go-bot", "!help", "user", conn)

					So(conn.Channel, ShouldEqual, "#go-bot")
					So(conn.Messages, ShouldResemble, []string{
						fmt.Sprintf(helpAboutCommand, CmdPrefix),
						fmt.Sprintf(availableCommands, "cmd"),
					})
				})

				Convey("If the command exists send a message to the channel", func() {

					messageReceived("#go-bot", "!help cmd", "user", conn)

					So(conn.Channel, ShouldEqual, "#go-bot")
					So(conn.Messages, ShouldResemble, []string{
						fmt.Sprintf(helpDescripton, cmdDescription),
						fmt.Sprintf(helpUsage, CmdPrefix, cmd, cmdExampleArgs),
					})

				})

				Convey("If the command does not exists, display the generic help", func() {
					messageReceived("#go-bot", "!help not_a_command", "user", conn)

					So(conn.Channel, ShouldEqual, "#go-bot")
					So(conn.Messages, ShouldResemble, []string{
						fmt.Sprintf(helpAboutCommand, CmdPrefix),
						fmt.Sprintf(availableCommands, "cmd"),
					})
				})
			})

		})

		Convey("When the command is V2", func() {
			conn = &ircConnectionMock{}

			Convey("it should send the message with the error to the channel", func() {
				RegisterCommandV2("cmd", "", "",
					func(c *Cmd) (CmdResult, error) {
						return CmdResult{
							Channel: "#channel",
							Message: "message"}, nil
					})

				messageReceived("#go-bot", "!cmd", "user", conn)

				So(conn.Channel, ShouldEqual, "#channel")
				So(conn.Messages, ShouldResemble, []string{"message"})
			})

			Convey("it should reply to the current channel if the command does not specify one", func() {
				RegisterCommandV2("cmd", "", "",
					func(c *Cmd) (CmdResult, error) {
						return CmdResult{
							Message: "message"}, nil
					})

				messageReceived("#go-bot", "!cmd", "user", conn)

				So(conn.Channel, ShouldEqual, "#go-bot")
				So(conn.Messages, ShouldResemble, []string{"message"})
			})
		})

		Convey("When the command is passive", func() {
			conn = &ircConnectionMock{}

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
				messageReceived("#go-bot", "test", "user", conn)

				So(conn.Channel, ShouldEqual, "#go-bot")
				So(len(conn.Messages), ShouldEqual, 2)
				So(conn.Messages, ShouldContain, "test")
				So(conn.Messages, ShouldContain, "pong")
			})

			Convey("If it is a private message, reply to the user", func() {
				conn.Nick = "go-bot"
				messageReceived("go-bot", "test", "sender-nick", conn)

				So(conn.Channel, ShouldEqual, "sender-nick")
				So(len(conn.Messages), ShouldEqual, 2)
				So(conn.Messages, ShouldContain, "test")
				So(conn.Messages, ShouldContain, "pong")
			})
		})
	})

}
