package commands

import (
	"regexp"
	"strings"
)

// Command is a struct which separates the user's input for easier handling of commands
type Command struct {
	Raw       string   // Raw is full string passed to the command
	Channel   string   // Channel where the command was called
	Nick   	 string   // User who sent the message
	IsCommand bool     // Confirmation if this is a command or just a regular message
	Message   string   // Full string without the prefix
	Command   string   // Command is the first argument passed to the bot
	Prefix    string   // Command prefix
	FullArg   string   // Full argument as a single string
	Args      []string // Arguments as array
}

// Parse the arguments returning the Command to execute and the arguments passed to it
func Parse(c string, prefix string, channel string, nick string) *Command {
	cmd := &Command{Raw: c}
	c = strings.TrimSpace(c)
	cmd.IsCommand = strings.HasPrefix(c, prefix)
	cmd.Channel = strings.TrimSpace(channel)
	cmd.Nick = strings.TrimSpace(nick)

	// we can stop here if no prefix is detected
	if !cmd.IsCommand {
		cmd.Message = c
		return cmd
	}

	// Trim the prefix and extra spaces
	cmd.Message = strings.TrimPrefix(c, prefix)
	cmd.Message = strings.TrimSpace(cmd.Message)

	// check if we have the command and not only the prefix
	cmd.IsCommand = cmd.Message != ""
	if !cmd.IsCommand {
		return cmd
	}
	cmd.Prefix = strings.TrimSpace(prefix)

	// get the command
	pieces := strings.SplitN(cmd.Message, " ", 2)
	cmd.Command = pieces[0]

	if len(pieces) > 1 {
		// get the arguments and remove extra spaces
		cmd.FullArg = removeExtraSpaces(pieces[1])
		cmd.Args = strings.Split(cmd.FullArg, " ")
	}

	return cmd
}

func removeExtraSpaces(args string) string {
	reg := regexp.MustCompile("\\s+") // Matches one or more spaces
	return reg.ReplaceAllString(strings.TrimSpace(args), " ")
}
