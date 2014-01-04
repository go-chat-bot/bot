package main

import (
	"strings"
)

// Struct which separates the user's input for easier handling of commands
type Command struct {
	Raw     string
	Command string
	Args    []string
}

// Parse the arguments returning the Command to execute and the arguments passed to it
func Parse(c string) *Command {
	cmd := &Command{Raw: c}

	values := strings.SplitN(strings.Trim(c, " "), " ", 2)

	cmd.Command = values[0]
	if len(values) > 1 {
		cmd.Args = strings.Split(values[1], " ")
	}

	return cmd
}
