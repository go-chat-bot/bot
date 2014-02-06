package main

import (
	"fmt"
	"github.com/fabioxgn/go-bot/commands"
	"github.com/thoj/go-ircevent"
	"log"
	"os"
)

const (
	configFile = "config.json"
)

var (
	irccon *irc.Connection
	config = &Config{}
)

func onPRIVMSG(e *irc.Event) {
	channel := e.Arguments[0]
	if channel == config.Nick {
		// channel = e.Nick
		// msg := e.Message
	} else {
		// Parse the raw message
		message := commands.Parse(e.Message, config.CmdPrefix)

		// Is it a command or just a regular message?
		if message.IsCommand {
			// It's a command
			commands.HandleCmd(message, channel, irccon.Privmsg)
		} else {
			// It's not a command
			// Test for passive commands (parse url, etc) ?
		}
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

func configureEvents() {
	irccon.AddCallback("001", func(e *irc.Event) {
		irccon.Join(config.Channels[0])
	})

	irccon.AddCallback("366", func(e *irc.Event) {
		irccon.Privmsg(config.Channels[0], "Hi there.\n")
	})

	irccon.AddCallback("PRIVMSG", onPRIVMSG)
}

func readConfig() {
	file, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	config.Read(file)
	fmt.Printf("%v\n", config)
}

func main() {
	readConfig()
	connect()
	configureEvents()
	irccon.Loop()
}
