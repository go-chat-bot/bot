package commands

func Part(cmd *Command) (msg string, err error) {
	if len(cmd.Args) > 0 {
		for _, channel := range cmd.Args {
			irccon.Part(channel)
		}
	} else {
		irccon.Part(cmd.Channel)
	}
	return
}

func init() {
	RegisterCommand("part", Part)

	man := Manual{
		helpDescripton: "Leave from the specified channels",
		helpUse: "#channel1 [#channel2 ... ]",
	}
	RegisterHelp("part", man)
}
