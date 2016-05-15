package bot

import (
	"regexp"
	"strings"
)

var (
	re = regexp.MustCompile("\\s+") // Matches one or more spaces
)

func parse(s string, channel string, user *User) *Cmd {
	c := &Cmd{Raw: s}
	s = strings.TrimSpace(s)

	if !strings.HasPrefix(s, CmdPrefix) {
		return nil
	}

	c.Channel = strings.TrimSpace(channel)
	c.User = user

	// Trim the prefix and extra spaces
	c.Message = strings.TrimPrefix(s, CmdPrefix)
	c.Message = strings.TrimSpace(c.Message)

	// check if we have the command and not only the prefix
	if c.Message == "" {
		return nil
	}

	// get the command
	pieces := strings.SplitN(c.Message, " ", 2)
	c.Command = pieces[0]

	if len(pieces) > 1 {
		// get the arguments and remove extra spaces
		c.RawArgs = strings.TrimSpace(pieces[1])
		c.Args = strings.Split(removeExtraSpaces(c.RawArgs), " ")
	}

	return c
}

func removeExtraSpaces(args string) string {
	return re.ReplaceAllString(strings.TrimSpace(args), " ")
}
