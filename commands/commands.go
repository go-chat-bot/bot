package commands

import (
	"fmt"
	"log"
)

// CommandFunc is the function to be executed when a user calls a command
// the arguments passed to the command will be passed to the function as a string slice
// the arguments are separated by spaces so: command arg1 arg2 will pass a slice with 2 itens
// arg1 and arg2 to the function
type CommandFunc func(args []string) string
type privMsgFunc func(channel string, msg string)

const (
	commandNotAvailable = "Command %v not available."
	availableCommands   = "Available commands"
	noCommandsAvailable = "No commands available."
)

var (
	commands = make(map[string]CommandFunc)
)

// RegisterCommand must be used to register a command (string) and CommandFunc.
// The CommandFunc will be executed when a users calls the bot passing the
// command string as the first argument
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
	cmds := ""
	for k := range commands {
		cmds += k + ", "
	}
	if cmds != "" {
		Msg(channel, fmt.Sprintf("%s: %s", availableCommands, cmds[:len(cmds)-2]))
	} else {
		Msg(channel, noCommandsAvailable)
	}
}
