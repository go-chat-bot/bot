package bot

type ircConnectionMock struct {
	Channel  string
	Messages []string
	Nick     string
}

func (m *ircConnectionMock) Privmsg(target, message string) {
	m.Channel = target
	m.Messages = append(m.Messages, message)
}

func (m ircConnectionMock) GetNick() string {
	return m.Nick
}
