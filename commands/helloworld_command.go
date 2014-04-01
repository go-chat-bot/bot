package commands

func Helloworld(cmd *Command) (msg string, err error) {
	msg = "Hello world!"
	return
}

func init() {
	RegisterCommand(&CustomCommand{
		Cmd:         "helloworld",
		CmdFunc:     Helloworld,
		Description: "Just send a 'Hello World' message on the channel.",
	})
}
