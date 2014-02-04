package commands

func Helloworld(args []string) string {
	return "Hello world!"
}

func init() {
	RegisterCommand("helloworld", Helloworld)
}
