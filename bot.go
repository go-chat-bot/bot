// Package bot provides a simple to use IRC bot
package bot

import (
	"math/rand"
	"time"
)

const (
	// CmdPrefix is the prefix used to identify a command.
	// !hello whould be identified as a command
	CmdPrefix = "!"
)

type MessageHandler func(target, message, sender string)
type NickHandler func() string

type Bot struct {
	channels       []string
	messageHandler MessageHandler
}

func NewBot(messageHandler MessageHandler, channels []string) *Bot {
	b := &Bot{}
	b.messageHandler = messageHandler
	b.channels = channels
	return b
}

// MessageReceived must be called by the protocol handler upon receiving a message
func (b *Bot) MessageReceived(channel, text, sender string) {
	command := parse(text, channel, sender)
	if command == nil {
		b.executePassiveCommands(&PassiveCmd{
			Raw:     text,
			Channel: channel,
			Nick:    sender,
		})
		return
	}

	switch command.Command {
	case helpCommand:
		b.help(command)
	default:
		b.handleCmd(command)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
