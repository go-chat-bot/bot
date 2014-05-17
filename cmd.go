package bot

import (
	"fmt"
	"log"
)

// Cmd holds the parsed user's input for easier handling of commands
type Cmd struct {
	Raw       string   // Raw is full string passed to the command
	Channel   string   // Channel where the command was called
	Nick      string   // User who sent the message
	IsCommand bool     // Confirmation if this is a command or just a regular message
	Message   string   // Full string without the prefix
	Command   string   // Command is the first argument passed to the bot
	FullArg   string   // Full argument as a single string
	Args      []string // Arguments as array
}

// TODO
type CustomCommand struct {
	Cmd         string
	CmdFunc     func(cmd *Cmd) (string, error)
	Description string
	Usage       string
}

type incomingMessage struct {
	Channel        string
	Text           string
	SenderNick     string
	BotCurrentNick string
}

const (
	commandNotAvailable   = "Command %v not available."
	noCommandsAvailable   = "No commands available."
	errorExecutingCommand = "Error executing %s: %s"
)

var (
	commands = make(map[string]*CustomCommand)
)

// RegisterCommand must be used to register a CustomCommand.
// The commands must be registered in the Ini() func
// The CustomCommand must have at least:
// Cmd: The string which the user will use to execute the command
// CmdFunc: The function which will be executed when the Cmd string is detected as a command
func RegisterCommand(c *CustomCommand) {
	commands[c.Cmd] = c
}

func isPrivateMsg(channel, currentNick string) bool {
	return channel == currentNick
}

func messageReceived(channel, text, senderNick string, conn ircConnection) {
	if isPrivateMsg(channel, conn.GetNick()) {
		channel = senderNick // should reply in private
	}

	command := parse(text, channel, senderNick)
	if command.Command == "help" {
		command = parse(CmdPrefix+command.FullArg, channel, senderNick)
		help(command, conn)
	} else if command.IsCommand {
		handleCmd(command, conn)
	} else {
		// It's not a command
		// TODO: Test for passive commands (parse url, etc) ?
	}
}

func handleCmd(c *Cmd, conn ircConnection) {
	customCmd := commands[c.Command]

	if customCmd == nil {
		log.Printf("Command not found %v", c.Command)
		return
	}

	log.Printf("HandleCmd %v %v", c.Command, c.FullArg)

	result, err := customCmd.CmdFunc(c)

	if err != nil {
		errorMsg := fmt.Sprintf(errorExecutingCommand, c.Command, err.Error())
		log.Printf(errorMsg)
		conn.Privmsg(c.Channel, errorMsg)
	}

	if result != "" {
		conn.Privmsg(c.Channel, result)
	}
}
