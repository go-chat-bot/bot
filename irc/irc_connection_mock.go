package irc

type ConnectionMock struct {
	PrivMsgFunc func(target, message string)
	NoticeFunc  func(target, message string)
}

func (m ConnectionMock) Privmsg(target, message string) {
	m.PrivMsgFunc(target, message)
}
func (m ConnectionMock) Notice(target, message string) {
	m.NoticeFunc(target, message)
}
func (m ConnectionMock) Join(channel string) {}
func (m ConnectionMock) Part(channel string) {}
func (m ConnectionMock) Quit()               {}
