package main

import (
	"fmt"
	"github.com/fabioxgn/go-bot/cmd"
	"github.com/fabioxgn/go-bot/irc"
	"strings"
)

const (
	helpDescripton    = "Description: %s"
	helpUsage         = "Usage: %s%s %s"
	availableCommands = "Available commands: %v"
	helpAboutCommand  = "Type: '%shelp <command>' to see details about a specific command."
)

func Help(c *cmd.Cmd, conn irc.Connection) {
	command := cmd.Commands[c.Command]
	if command == nil {
		showAvailabeCommands(c.Nick, c.Prefix, conn)
	} else {
		showHelp(c, command, conn)
	}
}

func showHelp(c *cmd.Cmd, help *cmd.CustomCommand, conn irc.Connection) {
	if help.Description != "" {
		conn.Notice(c.Nick, fmt.Sprintf(helpDescripton, help.Description))
	}
	if help.Usage != "" {
		conn.Notice(c.Nick, fmt.Sprintf(helpUsage, c.Prefix, c.Command, help.Usage))
	}
}

func showAvailabeCommands(nick, cmdPrefix string, conn irc.Connection) {
	cmds := make([]string, 0)
	for k := range cmd.Commands {
		cmds = append(cmds, k)
	}
	conn.Notice(nick, fmt.Sprintf(helpAboutCommand, cmdPrefix))
	conn.Notice(nick, fmt.Sprintf(availableCommands, strings.Join(cmds, ", ")))
}
