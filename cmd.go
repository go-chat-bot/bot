package bot

import (
	"fmt"
	"log"
	"sync"
)

// Cmd holds the parsed user's input for easier handling of commands
type Cmd struct {
	Raw     string   // Raw is full string passed to the command
	Channel string   // Channel where the command was called
	Nick    string   // User who sent the message
	Message string   // Full string without the prefix
	Command string   // Command is the first argument passed to the bot
	FullArg string   // Full argument as a single string
	Args    []string // Arguments as array
}

// PassiveCmd holds the information which will be passed to passive commands when receiving a message
type PassiveCmd struct {
	Raw     string // Raw message sent to the channel
	Channel string // Channel which the message was sent to
	Nick    string // Nick of the user which sent the message
}

type customCommand struct {
	Version     int
	Cmd         string
	CmdFuncV1   activeCmdFuncV1
	CmdFuncV2   activeCmdFuncV2
	Description string
	ExampleArgs string
}

type incomingMessage struct {
	Channel        string
	Text           string
	SenderNick     string
	BotCurrentNick string
}

// CmdResult is the result message of V2 commands
type CmdResult struct {
	Channel string // The channel where the bot should send the message
	Message string // The message to be sent
}

const (
	v1 = iota
	v2
)

const (
	commandNotAvailable   = "Command %v not available."
	noCommandsAvailable   = "No commands available."
	errorExecutingCommand = "Error executing %s: %s"
)

type passiveCmdFunc func(cmd *PassiveCmd) (string, error)
type activeCmdFuncV1 func(cmd *Cmd) (string, error)
type activeCmdFuncV2 func(cmd *Cmd) (CmdResult, error)

var (
	commands        = make(map[string]*customCommand)
	passiveCommands = make(map[string]passiveCmdFunc)
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

// RegisterPassiveCommand adds a new passive command to the bot.
// The command should be registered in the Init() func of your package
// Passive commands receives all the text posted to a channel without any parsing
// command: String used to identify the command, for internal use only (ex: logs)
// cmdFunc: Function which will be executed. It will received the raw message, channel and nick
func RegisterPassiveCommand(command string, cmdFunc func(cmd *PassiveCmd) (string, error)) {
	passiveCommands[command] = cmdFunc
}

func isPrivateMsg(channel, currentNick string) bool {
	return channel == currentNick
}

func (b *Bot) executePassiveCommands(cmd *PassiveCmd) {
	var wg sync.WaitGroup

	for k, v := range passiveCommands {
		cmdName := k
		cmdFunc := v

		wg.Add(1)

		go func() {
			defer wg.Done()

			log.Println("Executing passive command: ", cmdName)
			result, err := cmdFunc(cmd)
			if err != nil {
				log.Println(err)
			} else {
				b.messageHandler(cmd.Channel, result, cmd.Nick)
			}
		}()
	}

	wg.Wait()
}

func (b *Bot) handleCmd(c *Cmd) {
	cmd := commands[c.Command]

	if cmd == nil {
		log.Printf("Command not found %v", c.Command)
		return
	}

	log.Printf("HandleCmd %v %v", c.Command, c.FullArg)

	switch cmd.Version {
	case v1:
		message, err := cmd.CmdFuncV1(c)
		b.checkCmdError(err, c)
		if message != "" {
			b.messageHandler(c.Channel, message, c.Nick)
		}
	case v2:
		result, err := cmd.CmdFuncV2(c)
		b.checkCmdError(err, c)
		if result.Channel == "" {
			result.Channel = c.Channel
		}

		if result.Message != "" {
			b.messageHandler(result.Channel, result.Message, c.Nick)
		}
	}

}

func (b *Bot) checkCmdError(err error, c *Cmd) {
	if err != nil {
		errorMsg := fmt.Sprintf(errorExecutingCommand, c.Command, err.Error())
		log.Printf(errorMsg)
		b.messageHandler(c.Channel, errorMsg, c.Nick)
	}
}
