package jira

import (
	"fmt"
	"github.com/fabioxgn/go-bot"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestJira(t *testing.T) {
	url = "https://example.atlassian.net/browse/"
	Convey("Given a text", t, func() {
		cmd := &bot.PassiveCmd{}
		Convey("When the text does not match a jira issue syntax", func() {
			cmd.Raw = "My name is go-bot, I am awesome."
			s, err := jira(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})

		Convey("When the text match a jira issue syntax", func() {
			cmd.Raw = "My name is go-bot, I am awesome. MON-965"
			s, err := jira(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, fmt.Sprintf("%s%s", url, "MON-965"))
		})

		Convey("When the text has a jira issue in the midle of a word", func() {
			cmd.Raw = "My name is goBOT-123"
			s, err := jira(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})

		Convey("When the text has a jira issue syntax with only two numbers", func() {
			cmd.Raw = "BOT-12"
			s, err := jira(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, fmt.Sprintf("%s%s", url, "BOT-12"))
		})

		Convey("When the jira issue isn't preceeded by space", func() {
			cmd.Raw = "::BOT-122"
			s, err := jira(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})
	})
}
