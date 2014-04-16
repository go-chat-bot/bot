package irc

type ConnectionMock struct {
	PrivMsgFunc func(target, message string)
	NoticeFunc  func(target, message string)
	JoinFunc    func(channel string)
}

func (m ConnectionMock) Privmsg(target, message string) {
	m.PrivMsgFunc(target, message)
}
func (m ConnectionMock) Notice(target, message string) {
	m.NoticeFunc(target, message)
}
func (m ConnectionMock) Join(channel string) {
	m.JoinFunc(channel)
}
func (m ConnectionMock) Part(channel string) {}
func (m ConnectionMock) Quit()               {}
