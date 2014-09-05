package bot

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	DefaultChannel = "#go-bot"
	DefaultNick    = "user123"
	DefaultCommand = "command"
	DefaultFullArg = "arg1 arg2"
	DefaultArgs    = []string{
		"arg1",
		"arg2",
	}
)

func TestPaser(t *testing.T) {
	Convey("When given a message", t, func() {
		Convey("When the message is empty", func() {
			cmd := parse("", DefaultChannel, DefaultNick)

			So(cmd, ShouldBeNil)
		})

		Convey("When the message doesn't have the prefix", func() {
			Message := "regular message"
			cmd := parse(Message, DefaultChannel, DefaultNick)

			So(cmd, ShouldBeNil)
		})

		Convey("When the message is only the prefix", func() {
			cmd := parse(CmdPrefix, DefaultChannel, DefaultNick)

			So(cmd, ShouldBeNil)
		})

		Convey("When the message is valid command", func() {
			msg := fmt.Sprintf("%v%v", CmdPrefix, DefaultCommand)
			cmd := parse(msg, DefaultChannel, DefaultNick)

			So(cmd, ShouldNotBeNil)
			So(cmd.Command, ShouldEqual, DefaultCommand)
			So(cmd.Channel, ShouldEqual, DefaultChannel)
		})

		Convey("When the message is a command with args", func() {
			msg := fmt.Sprintf("%v%v %v", CmdPrefix, DefaultCommand, DefaultFullArg)
			cmd := parse(msg, DefaultChannel, DefaultNick)

			So(cmd, ShouldNotBeNil)
			So(cmd.Command, ShouldEqual, DefaultCommand)
			So(cmd.Channel, ShouldEqual, DefaultChannel)
			So(cmd.Args, ShouldResemble, DefaultArgs)
			So(cmd.FullArg, ShouldEqual, DefaultFullArg)
		})

		Convey("When the message has extra spaces", func() {
			msg := fmt.Sprintf(" %v %v %v  %v  ", CmdPrefix, DefaultCommand, DefaultArgs[0], DefaultArgs[1])
			cmd := parse(msg, DefaultChannel, DefaultNick)

			So(cmd, ShouldNotBeNil)
			So(cmd.Command, ShouldEqual, DefaultCommand)
			So(cmd.Channel, ShouldEqual, DefaultChannel)
			So(cmd.Args, ShouldResemble, DefaultArgs)
			So(cmd.FullArg, ShouldEqual, DefaultFullArg)
		})
	})
}
