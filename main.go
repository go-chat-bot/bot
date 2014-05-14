package main

import (
	"github.com/fabioxgn/go-bot/bot"
	_ "github.com/fabioxgn/go-bot/commands/cotacao"
	_ "github.com/fabioxgn/go-bot/commands/example"
	_ "github.com/fabioxgn/go-bot/commands/megasena"
	"log"
	"os"
	"strings"
)

func main() {
	config := &bot.Config{
		Server:   os.Getenv("IRC_SERVER"),
		Channels: strings.Split(os.Getenv("IRC_CHANNELS"), bot.ChannelSeparator),
		User:     os.Getenv("IRC_USER"),
		Nick:     os.Getenv("IRC_NICK"),
		UseTLS:   true,
	}
	log.Printf("%v\n", config)

	bot.Run(config)
}
