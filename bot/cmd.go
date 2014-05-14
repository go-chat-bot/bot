package bot

import (
	"errors"
	"fmt"
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
func RegisterCommand(c *CustomCommand) {
	Commands[c.Cmd] = c
}

// HandleCmd handles a command
func HandleCmd(c *Cmd, conn Connection) error {
	customCmd := Commands[c.Command]

	if customCmd == nil {
		log.Printf("Command not found %v", c.Command)
		return errors.New(fmt.Sprintf(commandNotAvailable, c.Command))
	}

	log.Printf("HandleCmd %v args %v", c.Command, c.FullArg)
	resultStr, err := customCmd.CmdFunc(c)
	if err != nil {
		conn.Privmsg(c.Channel, fmt.Sprintf(errorExecutingCommand, c.Command, err.Error()))
		return err
	}

	if resultStr != "" {
		conn.Privmsg(c.Channel, resultStr)
	}
	return nil
}
