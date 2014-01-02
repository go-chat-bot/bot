package main

import (
	"errors"
	"strings"
)

type Command struct {
	Raw     string
	Command string
	Args    []string
}

// Parses the arguments returning the Command to execute and the Arguments passed to it
func Parse(c string) (*Command, error) {
	if c == "" {
		return nil, errors.New("Empty params")
	}

	cmd := &Command{Raw: c}

	values := strings.SplitN(strings.Trim(c, " "), " ", 2)

	cmd.Command = values[0]
	if len(values) > 1 {
		cmd.Args = strings.Split(values[1], " ")
	}

	return cmd, nil
}
