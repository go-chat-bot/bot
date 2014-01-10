package commands

// From stackoverflow: http://stackoverflow.com/a/10030772
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Reverse reverses a string
func Reverse(args []string) string {
	return reverseString(ArgsToString(args))
}

func init() {
	RegisterCommand("reverse", Reverse)
}
