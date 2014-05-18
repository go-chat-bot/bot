package catfacts

import (
	"github.com/fabioxgn/go-bot/web"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const (
	result = `{"facts": [
    	    "Catz FTW!"
    	],
    	"success": "true"}`

	emptyResult = `{"facts": [], "success": "true"}`
)

func TestCatFacts(t *testing.T) {
	url := ""
	setURL := func(u string) {
		url = u
	}
	Convey("Given a text", t, func() {
		Convey("When the text does not have cat", func() {

			s, err := getFacts("My name is Catarina.", nil)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})

		Convey("When the api returns 0 results", func() {
			s, err := getFacts("I love Catz!", web.GetJSONFromString(emptyResult, setURL))

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})

		Convey("When the text has cat in the end of the sentence", func() {

			s, err := getFacts("I love Catz!", web.GetJSONFromString(result, setURL))

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "Catz FTW!")
			So(url, ShouldEqual, catFactsURL)
		})

		Convey("when the text have gato in the middle of the sentence", func() {

			s, err := getFacts("Eu tenho 2 gatos gordos.", web.GetJSONFromString(result, setURL))

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "Catz FTW!")
		})

	})
}
