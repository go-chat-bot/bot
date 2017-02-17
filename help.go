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
	cmd, _ := parse(CmdPrefix+c.RawArgs, c.ChannelData, c.User)
	if cmd == nil {
		b.showAvailabeCommands(c.Channel, c.User)
		return
	}

	command := commands[cmd.Command]
	if command == nil {
		b.showAvailabeCommands(c.Channel, c.User)
		return
	}

	b.showHelp(cmd, command)
}

func (b *Bot) showHelp(c *Cmd, help *customCommand) {
	if help.Description != "" {
		b.handlers.Response(c.Channel, fmt.Sprintf(helpDescripton, help.Description), c.User)
	}
	b.handlers.Response(c.Channel, fmt.Sprintf(helpUsage, CmdPrefix, c.Command, help.ExampleArgs), c.User)
}

func (b *Bot) showAvailabeCommands(channel string, sender *User) {
	var cmds []string
	for k := range commands {
		cmds = append(cmds, k)
	}
	b.handlers.Response(channel, fmt.Sprintf(helpAboutCommand, CmdPrefix), sender)
	b.handlers.Response(channel, fmt.Sprintf(availableCommands, strings.Join(cmds, ", ")), sender)
}
