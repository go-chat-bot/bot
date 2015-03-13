package puppet

import (
	"github.com/fabioxgn/go-bot"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCotacao(t *testing.T) {

	Convey("When say", t, func() {
		cmd := &bot.Cmd{}

		Convey("Should return error if less than 2 arguments", func() {
			cmd.Args = []string{
				"say",
			}
			_, err := sendMessage(cmd)

			So(err, ShouldNotBeNil)
		})

		Convey("Should return error if the first argument is not say or me", func() {
			cmd.Args = []string{
				"hi",
				"#channel",
				"go-bot",
			}
			_, err := sendMessage(cmd)

			So(err, ShouldNotBeNil)
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
