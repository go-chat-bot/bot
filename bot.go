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

// ResponseHandler must be implemented by the protocol to handle the bot responses
type ResponseHandler func(target, message, sender string)

// Handlers that must be registered to receive callbacks from the bot
type Handlers struct {
	Response ResponseHandler
}

var (
	handlers *Handlers
)

// New configures a new bot instance
func New(h *Handlers) {
	handlers = h
}

// MessageReceived must be called by the protocol upon receiving a message
func MessageReceived(channel, text, sender string) {
	command := parse(text, channel, sender)
	if command == nil {
		executePassiveCommands(&PassiveCmd{
			Raw:     text,
			Channel: channel,
			Nick:    sender,
		})
		return
	}

	switch command.Command {
	case helpCommand:
		help(command)
	default:
		handleCmd(command)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
