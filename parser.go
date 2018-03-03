package bot

import (
	"errors"
	"regexp"
	"strings"

	"github.com/mattn/go-shellwords"
)

var (
	re = regexp.MustCompile("\\s+") // Matches one or more spaces
)

func parse(s string, channel *ChannelData, user *User) (*Cmd, error) {
	c := &Cmd{Raw: s}
	s = strings.TrimSpace(s)

	if !strings.HasPrefix(s, CmdPrefix) {
		return nil, nil
	}

	c.Channel = strings.TrimSpace(channel.Channel)
	c.ChannelData = channel
	c.User = user

	// Trim the prefix and extra spaces
	c.Message = strings.TrimPrefix(s, CmdPrefix)
	c.Message = strings.TrimSpace(c.Message)

	// check if we have the command and not only the prefix
	if c.Message == "" {
		return nil, nil
	}

	// get the command
	pieces := strings.SplitN(c.Message, " ", 2)
	c.Command = pieces[0]

	if len(pieces) > 1 {
		// get the arguments and remove extra spaces
		c.RawArgs = strings.TrimSpace(pieces[1])
		parsedArgs, err := shellwords.Parse(c.RawArgs)
		if err != nil {
			return nil, errors.New("Error parsing arguments: " + err.Error())
		}
		c.Args = parsedArgs
	}

	c.MessageData = &Message{
		Text: c.Message,
	}

	return c, nil
}
