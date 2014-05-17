package bot

import (
	"errors"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMessageReceived(t *testing.T) {
	Convey("Given a new message in the channel", t, func() {
		commands = make(map[string]*CustomCommand)
		conn := &ircConnectionMock{}

		Convey("When the command is not found", func() {
			cmdFuncCalled := false
			RegisterCommand(&CustomCommand{
				Cmd: "cmd",
				CmdFunc: func(c *Cmd) (string, error) {
					cmdFuncCalled = true
					return "", nil
				},
			})

			messageReceived("#go-bot", "!not_a_cmd", "user", conn)

			So(cmdFuncCalled, ShouldBeFalse)
		})

		Convey("The command can return an error", func() {
			cmdError := errors.New("Error")
			RegisterCommand(&CustomCommand{
				Cmd:     "cmd",
				CmdFunc: func(c *Cmd) (string, error) { return "", cmdError },
			})

			msg := ""
			channel := ""
			conn.PrivMsgFunc = func(c string, m string) {
				channel = c
				msg = m
			}

			messageReceived("#go-bot", "!cmd", "user", conn)

			So(channel, ShouldEqual, "#go-bot")
			So(msg, ShouldEqual, fmt.Sprintf(errorExecutingCommand, "cmd", cmdError.Error()))
		})

		Convey("The command can return a string", func() {
			expectedMsg := "msg"
			RegisterCommand(&CustomCommand{
				Cmd:     "cmd",
				CmdFunc: func(c *Cmd) (string, error) { return expectedMsg, nil },
			})

			msg := ""
			channel := ""
			conn.PrivMsgFunc = func(c string, m string) {
				channel = c
				msg = m
			}

			messageReceived("#go-bot", "!cmd", "user", conn)

			So(channel, ShouldEqual, "#go-bot")
			So(msg, ShouldEqual, expectedMsg)
		})

		Convey("The command can be a private message", func() {
			RegisterCommand(&CustomCommand{
				Cmd: "cmd",
				CmdFunc: func(c *Cmd) (string, error) {
					return "hi", nil
				},
			})

			channel := ""
			conn.PrivMsgFunc = func(c string, m string) {
				channel = c
			}

			conn.Nick = "go-bot"
			messageReceived("go-bot", "!cmd", "sender-nick", conn)

			So(channel, ShouldEqual, "sender-nick")
		})

	})

}
