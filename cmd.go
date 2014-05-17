package bot

import (
	"fmt"
	"log"
	"strings"
)

// Cmd holds the parsed user's input for easier handling of commands
type Cmd struct {
	Raw     string   // Raw is full string passed to the command
	Channel string   // Channel where the command was called
	Nick    string   // User who sent the message
	Message string   // Full string without the prefix
	Command string   // Command is the first argument passed to the bot
	FullArg string   // Full argument as a single string
	Args    []string // Arguments as array
}

type customCommand struct {
	Cmd         string
	CmdFunc     func(cmd *Cmd) (string, error)
	Description string
	ExampleArgs string
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
	helpDescripton        = "Description: %s"
	helpUsage             = "Usage: %s%s %s"
	availableCommands     = "Available commands: %v"
	helpAboutCommand      = "Type: '%shelp <command>' to see details about a specific command."
	helpCommand           = "help"
)

var (
	commands = make(map[string]*customCommand)
)

// RegisterCommand adds a new command to the bot
// The command(s) should be registered in the Ini() func of the package
// command: String which the user will use to execute the command, example: reverse
// decription: Description of the command to use in !help, example: Reverses a string
// exampleArgs: Example args to be displayed in !help <command>, example: string to be reversed
// cmdFunc: Function which will be executed. It will received a parsed command as a Cmd value
func RegisterCommand(command, description, exampleArgs string, cmdFunc func(cmd *Cmd) (string, error)) {
	commands[command] = &customCommand{
		Cmd:         command,
		CmdFunc:     cmdFunc,
		Description: description,
		ExampleArgs: exampleArgs,
	}
}

func isPrivateMsg(channel, currentNick string) bool {
	return channel == currentNick
}

func messageReceived(channel, text, senderNick string, conn ircConnection) {
	if isPrivateMsg(channel, conn.GetNick()) {
		channel = senderNick // should reply in private
	}

	command := parse(text, channel, senderNick)
	if command == nil {
		handleMessage(text, channel)
		return
	}

	if command.Command == helpCommand {
		help(command, channel, senderNick, conn)
	} else {
		handleCmd(command, conn)
	}

}

func handleMessage(text, channel string) {
	//TODO handle passive commands, lik
}

func handleCmd(c *Cmd, conn ircConnection) {
	cmd := commands[c.Command]

	if cmd == nil {
		log.Printf("Command not found %v", c.Command)
		return
	}

	log.Printf("HandleCmd %v %v", c.Command, c.FullArg)

	result, err := cmd.CmdFunc(c)

	if err != nil {
		errorMsg := fmt.Sprintf(errorExecutingCommand, c.Command, err.Error())
		log.Printf(errorMsg)
		conn.Privmsg(c.Channel, errorMsg)
	}

	if result != "" {
		conn.Privmsg(c.Channel, result)
	}
}

func help(c *Cmd, channel, senderNick string, conn ircConnection) {
	cmd := parse(CmdPrefix+c.FullArg, channel, senderNick)
	if cmd == nil {
		showAvailabeCommands(channel, conn)
		return
	}

	command := commands[cmd.Command]
	if command == nil {
		showAvailabeCommands(c.Channel, conn)
		return
	}

	showHelp(cmd, command, conn)
}

func showHelp(c *Cmd, help *customCommand, conn ircConnection) {
	if help.Description != "" {
		conn.Privmsg(c.Channel, fmt.Sprintf(helpDescripton, help.Description))
	}
	conn.Privmsg(c.Channel, fmt.Sprintf(helpUsage, CmdPrefix, c.Command, help.ExampleArgs))
}

func showAvailabeCommands(channel string, conn ircConnection) {
	cmds := make([]string, 0)
	for k := range commands {
		cmds = append(cmds, k)
	}
	conn.Privmsg(channel, fmt.Sprintf(helpAboutCommand, CmdPrefix))
	conn.Privmsg(channel, fmt.Sprintf(availableCommands, strings.Join(cmds, ", ")))
}
