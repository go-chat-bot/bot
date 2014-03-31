package commands

import (
	"errors"
	"fmt"
	"log"
)

type IRCConnection interface {
	Privmsg(target, message string)
	Notice(target, message string)
	Join(channel string)
	Part(channel string)
	Quit()
}

// type privMsgFunc func(channel string, msg string)

// CommandFunc is the function to be executed when a user calls a command
type CommandFunc func(cmd *Command) (string, error)
type Manual struct {
	helpDescripton string
	helpUse        string
}

const (
	commandNotAvailable   = "Command %v not available."
	availableCommands     = "Available commands"
	noCommandsAvailable   = "No commands available."
	errorExecutingCommand = "Error executing %s: %s"
)

var (
	helps    = make(map[string]Manual)
	commands = make(map[string]CommandFunc)
	irccon   IRCConnection
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
func HandleCmd(cmd *Command, irc IRCConnection) error {
	cmdFunction := commands[cmd.Command]

	if cmdFunction == nil {
		log.Println("Command not found")
		return errors.New(fmt.Sprintf(commandNotAvailable, cmd.Command))
	}

	log.Printf("HandleCmd %v args %v", cmd.Command, cmd.FullArg)
	resultStr, err := cmdFunction(cmd)
	if err != nil {
		irc.Privmsg(cmd.Channel, fmt.Sprintf(errorExecutingCommand, cmd.Command, err.Error()))
		return err
	}

	if resultStr != "" {
		irc.Privmsg(cmd.Channel, resultStr)
	}
	return nil
}
