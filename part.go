package bot

func Part(command *Cmd, conn Connection) (msg string, err error) {
	if len(command.Args) > 0 {
		for _, channel := range command.Args {
			conn.Part(channel)
		}
	} else {
		conn.Part(command.Channel)
	}
	return
}
