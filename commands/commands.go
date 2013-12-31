package commands

type CommandFunc func(args []string) string

var (
	Commands = make(map[string]CommandFunc)
)

func RegisterCommand(command string, f CommandFunc) {
	Commands[command] = f
}
