package catfacts

import (
	"fmt"
	"github.com/fabioxgn/go-bot"
	"github.com/fabioxgn/go-bot/web"
	"regexp"
)

const (
	pattern     = "(?i)(\\s*?)(cat|gato|miau|meow|garfield|lolcat)(s|z)?(?![^ ?.!])"
	catFactsURL = "http://catfacts-api.appspot.com/api/facts?number=1"
	msgPrefix   = "I love cats! Here's a fact: %s"
)

type facts struct {
	Facts   []string `json:"facts"`
	Success string   `json:"success"`
}

var (
	re = regexp.MustCompile(pattern)
)

func getFacts(text string, get web.GetJSONFunc) (string, error) {
	if re.MatchString(text) {
		return getFact(get)
	}

	return "", nil
}

func getFact(get web.GetJSONFunc) (string, error) {
	data := &facts{}
	err := get(catFactsURL, data)
	if err != nil {
		return "", err
	}

	if len(data.Facts) == 0 {
		return "", nil
	}

	return fmt.Sprintf(msgPrefix, data.Facts[0]), nil
}

func catfacts(command *bot.PassiveCmd) (string, error) {
	return getFacts(command.Raw, web.GetJSON)
}

func init() {
	bot.RegisterPassiveCommand(
		"catfacts",
		catfacts)
}
