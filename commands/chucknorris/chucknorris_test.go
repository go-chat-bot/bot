package chucknorris

import (
	"github.com/fabioxgn/go-bot"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestChuckNorris(t *testing.T) {
	Convey("Given a text", t, func() {
		cmd := &bot.PassiveCmd{}
		Convey("When the text does not match a chuck norris name", func() {
			cmd.Raw = "My name is go-bot, I am awesome."
			s, err := chucknorris(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})

		Convey("When the text match a chuck name", func() {
			cmd.Raw = "My name is chuck"

			s, err := chucknorris(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldNotEqual, "")
		})

		Convey("When the text match norris", func() {
			cmd.Raw = "Hi, I'm Mr. Norris"

			s, err := chucknorris(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldNotEqual, "")
		})
	})
}
