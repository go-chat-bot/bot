package bot

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	"github.com/mattn/go-shellwords"
	unidecode "github.com/mozillazg/go-unidecode"
)

var (
	re = regexp.MustCompile("\\s+") // Matches one or more spaces
)

func parse(m *Message, channel *ChannelData, user *User) (*Cmd, error) {
	s := strings.TrimSpace(m.Text)

	if !strings.HasPrefix(s, CmdPrefix) {
		return nil, nil
	}

	c := &Cmd{
		Channel:     strings.TrimSpace(channel.Channel),
		ChannelData: channel,
		Message:     strings.TrimSpace(strings.TrimPrefix(s, CmdPrefix)),
		Raw:         m.Text,
		User:        user,
	}

	// check if we have the command and not only the prefix
	if c.Message == "" {
		return nil, nil
	}

	firstOccurrence := true
	firstUnicodeSpace := func(c rune) bool {
		isFirstSpace := unicode.IsSpace(c) && firstOccurrence
		if isFirstSpace {
			firstOccurrence = false
		}
		return isFirstSpace
	}

	// get the command
	pieces := strings.FieldsFunc(c.Message, firstUnicodeSpace)
	c.Command = strings.ToLower(unidecode.Unidecode(pieces[0]))

	if len(pieces) > 1 {
		// get the arguments and remove extra spaces
		c.RawArgs = strings.TrimSpace(pieces[1])
		parsedArgs, err := shellwords.Parse(c.RawArgs)
		if err != nil {
			return nil, errors.New("Error parsing arguments: " + err.Error())
		}
		c.Args = parsedArgs
	}

	m.Text = c.Message
	c.MessageData = m

	return c, nil
}
