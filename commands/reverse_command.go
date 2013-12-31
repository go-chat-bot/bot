package commands

// From stackoverflow: http://stackoverflow.com/a/10030772
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func Reverse(args []string) string {
	s := ""
	for _, value := range args {
		s += value + " "
	}
	return ReverseString(s[:len(s)-1])
}

func init() {
	RegisterCommand("reverse", Reverse)
}
