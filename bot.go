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

	// MsgBuffer is the max number of messages which can be buffered
	// while waiting to flush them to the chat service.
	MsgBuffer = 10
)

// Bot handles the bot instance
type Bot struct {
	handlers     *Handlers
	cron         *cron.Cron
	disabledCmds []string

	synchronousMessageSending bool

	msgsToSend chan *responseMessage
	done       chan struct{}
}

type responseMessage struct {
	target, message string
	sender          *User
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
		handlers:   h,
		cron:       cron.New(),
		msgsToSend: make(chan *responseMessage, MsgBuffer),
		done:       make(chan struct{}),
	}
	// Launch the background goroutine that isolates the possibly non-threadsafe
	// message sending logic of the underlying transport layer.
	go b.messageSender()

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
						b.SendMessage(channel, message, nil)
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
		b.SendMessage(channel.Channel, err.Error(), sender)
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

// SendMessage queues a message for a target recipient, optionally from a particular sender.
func (b *Bot) SendMessage(target string, message string, sender *User) {
	if b.synchronousMessageSending {
		b.handlers.Response(target, message, sender)
		return
	}

	select {
	case b.msgsToSend <- &responseMessage{
		target, message, sender,
	}:
	default:
		log.Printf("Failed to queue message to send. Must be busy.")
	}
}

func (b *Bot) messageSender() {
	for {
		select {
		case msg := <-b.msgsToSend:
			b.handlers.Response(msg.target, msg.message, msg.sender)
		case <-b.done:
			return
		}
	}
}

// Close will shut down the message sending capabilities of this bot. Call
// this when you are done using the bot.
func (b *Bot) Close() {
	close(b.done)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
