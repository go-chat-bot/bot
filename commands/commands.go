package commands

import (
	// "fmt"
   // "errors"
	"log"
	"github.com/thoj/go-ircevent"
)

// type privMsgFunc func(channel string, msg string)

// CommandFunc is the function to be executed when a user calls a command
type CommandFunc func(cmd *Command) (string, error)
type Manual struct {
	helpDescripton 	string
	helpUse 				string
}

const (
	commandNotAvailable = "Command %v not available."
	availableCommands   = "Available commands"
	noCommandsAvailable = "No commands available."
)

var (
	helps = make(map[string]Manual)
	commands = make(map[string]CommandFunc)
	irccon *irc.Connection
)

// RegisterCommand must be used to register a command (string) and CommandFunc.
// The CommandFunc will be executed when a users calls the bot passing the
// command string as the first argument
func RegisterCommand(command string, f CommandFunc) {
	commands[command] = f
}

func RegisterHelp(command string, h Manual) {
	helps[command] = h
}

// HandleCmd handles a command
func HandleCmd(cmd *Command, irc *irc.Connection) (err error) {
	cmdFunction := commands[cmd.Command]
	irccon = irc
	if cmdFunction == nil {
		// TODO: create error
		return err
	} else {
		log.Printf("cmd %v args %v", cmd.Command, cmd.FullArg)
		resultStr, resultError := cmdFunction(cmd)
		if resultError != nil {
			// TODO: create error
			return err
			irc.Privmsg(cmd.Channel, "Show the fucking error")
		} else if resultStr != "" {
			irc.Privmsg(cmd.Channel, resultStr)
		}
	}
	return
}
