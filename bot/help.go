package bot

import (
	"fmt"
	"strings"
)

const (
	helpDescripton    = "Description: %s"
	helpUsage         = "Usage: %s%s %s"
	availableCommands = "Available commands: %v"
	helpAboutCommand  = "Type: '%shelp <command>' to see details about a specific command."
)

func Help(c *Cmd, conn Connection) {
	command := Commands[c.Command]
	if command == nil {
		showAvailabeCommands(c.Channel, conn)
	} else {
		showHelp(c, command, conn)
	}
}

func showHelp(c *Cmd, help *CustomCommand, conn Connection) {
	if help.Description != "" {
		conn.Privmsg(c.Channel, fmt.Sprintf(helpDescripton, help.Description))
	}
	conn.Privmsg(c.Channel, fmt.Sprintf(helpUsage, CmdPrefix, c.Command, help.Usage))
}

func showAvailabeCommands(channel string, conn Connection) {
	cmds := make([]string, 0)
	for k := range Commands {
		cmds = append(cmds, k)
	}
	conn.Privmsg(channel, fmt.Sprintf(helpAboutCommand, CmdPrefix))
	conn.Privmsg(channel, fmt.Sprintf(availableCommands, strings.Join(cmds, ", ")))
}
