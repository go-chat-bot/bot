package example

import (
	"fmt"
	"github.com/fabioxgn/go-bot"
)

func hello(command *bot.Cmd) (msg string, err error) {
	msg = fmt.Sprintf("Hello %s", command.Nick)
	return
}

func init() {
	bot.RegisterCommand(
		"hello",
		"Sends a 'Hello' message to you on the channel.",
		"",
		hello)
}
