package catfacts

import (
	"fmt"
	"github.com/fabioxgn/go-bot"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	validResult = `{"facts": [
    	    "Catz FTW!"
    	],
    	"success": "true"}`

	emptyResult = `{"facts": [], "success": "true"}`
)

func TestCatFacts(t *testing.T) {
	apiResult := ""
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, apiResult)
		}))

	catFactsURL = ts.URL

	cmd := &bot.PassiveCmd{}

	Convey("Given a text", t, func() {

		Reset(func() {
			cmd.Raw = ""
			apiResult = ""
		})

		Convey("When the text does not have cat", func() {
			cmd.Raw = "My name is Catarina."
			s, err := catFacts(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})

		Convey("When the api returns 0 results", func() {
			apiResult = emptyResult
			cmd.Raw = "I love Catz!"

			s, err := catFacts(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})

		Convey("When the text has cat in the end of the sentence", func() {
			cmd.Raw = "I love Catz!"
			apiResult = validResult

			s, err := catFacts(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, fmt.Sprintf(msgPrefix, "Catz FTW!"))
		})

		Convey("When the text does not end with the world or puntuation", func() {
			cmd.Raw = "My name is Catzarina"

			s, err := catFacts(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})

		Convey("When the text has cat in the middle of a word", func() {
			cmd.Raw = "My name is aCats"

			s, err := catFacts(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})

		Convey("when the text have gato in the middle of the sentence", func() {
			cmd.Raw = "Eu tenho 2 gatos gordos."
			apiResult = validResult

			s, err := catFacts(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, fmt.Sprintf(msgPrefix, "Catz FTW!"))
		})

		Convey("When the api is unreachable", func() {
			cmd.Raw = "cat"
			catFactsURL = "127.0.0.1:0"

			_, err := catFacts(cmd)

			So(err, ShouldNotBeNil)
		})

	})
}
