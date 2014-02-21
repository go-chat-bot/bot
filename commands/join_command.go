package commands

func Join(cmd *Command) (msg string, err error) {
	for _, channel := range cmd.Args {
		irccon.Join(channel)
	}
	return
}

func init() {
	RegisterCommand("join", Join)

	man := Manual{
		helpDescripton: "Join the specified channels",
		helpUse: "#channel1 [#channel2 ... ]",
	}
	RegisterHelp("join", man)
}
