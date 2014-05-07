package main

import (
	"fmt"
	"github.com/fabioxgn/go-bot/cmd"
	"github.com/fabioxgn/go-bot/irc"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestHelp(t *testing.T) {
	cmd.Commands = make(map[string]*cmd.CustomCommand)
	Connection := irc.ConnectionMock{}

	channel := ""
	msg := []string{}
	Connection.NoticeFunc = func(target, message string) {
		channel = target
		msg = append(msg, message)
	}

	command := &cmd.CustomCommand{
		Cmd:         "cmd",
		Description: "Command Description",
		Usage:       "Command Usage",
	}

	availableCommand := &cmd.Cmd{
		Nick:    "unavailable",
		Command: command.Cmd,
		Prefix:  "!",
	}
	unavailableCommand := &cmd.Cmd{
		Nick:    "nick",
		Command: "unavaible",
		Prefix:  "!",
	}
	cmd.RegisterCommand(command)

	Convey("Given a help command", t, func() {
		msg = []string{}

		Convey("when the command is not registered", func() {
			Help(unavailableCommand, Connection)
			Convey("should send a notice to the user with the available commands", func() {
				So(channel, ShouldEqual, unavailableCommand.Nick)
				So(msg, ShouldResemble, []string{
					fmt.Sprintf(helpAboutCommand, unavailableCommand.Prefix),
					fmt.Sprintf(availableCommands, availableCommand.Command),
				})
			})

		})

		Convey("when the command is registered", func() {
			Help(availableCommand, Connection)
			Convey("should send a notice to the user with the command's Description and Usage", func() {
				So(channel, ShouldEqual, availableCommand.Nick)
				So(msg, ShouldResemble, []string{
					fmt.Sprintf(helpDescripton, command.Description),
					fmt.Sprintf(helpUsage, availableCommand.Prefix, command.Cmd, command.Usage),
				})
			})
		})

	})

}
