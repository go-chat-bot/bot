// Package irc implements IRC handlers for github.com/go-chat-bot/bot
package irc

import (
	"crypto/tls"
	"log"
	"strings"

	"github.com/go-chat-bot/bot"
	ircevent "github.com/thoj/go-ircevent"
)

// Config must contain the necessary data to connect to an IRC server
type Config struct {
	Server        string   // IRC server:port. Ex: ircevent.freenode.org:7000
	Channels      []string // Channels to connect. Ex: []string{"#go-bot", "#channel mypassword"}
	User          string   // The IRC username the bot will use
	Nick          string   // The nick the bot will use
	Password      string   // Server password
	UseTLS        bool     // Should connect using TLS?
	TLSServerName string   // Must supply if UseTLS is true
	Debug         bool     // This will log all IRC communication to standad output
}

var (
	ircConn *ircevent.Connection
	config  *Config
	b       *bot.Bot
)

func responseHandler(target string, message string, sender *bot.User) {
	channel := target
	if ircConn.GetNick() == target {
		channel = sender.Nick
	}

	for _, line := range strings.Split(message, "\n") {
		ircConn.Privmsg(channel, line)
	}
}

func onPRIVMSG(e *ircevent.Event) {
	b.MessageReceived(
		&bot.ChannelData{
			Protocol:  "irc",
			Server:    ircConn.Server,
			Channel:   e.Arguments[0],
			IsPrivate: e.Arguments[0] == ircConn.GetNick()},
		&bot.Message{
			Text: e.Message(),
		},
		&bot.User{
			ID:       e.Host,
			Nick:     e.Nick,
			RealName: e.User})
}

func onCTCPACTION(e *ircevent.Event) {
	b.MessageReceived(
		&bot.ChannelData{
			Protocol:  "irc",
			Server:    ircConn.Server,
			Channel:   e.Arguments[0],
			IsPrivate: false,
		},
		&bot.Message{
			Text:     e.Message(),
			IsAction: true,
		},
		&bot.User{
			ID:       e.Host,
			Nick:     e.Nick,
			RealName: e.User,
		})
}

func getServerName(server string) string {
	separatorIndex := strings.LastIndex(server, ":")
	if separatorIndex != -1 {
		return server[:separatorIndex]
	}
	return server
}

func onWelcome(e *ircevent.Event) {
	for _, channel := range config.Channels {
		ircConn.Join(channel)
	}
}

// Run reads the Config, connect to the specified IRC server and starts the bot.
// The bot will automatically join all the channels specified in the configuration
func Run(c *Config) {
	config = c

	ircConn = ircevent.IRC(c.User, c.Nick)
	ircConn.Password = c.Password
	ircConn.UseTLS = c.UseTLS
	ircConn.TLSConfig = &tls.Config{
		ServerName: getServerName(c.Server),
	}
	ircConn.VerboseCallbackHandler = c.Debug

	b = bot.New(&bot.Handlers{
		Response: responseHandler,
	})

	ircConn.AddCallback("001", onWelcome)
	ircConn.AddCallback("PRIVMSG", onPRIVMSG)
	ircConn.AddCallback("CTCP_ACTION", onCTCPACTION)

	err := ircConn.Connect(c.Server)
	if err != nil {
		log.Fatal(err)
	}
	ircConn.Loop()
}
