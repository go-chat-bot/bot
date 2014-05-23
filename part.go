package bot

import (
	"strings"
)

const (
	partMessage    = "As you wish, master!"
	partCommand    = "part"
	partUsage      = "Usage: !part"
	partNotAllowed = "Nope!"
)

func part(c *Cmd, channel, senderNick string, conn ircConnection) {
	for _, c := range config.Channels {
		if strings.EqualFold(c, channel) {
			conn.Privmsg(channel, partNotAllowed)
			return
		}

	}

	conn.Privmsg(channel, partMessage)
	conn.Part(channel)
}
