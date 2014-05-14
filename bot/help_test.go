package bot

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestHelp(t *testing.T) {
	Commands = make(map[string]*CustomCommand)
	Connection := ConnectionMock{}

	channel := ""
	msg := []string{}
	Connection.NoticeFunc = func(target, message string) {
		channel = target
		msg = append(msg, message)
	}

	command := &CustomCommand{
		Cmd:         "cmd",
		Description: "Command Description",
		Usage:       "Command Usage",
	}

	availableCommand := &Cmd{
		Nick:    "unavailable",
		Command: command.Cmd,
		Prefix:  "!",
	}
	unavailableCommand := &Cmd{
		Nick:    "nick",
		Command: "unavaible",
		Prefix:  "!",
	}
	RegisterCommand(command)

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
