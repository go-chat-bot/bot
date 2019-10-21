// Package bot provides a simple to use IRC, Slack and Telegram bot
package bot

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/robfig/cron/v3"
)

const (
	// CmdPrefix is the prefix used to identify a command.
	// !hello would be identified as a command
	CmdPrefix = "!"

	// MsgBuffer is the max number of messages which can be buffered
	// while waiting to flush them to the chat service.
	MsgBuffer = 100
)

// Bot handles the bot instance
type Bot struct {
	handlers     *Handlers
	cron         *cron.Cron
	disabledCmds []string
	msgsToSend   chan responseMessage
	done         chan struct{}

	// Protocol and Server are used by MesssageStreams to
	// determine if this is the correct bot to send a message on
	// see:
	// https://github.com/go-chat-bot/bot/issues/37#issuecomment-277661159
	// https://github.com/go-chat-bot/bot/issues/97#issuecomment-442827599
	Protocol string
	// Server and Protocol are used by MesssageStreams to
	// determine if this is the correct bot to send a message on
	// see:
	// https://github.com/go-chat-bot/bot/issues/37#issuecomment-277661159
	// https://github.com/go-chat-bot/bot/issues/97#issuecomment-442827599
	Server string
}

// Config configuration for this Bot instance
type Config struct {
	// Protocol and Server are used by MesssageStreams to
	/// determine if this is the correct bot to send a message on
	Protocol string
	// Server and Protocol are used by MesssageStreams to
	// determine if this is the correct bot to send a message on
	Server string
}

type responseMessage struct {
	target, message string
	sender          *User
	protoParams     interface{}
}

// OutgoingMessage collects all the information for a message to go out.
type OutgoingMessage struct {
	Target      string
	Message     string
	Sender      *User
	ProtoParams interface{}
}

// ResponseHandler must be implemented by the protocol to handle the bot responses
type ResponseHandler func(target, message string, sender *User)

// ResponseHandlerV2 may be implemented by the protocol to handle the bot responses
type ResponseHandlerV2 func(OutgoingMessage)

// ErrorHandler will be called when an error happens
type ErrorHandler func(msg string, err error)

// Handlers that must be registered to receive callbacks from the bot
type Handlers struct {
	Response   ResponseHandler
	ResponseV2 ResponseHandlerV2
	Errored    ErrorHandler
}

func logErrorHandler(msg string, err error) {
	log.Printf("%s: %s", msg, err.Error())
}

// New configures a new bot instance
func New(h *Handlers, bc *Config) *Bot {
	if h.Errored == nil {
		h.Errored = logErrorHandler
	}

	b := &Bot{
		handlers:   h,
		cron:       cron.New(),
		msgsToSend: make(chan responseMessage, MsgBuffer),
		done:       make(chan struct{}),
		Protocol:   bc.Protocol,
		Server:     bc.Server,
	}

	// Launch the background goroutine that isolates the possibly non-threadsafe
	// message sending logic of the underlying transport layer.
	go b.processMessages()

	b.startMessageStreams()

	b.startPeriodicCommands()
	return b
}

func (b *Bot) startMessageStreams() {
	for _, v := range messageStreamConfigs {

		go func(b *Bot, config *messageStreamConfig) {
			msMap.Lock()
			ms := &MessageStream{
				Data: make(chan MessageStreamMessage),
				Done: make(chan bool),
			}
			var err = config.msgFunc(ms)
			if err != nil {
				b.errored("MessageStream "+config.streamName+" failed ", err)
			}
			msKey := messageStreamKey{
				Protocol:   b.Protocol,
				Server:     b.Server,
				StreamName: config.streamName,
			}
			// thread safe write
			msMap.messageStreams[msKey] = ms
			msMap.Unlock()
			b.handleMessageStream(config.streamName, ms)
		}(b, v)
	}
}

func (b *Bot) startPeriodicCommands() {
	for _, config := range periodicCommands {
		func(b *Bot, config PeriodicConfig) {
			b.cron.AddFunc(config.CronSpec, func() {
				switch config.Version {
				case v1:
					for _, channel := range config.Channels {
						message, err := config.CmdFunc(channel)
						if err != nil {
							b.errored("Periodic command failed ", err)
						} else if message != "" {
							b.SendMessage(OutgoingMessage{
								Target:  channel,
								Message: message,
							})
						}
					}
				case v2:
					results, err := config.CmdFuncV2()
					if err != nil {
						b.errored("Periodic command failed ", err)
						return
					}
					for _, result := range results {
						b.SendMessage(OutgoingMessage{
							Target:      result.Channel,
							Message:     result.Message,
							ProtoParams: result.ProtoParams,
						})
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
	command, err := parse(message, channel, sender)
	if err != nil {
		b.SendMessage(OutgoingMessage{
			Target:  channel.Channel,
			Message: err.Error(),
			Sender:  sender,
		})
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

// SendMessage queues a message.
func (b *Bot) SendMessage(om OutgoingMessage) {
	message := b.executeFilterCommands(&FilterCmd{
		Target:  om.Target,
		Message: om.Message,
		User:    om.Sender,
	})
	if message == "" {
		return
	}

	select {
	case b.msgsToSend <- responseMessage{
		target:      om.Target,
		message:     message,
		sender:      om.Sender,
		protoParams: om.ProtoParams,
	}:
	default:
		b.errored("Failed to queue message to send.", errors.New("Too busy"))
	}
}

func (b *Bot) sendResponse(resp responseMessage) {
	if b.handlers.ResponseV2 != nil {
		b.handlers.ResponseV2(OutgoingMessage{
			Message:     resp.message,
			ProtoParams: resp.protoParams,
			Sender:      resp.sender,
			Target:      resp.target,
		})
		return
	}
	b.handlers.Response(resp.target, resp.message, resp.sender)
}

func (b *Bot) errored(msg string, err error) {
	if b.handlers.Errored != nil {
		b.handlers.Errored(msg, err)
	}
}

func (b *Bot) processMessages() {
	for {
		select {
		case msg := <-b.msgsToSend:
			b.sendResponse(msg)
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
