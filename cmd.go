package bot

import (
	"fmt"
	"log"
	"sync"
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
	IsPrivate bool   // Whether the channel is a group or private chat
}

// URI gives back an URI-fied string containing protocol, server and channel.
func (c *ChannelData) URI() string {
	return fmt.Sprintf("%s://%s/%s", c.Protocol, c.Server, c.Channel)
}

// Message holds the message info - for IRC and Slack networks, this can include whether the message was an action.
type Message struct {
	Text     string // The actual content of this Message
	IsAction bool   // True if this was a '/me does something' message
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
	CronSpec string                               // CronSpec that schedules some function
	Channels []string                             // A list of channels to notify
	CmdFunc  func(channel string) (string, error) // func to be executed at the period specified on CronSpec
}

// User holds user id, nick and real name
type User struct {
	ID       string
	Nick     string
	RealName string
	IsBot    bool
}

type customCommand struct {
	Version     int
	Cmd         string
	CmdFuncV1   activeCmdFuncV1
	CmdFuncV2   activeCmdFuncV2
	CmdFuncV3   activeCmdFuncV3
	Description string
	ExampleArgs string
}

// CmdResult is the result message of V2 commands
type CmdResult struct {
	Channel string // The channel where the bot should send the message
	Message string // The message to be sent
}

// CmdResultV3 is the result message of V3 commands
type CmdResultV3 struct {
	Channel string
	Message chan string
	Done    chan bool
}

const (
	v1 = iota
	v2
	v3
)

const (
	commandNotAvailable   = "Command %v not available."
	noCommandsAvailable   = "No commands available."
	errorExecutingCommand = "Error executing %s: %s"
)

type passiveCmdFunc func(cmd *PassiveCmd) (string, error)
type activeCmdFuncV1 func(cmd *Cmd) (string, error)
type activeCmdFuncV2 func(cmd *Cmd) (CmdResult, error)
type activeCmdFuncV3 func(cmd *Cmd) (CmdResultV3, error)

var (
	commands         = make(map[string]*customCommand)
	passiveCommands  = make(map[string]passiveCmdFunc)
	periodicCommands = make(map[string]PeriodicConfig)
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

// RegisterPassiveCommand adds a new passive command to the bot.
// The command should be registered in the Init() func of your package
// Passive commands receives all the text posted to a channel without any parsing
// command: String used to identify the command, for internal use only (ex: logs)
// cmdFunc: Function which will be executed. It will received the raw message, channel and nick
func RegisterPassiveCommand(command string, cmdFunc func(cmd *PassiveCmd) (string, error)) {
	passiveCommands[command] = cmdFunc
}

// RegisterPeriodicCommand adds a command that is run periodically.
// The command should be registered in the Init() func of your package
// config: PeriodicConfig which specify CronSpec and a channel list
// cmdFunc: A no-arg function which gets triggered periodically
func RegisterPeriodicCommand(command string, config PeriodicConfig) {
	periodicCommands[command] = config
}

// Disable allows disabling commands that were registered.
// It is usefull when running multiple bot instances to disabled some plugins like url which
// is already present on some protocols.
func (b *Bot) Disable(cmds []string) {
	b.disabledCmds = append(b.disabledCmds, cmds...)
}

func (b *Bot) executePassiveCommands(cmd *PassiveCmd) {
	var wg sync.WaitGroup
	mutex := &sync.Mutex{}

	for k, v := range passiveCommands {
		if b.isDisabled(k) {
			continue
		}

		cmdFunc := v
		wg.Add(1)

		go func() {
			defer wg.Done()

			result, err := cmdFunc(cmd)
			if err != nil {
				log.Println(err)
			} else {
				mutex.Lock()
				b.handlers.Response(cmd.Channel, result, cmd.User)
				mutex.Unlock()
			}
		}()
	}
	wg.Wait()
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
		log.Printf("Command not found %v", c.Command)
		return
	}

	switch cmd.Version {
	case v1:
		message, err := cmd.CmdFuncV1(c)
		b.checkCmdError(err, c)
		if message != "" {
			b.handlers.Response(c.Channel, message, c.User)
		}
	case v2:
		result, err := cmd.CmdFuncV2(c)
		b.checkCmdError(err, c)
		if result.Channel == "" {
			result.Channel = c.Channel
		}

		if result.Message != "" {
			b.handlers.Response(result.Channel, result.Message, c.User)
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
					b.handlers.Response(result.Channel, message, c.User)
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
		log.Printf(errorMsg)
		b.handlers.Response(c.Channel, errorMsg, c.User)
	}
}
