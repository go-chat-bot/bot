package bot

func Join(command Cmd, conn Connection) {
	for _, channel := range command.Args {
		conn.Join(channel)
	}
}
