package bot

import (
	"github.com/thoj/go-ircevent"
	"log"
)

const (
	CmdPrefix = "!"
)

type Config struct {
	Server   string
	Channels []string
	User     string
	Nick     string
	UseTLS   bool
}

var (
	irccon *irc.Connection
	config *Config
)

func isPrivateMsg(channel string) bool {
	return channel == config.Nick
}

func onPRIVMSG(e *irc.Event) {
	channel := e.Arguments[0]

	if isPrivateMsg(channel) {
		channel = e.Nick //e.Nick is who sent the pvt message
	}

	command := Parse(e.Message(), channel, e.Nick)
	if command.Command == "help" {
		command = Parse(CmdPrefix+command.FullArg, channel, e.Nick)
		Help(command, irccon)
	} else if command.IsCommand {
		HandleCmd(command, irccon)
	} else {
		// It's not a command
		// TODO: Test for passive commands (parse url, etc) ?
	}
}

func connect() {
	irccon = irc.IRC(config.User, config.Nick)
	irccon.UseTLS = config.UseTLS
	err := irccon.Connect(config.Server)
	if err != nil {
		log.Fatal(err)
	}
}

func onWelcome(e *irc.Event) {
	for _, channel := range config.Channels {
		irccon.Join(channel)
	}
}

func onEndOfNames(e *irc.Event) {
	log.Println("onEndOfNames: %v", e.Arguments)
	irccon.Privmsg(e.Arguments[1], "Hi there.\n")
}

func configureEvents() {
	irccon.AddCallback("001", onWelcome)
	irccon.AddCallback("366", onEndOfNames)
	irccon.AddCallback("PRIVMSG", onPRIVMSG)
}

func Run(c *Config) {
	config = c
	connect()
	configureEvents()
	irccon.Loop()
}
