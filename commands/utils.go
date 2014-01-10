package commands

// ArgsToString converts the argumets passed to a command to string, using space as separator
func ArgsToString(args []string) string {
	s := ""
	for _, value := range args {
		s += value + " "
	}
	return s[:len(s)-1]
}
