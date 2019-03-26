package catfacts

import (
	"fmt"
	"github.com/go-chat-bot/bot"
	"github.com/go-chat-bot/plugins/web"
	"regexp"
)

const (
	pattern   = "(?i)\\b(cat|gato|miau|meow|garfield|lolcat)[s|z]{0,1}\\b"
	msgPrefix = "I love cats! Here's a fact: %s"
)

type catFact struct {
	Fact   string   `json:"fact"`
	Length int      `json:"length"`
}

var (
	re          = regexp.MustCompile(pattern)
	catFactsURL = "http://catfact.ninja/fact"
)

func catFacts(command *bot.PassiveCmd) (string, error) {
	if !re.MatchString(command.Raw) {
		return "", nil
	}
	data := &catFact{}
	err := web.GetJSON(catFactsURL, data)
	if err != nil {
		return "", err
	}

	if len(data.Fact) == 0 {
		return "", nil
	}

	return fmt.Sprintf(msgPrefix, data.Fact), nil}

func init() {
	bot.RegisterPassiveCommand(
		"catfacts",
		catFacts)
}
