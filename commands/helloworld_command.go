package commands

func Helloworld(cmd *Command) (msg string, err error) {
	msg = "Hello world!"
	return
}

func init() {
	RegisterCommand("helloworld", Helloworld)

	man := Manual{
		helpDescripton: "Just send a 'Hello World' message on the channel.",
	}
	RegisterHelp("helloworld", man)
}
