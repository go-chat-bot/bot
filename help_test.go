package bot

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestHelp(t *testing.T) {
	commands = make(map[string]*CustomCommand)
	Connection := ConnectionMock{}

	channel := ""
	msg := []string{}
	Connection.PrivMsgFunc = func(target, message string) {
		channel = target
		msg = append(msg, message)
	}

	command := &CustomCommand{
		Cmd:         "cmd",
		Description: "Command Description",
		Usage:       "Command Usage",
	}

	availableCommand := &Cmd{
		Channel: "unavailable",
		Command: command.Cmd,
	}
	unavailableCommand := &Cmd{
		Channel: "nick",
		Command: "unavaible",
	}
	RegisterCommand(command)

	Convey("Given a help command", t, func() {
		msg = []string{}

		Convey("when the command is not registered", func() {
			Help(unavailableCommand, Connection)
			Convey("should send a message to the channel with the available commands", func() {
				So(channel, ShouldEqual, unavailableCommand.Channel)
				So(msg, ShouldResemble, []string{
					fmt.Sprintf(helpAboutCommand, CmdPrefix),
					fmt.Sprintf(availableCommands, availableCommand.Command),
				})
			})

		})

		Convey("when the command is registered", func() {
			Help(availableCommand, Connection)
			Convey("should send a message to the channel with the command's Description and Usage", func() {
				So(channel, ShouldEqual, availableCommand.Channel)
				So(msg, ShouldResemble, []string{
					fmt.Sprintf(helpDescripton, command.Description),
					fmt.Sprintf(helpUsage, CmdPrefix, command.Cmd, command.Usage),
				})
			})
		})

	})

}
