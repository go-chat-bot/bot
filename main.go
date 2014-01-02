package main

import (
	"fmt"
	"github.com/fabioxgn/go-bot/commands"
	"github.com/thoj/go-ircevent"
	"log"
	"os"
	"strings"
)

const (
	CONFIG_FILE = "config.json"
)

var (
	irccon *irc.Connection
	config = &Config{}
)

func printAvailableCommands(channel string) {
	irccon.Privmsg(channel, "Available Commands:")
	cmds := ""
	for k, _ := range commands.Commands {
		cmds += k + ", "
	}
	irccon.Privmsg(channel, cmds[:len(cmds)-2])
}

func onPRIVMSG(e *irc.Event) {
	log.Println(e.Message)
	if !strings.Contains(e.Message, config.Cmd) {
		return
	}

	channel := e.Arguments[0]
	cmd, err := Parse(StrAfter(e.Message, config.Cmd))
	if err != nil {
		irccon.Privmsg(channel, err.Error())
		return
	}

	log.Printf("cmd: %v", cmd)

	irc_cmd := commands.Commands[cmd.Command]
	if irc_cmd == nil {
		irccon.Privmsg(channel, fmt.Sprintf("Command %v not found.", cmd.Command))
		printAvailableCommands(channel)
	} else {
		log.Printf("cmd %v args %v", cmd.Command, cmd.Args)
		irccon.Privmsg(channel, irc_cmd(cmd.Args))
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
	configFile, err := os.Open(CONFIG_FILE)
	if err != nil {
		panic(err)
	}
	config.Read(configFile)
	fmt.Printf("%v", config)
}

func main() {
	readConfig()
	connect()
	configureEvents()
	irccon.Loop()
}
