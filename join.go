package main

import (
	"github.com/fabioxgn/go-bot/cmd"
	"github.com/fabioxgn/go-bot/irc"
)

func Join(command cmd.Cmd, conn irc.Connection) {
	for _, channel := range command.Args {
		conn.Join(channel)
	}
}
