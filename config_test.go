package main

import (
	"strings"
	"testing"
)

func TestReadConfig(t *testing.T) {
	config := `{"Server": "irc.freenode.net:7000",
		"Channels": ["#lightirc", "#go-bot"],
		"User": "go-user",		
		"Nick": "go-bot",
		"Cmd": "!go-bot",
		"UseTLS": true}`

	c := &Config{}
	c.Read(strings.NewReader(config))

	server := "irc.freenode.net:7000"
	if c.Server != server {
		t.Errorf("Expected '%v' got '%v'", server, c.Server)
	}

	channels := []string{"#lightirc", "#go-bot"}
	if c.Channels[0] != channels[0] {
		t.Errorf("Expected '%v' got '%v'", channels[0], c.Channels[0])
	}
	if c.Channels[1] != channels[1] {
		t.Errorf("Expected '%v' got '%v'", channels[1], c.Channels[1])
	}

	user := "go-user"
	if c.User != user {
		t.Errorf("Expected '%v' got '%v'", user, c.User)
	}

	nick := "go-bot"
	if c.Nick != nick {
		t.Errorf("Expected '%v' got '%v'", nick, c.Nick)
	}

	cmd := "!go-bot"
	if c.Cmd != cmd {
		t.Errorf("Expected '%v' got '%v'", cmd, c.Cmd)
	}

	useTLS := true
	if c.UseTLS != useTLS {
		t.Errorf("Expected '%v' got '%v'", useTLS, c.UseTLS)
	}
}
