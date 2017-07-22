// Package bot provides a simple to use IRC, Slack and Telegram bot
package bot

import (
	"log"
	"math/rand"
	"time"

	"github.com/robfig/cron"
)

const (
	// CmdPrefix is the prefix used to identify a command.
	// !hello would be identified as a command
	CmdPrefix = "!"
)

// Bot handles the bot instance
type Bot struct {
	handlers     *Handlers
	cron         *cron.Cron
	disabledCmds []string
}

// ResponseHandler must be implemented by the protocol to handle the bot responses
type ResponseHandler func(target, message string, sender *User)

// Handlers that must be registered to receive callbacks from the bot
type Handlers struct {
	Response ResponseHandler
}

// New configures a new bot instance
func New(h *Handlers) *Bot {
	b := &Bot{
		handlers: h,
		cron:     cron.New(),
	}
	b.startPeriodicCommands()
	return b
}

func (b *Bot) startPeriodicCommands() {
	for _, config := range periodicCommands {
		func(b *Bot, config PeriodicConfig) {
			b.cron.AddFunc(config.CronSpec, func() {
				for _, channel := range config.Channels {
					message, err := config.CmdFunc(channel)
					if err != nil {
						log.Print("Periodic command failed ", err)
					} else if message != "" {
						b.handlers.Response(channel, message, nil)
					}
				}
			})
		}(b, config)
	}
	if len(b.cron.Entries()) > 0 {
		b.cron.Start()
	}
}

// MessageReceived must be called by the protocol upon receiving a message
func (b *Bot) MessageReceived(channel *ChannelData, message *Message, sender *User) {
	command, err := parse(message.Text, channel, sender)
	if err != nil {
		b.handlers.Response(channel.Channel, err.Error(), sender)
		return
	}

	if command == nil {
		b.executePassiveCommands(&PassiveCmd{
			Raw:         message.Text,
			MessageData: message,
			Channel:     channel.Channel,
			ChannelData: channel,
			User:        sender,
		})
		return
	}

	if b.isDisabled(command.Command) {
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
