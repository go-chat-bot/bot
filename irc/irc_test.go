package irc

import (
	"testing"
)

const (
	serverName = "irc.server.com"
)

func TestGetServerName(t *testing.T) {
	if getServerName(serverName+":6667") != serverName {
		t.Fatal("Port should be removed from server name")
	}
}
