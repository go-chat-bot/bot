package main

import (
	"fmt"
	"github.com/thoj/go-ircevent"
	"log"
)

const (
	CHANNEL = "#lightirc"
)

func main() {
	irccon := irc.IRC("go-bot", "go-bot")
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

	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		fmt.Printf("%v", e)
		if e.Message == "die" {
			irccon.Quit()
		}
		log.Println(e.Message)
	})

	irccon.Loop()
}
