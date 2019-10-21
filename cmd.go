package bot

import (
	"errors"
	"fmt"
	"sync"
)

var (
	// ErrProtocolServerMismatch server && proto must match
	ErrProtocolServerMismatch = errors.New("the specified protocol and server do not correspond to this bot instance")
	errNoChannelSpecified     = errors.New("no channel was specified for this message")
)

// Cmd holds the parsed user's input for easier handling of commands
type Cmd struct {
	Raw         string       // Raw is full string passed to the command
	Channel     string       // Channel where the command was called
	ChannelData *ChannelData // More info about the channel, including network
	User        *User        // User who sent the message
	Message     string       // Full string without the prefix
	MessageData *Message     // Message with extra flags
	Command     string       // Command is the first argument passed to the bot
	RawArgs     string       // Raw arguments after the command
	Args        []string     // Arguments as array
}

// ChannelData holds the improved channel info, which includes protocol and server
type ChannelData struct {
	Protocol  string // What protocol the message was sent on (irc, slack, telegram)
	Server    string // The server hostname the message was sent on
	Channel   string // The channel name the message appeared in
	HumanName string // The human readable name of the channel.
	IsPrivate bool   // Whether the channel is a group or private chat
}

// URI gives back an URI-fied string containing protocol, server and channel.
func (c *ChannelData) URI() string {
	return fmt.Sprintf("%s://%s/%s", c.Protocol, c.Server, c.Channel)
}

// Message holds the message info - for IRC and Slack networks, this can include whether the message was an action.
type Message struct {
	Text     string      // The actual content of this Message
	IsAction bool        // True if this was a '/me does something' message
	ProtoMsg interface{} // The underlying object that we got from the protocol pkg
}

// FilterCmd holds information about what is output being filtered - message and
// channel where it is being sent
type FilterCmd struct {
	Target  string // Channel or user the message is being sent to
	Message string // Message text being sent
	User    *User  // User who triggered original message
}

// PassiveCmd holds the information which will be passed to passive commands when receiving a message
type PassiveCmd struct {
	Raw         string       // Raw message sent to the channel
	MessageData *Message     // Message with extra
	Channel     string       // Channel which the message was sent to
	ChannelData *ChannelData // Channel and network info
	User        *User        // User who sent this message
}

// PeriodicConfig holds a cron specification for periodically notifying the configured channels
type PeriodicConfig struct {
	Version   int
	CronSpec  string                               // CronSpec that schedules some function
	Channels  []string                             // A list of channels to notify, ignored for V2
	CmdFunc   func(channel string) (string, error) // func to be executed at the period specified on CronSpec
	CmdFuncV2 func() ([]CmdResult, error)          // func v2 to be executed at the period specified on CronSpec
}

// User holds user id, nick and real name
type User struct {
	ID       string
	Nick     string
	RealName string
	IsBot    bool
}

// MessageStream allows event information to be transmitted to an arbitrary channel
// https://github.com/go-chat-bot/bot/issues/97
type MessageStream struct {
	Data chan MessageStreamMessage
	// Done is almost never called, usually the bot should just leave the chan open
	Done chan bool
}

// MessageStreamMessage the actual Message passed back to MessageStream in a chan
type MessageStreamMessage struct {
	Message     string
	ChannelData *ChannelData
}

type customCommand struct {
	Version       int
	Cmd           string
	CmdFuncV1     activeCmdFuncV1
	CmdFuncV2     activeCmdFuncV2
	CmdFuncV3     activeCmdFuncV3
	PassiveFuncV1 passiveCmdFuncV1
	PassiveFuncV2 passiveCmdFuncV2
	FilterFuncV1  filterCmdFuncV1
	Description   string
	ExampleArgs   string
}

// CmdResult is the result message of V2 commands
type CmdResult struct {
	Channel     string // The channel where the bot should send the message
	Message     string // The message to be sent
	ProtoParams interface{}
}

// CmdResultV3 is the result message of V3 commands
type CmdResultV3 struct {
	Channel     string
	Message     chan string
	Done        chan bool
	ProtoParams interface{}
}

const (
	v1 = iota
	v2
	v3
	pv1
	pv2
	fv1
)

const (
	commandNotAvailable   = "Command %v not available."
	noCommandsAvailable   = "No commands available."
	errorExecutingCommand = "Error executing %s: %s"
)

type passiveCmdFuncV1 func(cmd *PassiveCmd) (string, error)
type passiveCmdFuncV2 func(cmd *PassiveCmd) (CmdResultV3, error)

type activeCmdFuncV1 func(cmd *Cmd) (string, error)
type activeCmdFuncV2 func(cmd *Cmd) (CmdResult, error)
type activeCmdFuncV3 func(cmd *Cmd) (CmdResultV3, error)

type filterCmdFuncV1 func(cmd *FilterCmd) (string, error)

type messageStreamFunc func(ms *MessageStream) error

type messageStreamSyncMap struct {
	sync.RWMutex
	messageStreams map[messageStreamKey]*MessageStream
}
type messageStreamKey struct {
	StreamName string
	Server     string
	Protocol   string
}

// messageStreamConfig holds the registered function for the streamname
type messageStreamConfig struct {
	version    int
	streamName string
	msgFunc    messageStreamFunc
}

var (
	commands         = make(map[string]*customCommand)
	passiveCommands  = make(map[string]*customCommand)
	filterCommands   = make(map[string]*customCommand)
	periodicCommands = make(map[string]PeriodicConfig)

	messageStreamConfigs []*messageStreamConfig

	msMap = &messageStreamSyncMap{
		messageStreams: make(map[messageStreamKey]*MessageStream),
	}
)

// RegisterCommand adds a new command to the bot.
// The command(s) should be registered in the Init() func of your package
// command: String which the user will use to execute the command, example: reverse
// decription: Description of the command to use in !help, example: Reverses a string
// exampleArgs: Example args to be displayed in !help <command>, example: string to be reversed
// cmdFunc: Function which will be executed. It will received a parsed command as a Cmd value
func RegisterCommand(command, description, exampleArgs string, cmdFunc activeCmdFuncV1) {
	commands[command] = &customCommand{
		Version:     v1,
		Cmd:         command,
		CmdFuncV1:   cmdFunc,
		Description: description,
		ExampleArgs: exampleArgs,
	}
}

// RegisterCommandV2 adds a new command to the bot.
// It is the same as RegisterCommand but the command can specify the channel to reply to
func RegisterCommandV2(command, description, exampleArgs string, cmdFunc activeCmdFuncV2) {
	commands[command] = &customCommand{
		Version:     v2,
		Cmd:         command,
		CmdFuncV2:   cmdFunc,
		Description: description,
		ExampleArgs: exampleArgs,
	}
}

// RegisterCommandV3 adds a new command to the bot.
// It is the same as RegisterCommand but the command return a chan
func RegisterCommandV3(command, description, exampleArgs string, cmdFunc activeCmdFuncV3) {
	commands[command] = &customCommand{
		Version:     v3,
		Cmd:         command,
		CmdFuncV3:   cmdFunc,
		Description: description,
		ExampleArgs: exampleArgs,
	}
}

// RegisterMessageStream adds a new message stream to the bot.
// The command should be registered in the Init() func of your package
// MessageStreams send messages to a channel
// streamName: String used to identify the command, for internal use only (ex: webhook)
// messageStreamFunc: Function which will be executed. It will received a MessageStream with a chan to push
func RegisterMessageStream(streamName string, msgFunc messageStreamFunc) {
	messageStreamConfigs = append(messageStreamConfigs, &messageStreamConfig{
		version:    v1,
		streamName: streamName,
		msgFunc:    msgFunc,
	})
}

// RegisterPassiveCommand adds a new passive command to the bot.
// The command should be registered in the Init() func of your package
// Passive commands receives all the text posted to a channel without any parsing
// command: String used to identify the command, for internal use only (ex: logs)
// cmdFunc: Function which will be executed. It will received the raw message, channel and nick
func RegisterPassiveCommand(command string, cmdFunc passiveCmdFuncV1) {
	passiveCommands[command] = &customCommand{
		Version:       pv1,
		Cmd:           command,
		PassiveFuncV1: cmdFunc,
	}
}

// RegisterPassiveCommandV2 adds a new passive command to the bot.
// The command should be registered in the Init() func of your package
// Passive commands receives all the text posted to a channel without any parsing
// command: String used to identify the command, for internal use only (ex: logs)
// cmdFunc: Function which will be executed. It will received the raw message, channel and nick
func RegisterPassiveCommandV2(command string, cmdFunc passiveCmdFuncV2) {
	passiveCommands[command] = &customCommand{
		Version:       pv2,
		Cmd:           command,
		PassiveFuncV2: cmdFunc,
	}
}

// RegisterFilterCommand adds a command that is run every time bot is about to
// send a message. The comand should be registered in the Init() func of your
// package.
// Filter commands receive message and its destination and should return
// modified version. Returning empty string prevents message being sent
// completely
// command: String used to identify the command, for internal use only (ex: silence)
// cmdFunc: Function which will be executed. It will receive the message, target
// channel and nick who triggered original message
func RegisterFilterCommand(command string, cmdFunc filterCmdFuncV1) {
	filterCommands[command] = &customCommand{
		Version:      fv1,
		Cmd:          command,
		FilterFuncV1: cmdFunc,
	}
}

// RegisterPeriodicCommand adds a command that is run periodically.
// The command should be registered in the Init() func of your package
// config: PeriodicConfig which specify CronSpec and a channel list
// cmdFunc: A function with single string argument (channel) which gets triggered periodically
func RegisterPeriodicCommand(command string, config PeriodicConfig) {
	config.Version = v1
	periodicCommands[command] = config
}

// RegisterPeriodicCommandV2 adds a command that is run periodically.
// The command should be registered in the Init() func of your package
// config: PeriodicConfig which specifies CronSpec
// cmdFuncV2: A no-arg function which gets triggered periodically
// It should return slice of CmdResults (channel and message to send to it)
func RegisterPeriodicCommandV2(command string, config PeriodicConfig) {
	config.Version = v2
	periodicCommands[command] = config
}

// Disable allows disabling commands that were registered.
// It is useful when running multiple bot instances to disabled some plugins like url which
// is already present on some protocols.
func (b *Bot) Disable(cmds []string) {
	b.disabledCmds = append(b.disabledCmds, cmds...)
}

func (b *Bot) executePassiveCommands(cmd *PassiveCmd) {
	var wg sync.WaitGroup

	for k, v := range passiveCommands {
		if b.isDisabled(k) {
			continue
		}

		wg.Add(1)

		go func(cmdFunc *customCommand) {
			defer wg.Done()

			switch cmdFunc.Version {
			case pv1:
				result, err := cmdFunc.PassiveFuncV1(cmd)
				if err != nil {
					b.errored(fmt.Sprintf("Error executing %s", cmdFunc.Cmd), err)
				} else {
					b.SendMessage(OutgoingMessage{
						Target:  cmd.Channel,
						Message: result,
						Sender:  cmd.User,
					})
				}
			case pv2:
				result, err := cmdFunc.PassiveFuncV2(cmd)
				if err != nil {
					b.errored(fmt.Sprintf("Error executing %s", cmdFunc.Cmd), err)
					return
				}
				for {
					select {
					case message := <-result.Message:
						if message != "" {
							b.SendMessage(OutgoingMessage{
								Target:      result.Channel,
								Message:     message,
								Sender:      cmd.User,
								ProtoParams: result.ProtoParams,
							})
						}
					case <-result.Done:
						return
					}
				}
			default:
			}
		}(v)
	}
	wg.Wait()
}

func (b *Bot) executeFilterCommands(cmd *FilterCmd) string {
	for k, filter := range filterCommands {
		switch filter.Version {
		case fv1:
			filtered, err := filter.FilterFuncV1(cmd)
			if err != nil {
				b.errored(fmt.Sprintf("Error executing filter %s", k), err)
				continue
			}
			cmd.Message = filtered
		}
	}
	return cmd.Message
}

func (b *Bot) isDisabled(cmd string) bool {
	for _, c := range b.disabledCmds {
		if c == cmd {
			return true
		}
	}
	return false
}

func (b *Bot) handleCmd(c *Cmd) {
	cmd := commands[c.Command]

	if cmd == nil {
		b.errored(fmt.Sprintf("Command not found %v", c.Command), errors.New("Command not found"))
		return
	}

	switch cmd.Version {
	case v1:
		message, err := cmd.CmdFuncV1(c)
		b.checkCmdError(err, c)
		if message != "" {
			b.SendMessage(OutgoingMessage{
				Target:  c.Channel,
				Message: message,
				Sender:  c.User,
			})
		}
	case v2:
		result, err := cmd.CmdFuncV2(c)
		b.checkCmdError(err, c)
		if result.Channel == "" {
			result.Channel = c.Channel
		}

		if result.Message != "" {
			b.SendMessage(OutgoingMessage{
				Target:      result.Channel,
				Message:     result.Message,
				Sender:      c.User,
				ProtoParams: result.ProtoParams,
			})
		}
	case v3:
		result, err := cmd.CmdFuncV3(c)
		b.checkCmdError(err, c)
		if result.Channel == "" {
			result.Channel = c.Channel
		}
		for {
			select {
			case message := <-result.Message:
				if message != "" {
					b.SendMessage(OutgoingMessage{
						Target:      result.Channel,
						Message:     message,
						Sender:      c.User,
						ProtoParams: result.ProtoParams,
					})
				}
			case <-result.Done:
				return
			}
		}
	}
}

func (b *Bot) checkCmdError(err error, c *Cmd) {
	if err != nil {
		errorMsg := fmt.Sprintf(errorExecutingCommand, c.Command, err.Error())
		b.errored(errorMsg, err)
		b.SendMessage(OutgoingMessage{
			Target:  c.Channel,
			Message: errorMsg,
			Sender:  c.User,
		})
	}
}

// handleMessageStream
// if there are two bots (telegram, irc) and three messsages(a, b, c) then there will be six entries in messageStreams[key]
// when a message is sent into a chan it has a good chance of arriving at the wrong bot instance
// for every message we check to see if it matched this b.Protocol and b.Server
// if it doesn't we lookup the entry in messageStreams[key] send it to *that* Data chan
func (b *Bot) handleMessageStream(streamName string, ms *MessageStream) {
	for {
		select {
		case d := <-ms.Data:

			if d.ChannelData.Protocol != b.Protocol || d.ChannelData.Server != b.Server {
				// then lookup who it *should* be sent to and send it back into *that* chan
				key := messageStreamKey{Protocol: d.ChannelData.Protocol, Server: d.ChannelData.Server, StreamName: streamName}
				msMap.RLock()
				msMap.messageStreams[key].Data <- d
				msMap.RUnlock()
				continue
			}

			// this message is meant for us!

			if d.ChannelData.Channel == "" {
				b.errored("handleMessageStream: "+d.Message, errNoChannelSpecified)
				continue
			}
			if d.Message != "" {
				b.SendMessage(OutgoingMessage{
					Target:  d.ChannelData.Channel,
					Message: d.Message,
				})
			}
		case <-ms.Done:
			return
		}
	}
}
