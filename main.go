package main

import (
	"fmt"
	"github.com/fabioxgn/go-bot/commands"
	"github.com/thoj/go-ircevent"
	"log"
	"strings"
)

const (
	CHANNEL        = "#bandodeputos"
	CMD_IDENTIFIER = "!go-bot"
)

var (
	irccon = irc.IRC("go-bot", "go-bot")
)

func printAvailableCommands(c string) {
	irccon.Privmsg(CHANNEL, fmt.Sprintf("Command %v not found.", c))
	irccon.Privmsg(CHANNEL, "Available Commands:")
	cmds := ""
	for k, _ := range commands.Commands {
		cmds += k + ", "
	}
	irccon.Privmsg(CHANNEL, cmds[:len(cmds)-2])
}

func OnPRIVMSG(e *irc.Event) {
	log.Println(e.Message)
	if !strings.Contains(e.Message, CMD_IDENTIFIER) {
		return
	}

	cmd, err := Parse(StrAfter(e.Message, CMD_IDENTIFIER))
	if err != nil {
		irccon.Privmsg(CHANNEL, err.Error())
		return
	}

	log.Printf("cmd: %v", cmd)

	irc_cmd := commands.Commands[cmd.Command]
	if irc_cmd == nil {
		printAvailableCommands(cmd.Command)
	} else {
		log.Printf("cmd %v args %v", cmd.Command, cmd.Args)
		irccon.Privmsg(CHANNEL, irc_cmd(cmd.Args))
	}
}

func main() {
	irccon.UseTLS = true
	err := irccon.Connect("irc.freenode.net:7000")
	if err != nil {
		log.Fatal(err)
	}
	irccon.AddCallback("001", func(e *irc.Event) {
		irccon.Join(CHANNEL)
	})

	irccon.AddCallback("366", func(e *irc.Event) {
		irccon.Privmsg(CHANNEL, "Hi there.\n")
	})

	irccon.AddCallback("PRIVMSG", OnPRIVMSG)

	irccon.Loop()
}
