package megasena

import (
	"github.com/fabioxgn/go-bot"
	. "github.com/smartystreets/goconvey/convey"
	"regexp"
	"testing"
)

func TestSortear(t *testing.T) {
    Convey("Sortear", t, func() {
    	So(sortear(6), ShouldEqual, "01 02 03 04 05 06")
    })
}

func TestMegaSena(t *testing.T) {
    Convey("Megasena", t, func() {
        
    	cmd := &bot.Cmd{
    		Command: "megasena",
    		Nick:    "nick",
    	}
	
	    Convey("Quando o argumento for gerar", func() {
	        cmd.Args = []string{"gerar"}
	        got, err := megasena(cmd)

	        So(err, ShouldBeNil)

	        match, err := regexp.MatchString("nick: (\\d{2} {1}){5}\\d{2}", got)
	        
	        So(err, ShouldBeNil)
	        So(match, ShouldBeTrue)
	    })
    })
}
