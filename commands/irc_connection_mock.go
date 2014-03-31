package commands

type IRCConnectionMock struct {
	PrivMsgFunc func(target, message string)
}

func (m IRCConnectionMock) Privmsg(target, message string) {
	m.PrivMsgFunc(target, message)
}
func (m IRCConnectionMock) Notice(target, message string) {}
func (m IRCConnectionMock) Join(channel string)           {}
func (m IRCConnectionMock) Part(channel string)           {}
func (m IRCConnectionMock) Quit()                         {}
