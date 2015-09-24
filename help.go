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

func (b *Bot) help(c *Cmd) {
	cmd := parse(CmdPrefix+c.FullArg, c.Channel, c.Nick)
	if cmd == nil {
		b.showAvailabeCommands(c.Channel, c.Nick)
		return
	}

	command := commands[cmd.Command]
	if command == nil {
		b.showAvailabeCommands(c.Channel, c.Nick)
		return
	}

	b.showHelp(cmd, command)
}

func (b *Bot) showHelp(c *Cmd, help *customCommand) {
	if help.Description != "" {
		b.messageHandler(c.Channel, fmt.Sprintf(helpDescripton, help.Description), c.Nick)
	}
	b.messageHandler(c.Channel, fmt.Sprintf(helpUsage, CmdPrefix, c.Command, help.ExampleArgs), c.Nick)
}

func (b *Bot) showAvailabeCommands(channel, sender string) {
	var cmds []string
	for k := range commands {
		cmds = append(cmds, k)
	}
	b.messageHandler(channel, fmt.Sprintf(helpAboutCommand, CmdPrefix), sender)
	b.messageHandler(channel, fmt.Sprintf(availableCommands, strings.Join(cmds, ", ")), sender)
}
