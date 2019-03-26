package guid

import (
	"strings"

	uuid "github.com/beevik/guid"
	"github.com/go-chat-bot/bot"
)

const (
	msgInvalidAmountOfParams = "Invalid amount of parameters"
	msgInvalidParam          = "Invalid parameter"
)

func guid(command *bot.Cmd) (string, error) {

	if len(command.Args) > 1 {
		return msgInvalidAmountOfParams, nil
	}

	if len(command.Args) == 1 {
		if command.Args[0] == "upper" {
			return strings.ToUpper(uuid.NewString()), nil
		}
		return msgInvalidParam, nil
	}
	return uuid.NewString(), nil
}

func init() {
	bot.RegisterCommand(
		"guid",
		"Generates GUID",
		"",
		guid)
}
