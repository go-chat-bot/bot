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
			s, err := getIssueURL("My name is go-bot, I am awesome. MON-965")

			So(err, ShouldBeNil)
			So(s, ShouldEqual, fmt.Sprintf("%s%s", url, "MON-965"))
		})

		Convey("When the text has a jira issue in the midle of a word", func() {
			s, err := getIssueURL("My name is goBOT-123")

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})

		Convey("When the text has a jira issue syntax with only two numbers", func() {
			s, err := getIssueURL("BOT-12")

			So(err, ShouldBeNil)
			So(s, ShouldEqual, fmt.Sprintf("%s%s", url, "BOT-12"))
		})

		Convey("When the jira issue isn't preceeded by space", func() {
			s, err := getIssueURL("::BOT-122")

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})
	})

}
