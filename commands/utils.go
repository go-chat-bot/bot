package commands

func ArgsToString(args []string) string {
	s := ""
	for _, value := range args {
		s += value + " "
	}
	return s[:len(s)-1]
}
