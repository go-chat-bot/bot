package irc

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	serverName = "irc.server.com"
)

func TestGetServerName(t *testing.T) {
	Convey("Given server name", t, func() {
		Convey("When there is no port specified", func() {
			So(getServerName(serverName), ShouldEqual, serverName)
		})

		Convey("When there is a port specified", func() {
			So(getServerName(serverName+":6667"), ShouldEqual, serverName)
		})
	})
}
