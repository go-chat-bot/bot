package bot

import (
	"regexp"
	"strings"
)

// Parse the arguments returning the Command to execute and the arguments passed to it
func Parse(s string, prefix string, channel string, nick string) *Cmd {
	c := &Cmd{Raw: s}
	s = strings.TrimSpace(s)
	c.IsCommand = strings.HasPrefix(s, prefix)
	c.Channel = strings.TrimSpace(channel)
	c.Nick = strings.TrimSpace(nick)

	// we can stop here if no prefix is detected
	if !c.IsCommand {
		c.Message = s
		return c
	}

	// Trim the prefix and extra spaces
	c.Message = strings.TrimPrefix(s, prefix)
	c.Message = strings.TrimSpace(c.Message)

	// check if we have the command and not only the prefix
	c.IsCommand = c.Message != ""
	if !c.IsCommand {
		return c
	}
	c.Prefix = strings.TrimSpace(prefix)

	// get the command
	pieces := strings.SplitN(c.Message, " ", 2)
	c.Command = pieces[0]

	if len(pieces) > 1 {
		// get the arguments and remove extra spaces
		c.FullArg = removeExtraSpaces(pieces[1])
		c.Args = strings.Split(c.FullArg, " ")
	}

	return c
}

func removeExtraSpaces(args string) string {
	reg := regexp.MustCompile("\\s+") // Matches one or more spaces
	return reg.ReplaceAllString(strings.TrimSpace(args), " ")
}
