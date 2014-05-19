package chucknorris

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestChuckNorris(t *testing.T) {
	Convey("Given a text", t, func() {
		Convey("When the text does not match a chuck norris name", func() {

			s, err := getChuckNorrisFact("My name is go-bot, I am awesome.")

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})

		Convey("When the text match a chuck name", func() {
			text := "My name is chuck"

			s, err := getChuckNorrisFact(text)

			So(err, ShouldBeNil)
			So(s, ShouldNotEqual, "")
		})
	})
}
