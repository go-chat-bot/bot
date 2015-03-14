package puppet

import (
	"github.com/fabioxgn/go-bot"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPuppet(t *testing.T) {
	Convey("When say", t, func() {
		cmd := &bot.Cmd{}

		Convey("Should return usage if less than 3 arguments", func() {
			cmd.Args = []string{
				"say",
				"#go-bot",
			}
			cmd.Channel = "#channel"
			result, err := sendMessage(cmd)

			So(err, ShouldBeNil)
			So(result.Message, ShouldEqual, seeUsage)
			So(result.Channel, ShouldBeEmpty)
		})

		Convey("Should return error if the first argument is not say or act", func() {
			cmd.Args = []string{
				"hi",
				"#channel",
				"go-bot",
			}
			result, err := sendMessage(cmd)

			So(err, ShouldBeNil)
			So(result.Message, ShouldEqual, seeUsage)
			So(result.Channel, ShouldBeEmpty)
		})

		Convey("Should send a message to the specific channel", func() {

			cmd.Args = []string{
				"say",
				"#channel",
				"message with spaces",
			}
			result, err := sendMessage(cmd)

			So(err, ShouldBeNil)
			So(result.Channel, ShouldEqual, "#channel")
			So(result.Message, ShouldEqual, "message with spaces")
		})
	})
}
