package example

import (
	"fmt"
	"github.com/fabioxgn/go-bot/cmd"
)

func hello(command *cmd.Cmd) (msg string, err error) {
	msg = fmt.Sprintf("Hello %s", command.Nick)
	return
}

func init() {
	cmd.RegisterCommand(&cmd.CustomCommand{
		Cmd:         "hello",
		CmdFunc:     hello,
		Description: "Sends a 'Hello' message to you on the channel.",
	})
}
