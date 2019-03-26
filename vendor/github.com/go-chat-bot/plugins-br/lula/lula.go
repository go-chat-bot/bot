package lula

import (
	"regexp"

	"github.com/go-chat-bot/bot"
)

const (
	pattern = "(?i)\\b(lula)\\b"
	resp    = "O Lula tá preso, babaca!"
)

var (
	re = regexp.MustCompile(pattern)
)

func lula(command *bot.PassiveCmd) (string, error) {
	if re.MatchString(command.Raw) {
		return resp, nil
	}
	return "", nil
}

func init() {
	bot.RegisterPassiveCommand("lula", lula)
}
