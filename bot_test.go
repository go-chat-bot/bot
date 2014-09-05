package bot

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const (
	serverName = "irc.server.com"
)

func TestGetServerName(t *testing.T) {
	Convey("Given a config message", t, func() {
		config = &Config{}
		Convey("When there is no port specified", func() {
			config.Server = serverName
			So(getServerName(), ShouldEqual, serverName)
		})
		Convey("When there is a port specified", func() {
			config.Server = serverName + ":6667"
			So(getServerName(), ShouldEqual, serverName)
		})
	})
}
