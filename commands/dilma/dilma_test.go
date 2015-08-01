package dilma

import (
	"github.com/fabioxgn/go-bot"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestDilma(t *testing.T) {
	Convey("Given a text", t, func() {
		cmd := &bot.PassiveCmd{}

		Convey("When the text does not match dilma", func() {
			cmd.Raw = "My name is go-bot, I am awesome."
			s, err := dilma(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})

		Convey("When the text match dilma", func() {
			cmd.Raw = "eu não votei na dilma!"

			s, err := dilma(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldNotEqual, "")
			So(strings.HasPrefix(s, ":dilma: "), ShouldBeTrue)
		})
	})
}
