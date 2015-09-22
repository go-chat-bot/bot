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
	helpCommand       = "help"
)

func help(c *Cmd, channel, senderNick string, conn connection) {
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

func showHelp(c *Cmd, help *customCommand, conn connection) {
	if help.Description != "" {
		conn.Privmsg(c.Channel, fmt.Sprintf(helpDescripton, help.Description))
	}
	conn.Privmsg(c.Channel, fmt.Sprintf(helpUsage, CmdPrefix, c.Command, help.ExampleArgs))
}

func showAvailabeCommands(channel string, conn connection) {
	var cmds []string
	for k := range commands {
		cmds = append(cmds, k)
	}
	conn.Privmsg(channel, fmt.Sprintf(helpAboutCommand, CmdPrefix))
	conn.Privmsg(channel, fmt.Sprintf(availableCommands, strings.Join(cmds, ", ")))
}
