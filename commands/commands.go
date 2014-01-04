package commands

import (
	"fmt"
	"log"
)

type CommandFunc func(args []string) string
type privMsgFunc func(channel string, msg string)

type IRC interface {
	Privmsg(string, string)
}

const (
	commandNotAvailable = "Command %v not available."
)

var (
	commands = make(map[string]CommandFunc)
)

func RegisterCommand(command string, f CommandFunc) {
	commands[command] = f
}

// HandleCmd handles a command and respond to channel or user
func HandleCmd(cmd *Command, channel string, Msg privMsgFunc) {
	cmdFunction := commands[cmd.Command]
	if cmdFunction == nil {
		Msg(channel, fmt.Sprintf(commandNotAvailable, cmd.Command))
		printAvailableCommands(channel, Msg)
	} else {
		log.Printf("cmd %v args %v", cmd.Command, cmd.Args)
		Msg(channel, cmdFunction(cmd.Args))
	}
}

func printAvailableCommands(channel string, Msg privMsgFunc) {
	availableCommands := "Available Commands: "
	for k := range commands {
		availableCommands += k + ", "
	}
	Msg(channel, availableCommands[:len(availableCommands)-2])
}
