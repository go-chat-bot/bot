package bot

const (
	partMessage = "As you wish, master!"
	partCommand = "part"
	partUsage   = "Usage: !part"
)

func part(c *Cmd, channel, senderNick string, conn ircConnection) {
	conn.Privmsg(channel, partMessage)
	conn.Part(channel)
}
