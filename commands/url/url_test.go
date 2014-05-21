package url

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestURL(t *testing.T) {
	Convey("Give a text", t, func() {
		getExecuted := false
		getResult := []byte{}
		get := func(string) ([]byte, error) {
			getExecuted = true
			return getResult, nil
		}

		Convey("If the text is not a URL", func() {
			getExecuted = false
			title, err := getTitle("foo bar", get)

			So(getExecuted, ShouldBeFalse)
			So(err, ShouldBeNil)
			So(title, ShouldBeBlank)
		})

		Convey("If the text contains a http URL", func() {

			Convey("If the url contains no title", func() {
				getExecuted = false
				title, err := getTitle("foo http://google.com bar", get)

				So(getExecuted, ShouldBeTrue)
				So(err, ShouldBeNil)
				So(title, ShouldBeBlank)
			})

			Convey("If the url contains a title", func() {
				getExecuted = false
				getResult = []byte("<title>Google</title>")
				title, err := getTitle("foo http://google.com bar", get)

				So(getExecuted, ShouldBeTrue)
				So(err, ShouldBeNil)
				So(title, ShouldEqual, "Google")
			})
		})

		Convey("If the text is a https URL", func() {
			getExecuted = false
			getResult = []byte("<title>Google</title>")
			title, err := getTitle("foo https://google.com bar", get)

			So(getExecuted, ShouldBeTrue)
			So(err, ShouldBeNil)
			So(title, ShouldEqual, "Google")
		})

		Convey("If title contains a new line", func() {
			getExecuted = false
			getResult = []byte("<title>Google\n</title>")
			title, err := getTitle("https://google.com", get)

			So(getExecuted, ShouldBeTrue)
			So(err, ShouldBeNil)
			So(title, ShouldEqual, "Google")
		})

		Convey("if the url doesn't have a protocol", func() {
			getExecuted = false
			getResult = []byte("<title>Google</title>")
			title, err := getTitle("foo google.com bar", get)

			So(getExecuted, ShouldBeTrue)
			So(err, ShouldBeNil)
			So(title, ShouldEqual, "Google")
		})

		Convey("if the text has fewer than 4 characters", func() {
			getExecuted = false
			_, _ = getTitle("a.a", get)

			So(getExecuted, ShouldBeFalse)
		})

	})
}
