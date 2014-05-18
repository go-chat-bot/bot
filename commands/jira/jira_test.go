package jira

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestJira(t *testing.T) {
	url = "https://monde-sistemas.atlassian.net/browse/"

	Convey("Given a text", t, func() {
		Convey("When the text does not match a jira issue sintax", func() {

			s, err := getIssueURL("My name is go-bot, I am awesome.")

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})

		Convey("When the text match a jira issue sintax", func() {
			text := "My name is go-bot, I am awesome. MON-965"

			s, err := getIssueURL(text)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, fmt.Sprintf("%s%s", url, "MON-965"))
		})
	})
}
