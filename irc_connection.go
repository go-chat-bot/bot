package bot

type ircConnection interface {
	Privmsg(target, message string)
	GetNick() string
}
