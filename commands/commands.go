package commands

import (
	"fmt"
	"log"
)

type CommandFunc func(args []string) string

type IRC interface {
	Privmsg(string, string)
}

var (
	commands = make(map[string]CommandFunc)
)

func RegisterCommand(command string, f CommandFunc) {
	commands[command] = f
}

// HandleCmd handles a command and respond to channel or user
func HandleCmd(cmd *Command, channel string, irc IRC) {
	cmdFunction := commands[cmd.Command]
	if cmdFunction == nil {
		irc.Privmsg(channel, fmt.Sprintf("Command %v not found.", cmd.Command))
		printAvailableCommands(channel, irc)
	} else {
		log.Printf("cmd %v args %v", cmd.Command, cmd.Args)
		irc.Privmsg(channel, cmdFunction(cmd.Args))
	}
}

func printAvailableCommands(channel string, irc IRC) {
	irc.Privmsg(channel, "Available Commands:")
	cmds := ""
	for k := range commands {
		cmds += k + ", "
	}
	irc.Privmsg(channel, cmds[:len(cmds)-2])
}
