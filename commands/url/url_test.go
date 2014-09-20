package url

import (
	"fmt"
	"github.com/fabioxgn/go-bot"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestURL(t *testing.T) {
	cmd := &bot.PassiveCmd{}
	getExecuted := false
	getResult := ""

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			getExecuted = true
			fmt.Fprintln(w, getResult)
		}))

	url := ts.URL

	Convey("Given a text", t, func() {

		Reset(func() {
			getExecuted = false
			getResult = ""
		})

		Convey("If the text is not a URL", func() {
			cmd.Raw = "foo bar"
			title, err := urlTitle(cmd)

			So(getExecuted, ShouldBeFalse)
			So(err, ShouldBeNil)
			So(title, ShouldBeBlank)
		})

		Convey("If the url contains no title", func() {
			cmd.Raw = "foo " + url

			title, err := urlTitle(cmd)

			So(getExecuted, ShouldBeTrue)
			So(err, ShouldBeNil)
			So(title, ShouldBeBlank)
		})

		Convey("If the url contains a title", func() {
			getResult = "<title>Google</title>"
			cmd.Raw = fmt.Sprintf("foo %s bar", url)

			title, err := urlTitle(cmd)

			So(getExecuted, ShouldBeTrue)
			So(err, ShouldBeNil)
			So(title, ShouldEqual, "Google")
		})

		Convey("If the text is a https URL", func() {
			httpsURL := "https://google.com"

			extractedURL := extractURL(fmt.Sprintf("foo %s bar", httpsURL))

			So(extractedURL, ShouldEqual, httpsURL)
		})

		Convey("If title starts or ends with a new line", func() {
			getResult = "<title>\nGoogle\n</title>"
			cmd.Raw = url

			title, err := urlTitle(cmd)

			So(err, ShouldBeNil)
			So(title, ShouldEqual, "Google")
		})

		Convey("If an error occurs while fetching the url", func() {
			cmd.Raw = "127.0.0.1:0"

			_, err := urlTitle(cmd)

			So(err, ShouldNotBeNil)
		})

		Convey("if the url doesn't have a protocol", func() {
			noProtocolURL := "google.com"

			extractedURL := extractURL(fmt.Sprintf("foo %s bar", noProtocolURL))

			So(extractedURL, ShouldEqual, "http://google.com")
		})

		Convey("if the text has fewer than 4 characters", func() {
			So(extractURL("a.a"), ShouldEqual, "")
		})

		Convey("if the url is invalid", func() {
			So(extractURL(":googlecom"), ShouldEqual, "")
		})

	})
}
