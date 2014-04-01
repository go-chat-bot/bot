package cmd

import (
	"errors"
	"fmt"
	"github.com/fabioxgn/go-bot/irc"
	"log"
)

// Cmd is a struct which separates the user's input for easier handling of commands
type Cmd struct {
	Raw       string   // Raw is full string passed to the command
	Channel   string   // Channel where the command was called
	Nick      string   // User who sent the message
	IsCommand bool     // Confirmation if this is a command or just a regular message
	Message   string   // Full string without the prefix
	Command   string   // Command is the first argument passed to the bot
	Prefix    string   // Command prefix
	FullArg   string   // Full argument as a single string
	Args      []string // Arguments as array
}

type CustomCommand struct {
	Cmd         string
	CmdFunc     func(cmd *Cmd) (string, error)
	Description string
	Usage       string
}

const (
	commandNotAvailable   = "Command %v not available."
	noCommandsAvailable   = "No commands available."
	errorExecutingCommand = "Error executing %s: %s"
)

var (
	Commands = make(map[string]*CustomCommand)
)

// RegisterCommand must be used to register a command (string) and CommandFunc.
// The CommandFunc will be executed when a users calls the bot passing the
// command string as the first argument
func RegisterCommand(cmd *CustomCommand) {
	Commands[cmd.Cmd] = cmd
}

// HandleCmd handles a command
func HandleCmd(cmd *Cmd, conn irc.Connection) error {
	customCmd := Commands[cmd.Command]

	if customCmd == nil {
		log.Println("Command not found")
		return errors.New(fmt.Sprintf(commandNotAvailable, cmd.Command))
	}

	log.Printf("HandleCmd %v args %v", cmd.Command, cmd.FullArg)
	resultStr, err := customCmd.CmdFunc(cmd)
	if err != nil {
		conn.Privmsg(cmd.Channel, fmt.Sprintf(errorExecutingCommand, cmd.Command, err.Error()))
		return err
	}

	if resultStr != "" {
		conn.Privmsg(cmd.Channel, resultStr)
	}
	return nil
}
