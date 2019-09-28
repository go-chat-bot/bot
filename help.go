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
	msg := &Message{
		Text: CmdPrefix + c.RawArgs,
	}
	cmd, _ := parse(msg, c.ChannelData, c.User)
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
		b.SendMessage(OutgoingMessage{
			Target:  c.Channel,
			Message: fmt.Sprintf(helpDescripton, help.Description),
			Sender:  c.User,
		})
	}
	b.SendMessage(OutgoingMessage{
		Target:  c.Channel,
		Message: fmt.Sprintf(helpUsage, CmdPrefix, c.Command, help.ExampleArgs),
		Sender:  c.User,
	})
}

func (b *Bot) showAvailabeCommands(channel string, sender *User) {
	var cmds []string
	for k := range commands {
		cmds = append(cmds, k)
	}
	b.SendMessage(OutgoingMessage{
		Target:  channel,
		Message: fmt.Sprintf(helpAboutCommand, CmdPrefix),
		Sender:  sender,
	})
	b.SendMessage(OutgoingMessage{
		Target:  channel,
		Message: fmt.Sprintf(availableCommands, strings.Join(cmds, ", ")),
		Sender:  sender,
	})
}
