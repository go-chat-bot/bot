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

func part(c *Cmd, channel, senderNick string, conn connection) {
	for _, configChannel := range config.Channels {
		channelName := strings.Split(configChannel, " ")
		if strings.EqualFold(channelName[0], channel) {
			conn.Privmsg(channel, partNotAllowed)
			return
		}
	}

	conn.Privmsg(channel, partMessage)
	conn.Part(channel)
}
