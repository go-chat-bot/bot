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

//TODO: Register help for built in commands
// func init() {
// 	RegisterCommand("part", Part)

// 	man := Manual{
// 		helpDescripton: "Leave from the specified channels",
// 		helpUse:        "#channel1 [#channel2 ... ]",
// 	}
// 	RegisterHelp("part", man)
// }
