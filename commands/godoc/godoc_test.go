package godoc

import (
	"fmt"
	"github.com/fabioxgn/go-bot/web"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const (
	results = `{
    	"results": [
	        {
	            "path": "github.com/fabioxgn/go-bot",
	            "synopsis": "IRC bot written in go"
	        }
    	]
	}`

	emptyResults = `{"results":[]}`
)

func TestGoDoc(t *testing.T) {
	Convey("Given a search query text", t, func() {
		url := ""
		setURL := func(u string) {
			url = u
		}

		Convey("When the result is empty", func() {
			s, err := searchGodoc("non existant package", web.GetJSONFromString(emptyResults, setURL))

			So(err, ShouldBeNil)
			So(s, ShouldEqual, noPackagesFound)
			So(url, ShouldEqual, fmt.Sprintf(godocSearchURL, "non existant package"))
		})

		Convey("When the result is ok", func() {
			s, err := searchGodoc("go-bot", web.GetJSONFromString(results, setURL))

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "github.com/fabioxgn/go-bot IRC bot written in go")
			So(url, ShouldEqual, fmt.Sprintf(godocSearchURL, "go-bot"))
		})

		Convey("When the query is empty", func() {
			url = ""
			s, err := searchGodoc("", web.GetJSONFromString(results, setURL))

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
			So(url, ShouldEqual, "")
		})
	})

}
