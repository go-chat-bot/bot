package commands

// From stackoverflow: http://stackoverflow.com/a/10030772
func Reverse(cmd *Command) (msg string, err error) {
	runes := []rune(cmd.FullArg)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	msg = string(runes)
	return
}

func init() {
	RegisterCommand("reverse", Reverse)

	man := Manual{
		helpDescripton: "Reverse the whole string",
		helpUse: "all your base are belong to us",
	}
	RegisterHelp("reverse", man)
}
