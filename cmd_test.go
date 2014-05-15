package bot

import (
	"errors"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCmd(t *testing.T) {
	Convey("Given a command", t, func() {
		commands = make(map[string]*CustomCommand)
		conn := &ircConnectionMock{}

		Convey("When the command is not found", func() {
			cmd1 := &Cmd{
				Command: "cmd",
			}

			err := handleCmd(cmd1, nil)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, fmt.Sprintf(commandNotAvailable, "cmd"))
		})

		Convey("When the command exists", func() {
			cmd := &Cmd{
				Command: "cmd",
				Channel: "#go-bot",
			}

			Convey("it can return an error", func() {
				cmdError := errors.New("Error")
				RegisterCommand(&CustomCommand{
					Cmd:     cmd.Command,
					CmdFunc: func(c *Cmd) (string, error) { return "", cmdError },
				})

				msg := ""
				channel := ""
				conn.PrivMsgFunc = func(c string, m string) {
					channel = c
					msg = m
				}

				err := handleCmd(cmd, conn)

				So(err, ShouldEqual, cmdError)
				So(channel, ShouldEqual, cmd.Channel)
				So(msg, ShouldEqual, fmt.Sprintf(errorExecutingCommand, cmd.Command, cmdError.Error()))
			})

			Convey("it can return a string", func() {
				expectedMsg := "msg"
				RegisterCommand(&CustomCommand{
					Cmd:     cmd.Command,
					CmdFunc: func(c *Cmd) (string, error) { return expectedMsg, nil },
				})

				msg := ""
				channel := ""
				conn.PrivMsgFunc = func(c string, m string) {
					channel = c
					msg = m
				}

				err := handleCmd(cmd, conn)

				So(err, ShouldBeNil)
				So(channel, ShouldEqual, cmd.Channel)
				So(msg, ShouldEqual, expectedMsg)
			})

		})

	})

}
