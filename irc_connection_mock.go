package bot

type ircConnectionMock struct {
	PrivMsgFunc func(target, message string)
	NoticeFunc  func(target, message string)
	JoinFunc    func(channel string)
	PartFunc    func(channel string)
	Nick        string
}

func (m ircConnectionMock) Privmsg(target, message string) {
	m.PrivMsgFunc(target, message)
}

func (m ircConnectionMock) Notice(target, message string) {
	m.NoticeFunc(target, message)
}

func (m ircConnectionMock) Join(channel string) {
	m.JoinFunc(channel)
}

func (m ircConnectionMock) Part(channel string) {
	m.PartFunc(channel)
}

func (m ircConnectionMock) Quit() {}

func (m ircConnectionMock) GetNick() string {
	return m.Nick
}
