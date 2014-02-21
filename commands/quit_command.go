package commands

// var (
// 	helpDescripton = ""
// 	helpUse = ""
// )

func Quit(cmd *Command) (msg string, err error) {
	irccon.Quit()
	return
}

func init() {
	RegisterCommand("quit", Quit)

	man := Manual{
		helpDescripton: "Disconnect from the current network",
	}
	RegisterHelp("quit", man)
}
