package puppet

import (
	"errors"
	"github.com/fabioxgn/go-bot"
	"strings"
)

const (
	seeUsage = "Invalid args, see usage."
)

func sendMessage(command *bot.Cmd) (result bot.CmdResult, err error) {
	result = bot.CmdResult{}
	if len(command.Args) < 2 {
		return result, errors.New(seeUsage)
	}

	if command.Args[0] != "say" && command.Args[0] != "me" {
		return result, errors.New(seeUsage)
	}

	result.Channel = command.Args[1]
	result.Message = strings.Join(command.Args[2:], " ")
	return result, nil
}

func init() {
	bot.RegisterCommandV2(
		"puppet",
		"Allows you to control what the bot says or acts",
		"say|me your message",
		sendMessage)
}
