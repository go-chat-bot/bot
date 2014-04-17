package main

import (
	"github.com/fabioxgn/go-bot/cmd"
	"github.com/fabioxgn/go-bot/irc"
)

func Part(command *cmd.Cmd, conn irc.Connection) (msg string, err error) {
	if len(command.Args) > 0 {
		for _, channel := range command.Args {
			conn.Part(channel)
		}
	} else {
		conn.Part(command.Channel)
	}
	return
}
