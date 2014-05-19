package jira

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestJira(t *testing.T) {
	url = "https://monde-sistemas.atlassian.net/browse/"

	Convey("Given a text", t, func() {
		Convey("When the text does not match a jira issue syntax", func() {

			s, err := getIssueURL("My name is go-bot, I am awesome.")

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})

		Convey("When the text match a jira issue syntax", func() {
			text := "My name is go-bot, I am awesome. MON-965"

			s, err := getIssueURL(text)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, fmt.Sprintf("%s%s", url, "MON-965"))
		})

		Convey("When the text has a jira issue in the midle of a word", func() {
			text := "My name is goBOT-123"

			s, err := getIssueURL(text)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})

		Convey("When the text has a jira issue syntax with only two numbers", func() {
			text := "BOT-12"

			s, err := getIssueURL(text)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})
	})
}
