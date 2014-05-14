package main

import (
	"fmt"
	"github.com/fabioxgn/go-bot/cmd"
	"github.com/fabioxgn/go-bot/cmd/parser"
	_ "github.com/fabioxgn/go-bot/commands/cotacao"
	_ "github.com/fabioxgn/go-bot/commands/example"
	_ "github.com/fabioxgn/go-bot/commands/megasena"
	"github.com/thoj/go-ircevent"
	"log"
	"os"
	"strings"
)

const (
	channelSeparator = ","
)

type Config struct {
	Server    string
	Channels  []string
	User      string
	Nick      string
	CmdPrefix string
	UseTLS    bool
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

	command := parser.Parse(e.Message(), config.CmdPrefix, channel, e.Nick)
	if command.Command == "help" {
		command = parser.Parse(config.CmdPrefix+command.FullArg, config.CmdPrefix, channel, e.Nick)
		Help(command, irccon)
	} else if command.IsCommand {
		cmd.HandleCmd(command, irccon)
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

func readConfig() {
	//TODO: Validar config
	config = &Config{
		Server:    os.Getenv("IRC_SERVER"),
		Channels:  strings.Split(os.Getenv("IRC_CHANNELS"), channelSeparator),
		User:      os.Getenv("IRC_USER"),
		Nick:      os.Getenv("IRC_NICK"),
		CmdPrefix: "!",
		UseTLS:    true,
	}
	fmt.Printf("%v\n", config)
}

func main() {
	readConfig()
	connect()
	configureEvents()
	irccon.Loop()
}
