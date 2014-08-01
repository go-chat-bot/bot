// Package bot provides a simple to use IRC bot
package bot

import (
	"crypto/tls"
	"github.com/thoj/go-ircevent"
	"log"
	"math/rand"
	"strings"
	"time"
)

const (
	// CmdPrefix is the prefix used to identify a command.
	// !hello whould be identified as a command
	CmdPrefix = "!"
)

// Config must contain the necessary data to connect to an IRC server
type Config struct {
	Server        string   // IRC server:port. Ex: irc.freenode.org:7000
	Channels      []string // Channels to connect. Ex: []string{"#go-bot", "#channel mypassword"}
	User          string   // The IRC username the bot will use
	Nick          string   // The nick the bot will use
	Password      string   // Server password
	UseTLS        bool     // Should connect using TLS?
	TLSServerName string   // Must supply if UseTLS is true
	Debug         bool     // This will log all IRC communication to standad output
}

type ircConnection interface {
	Privmsg(target, message string)
	GetNick() string
	Join(target string)
	Part(target string)
}

var (
	irccon *irc.Connection
	config *Config
)

func onPRIVMSG(e *irc.Event) {
	messageReceived(e.Arguments[0], e.Message(), e.Nick, irccon)
}

func getServerName() string {
	separatorIndex := strings.LastIndex(config.Server, ":")
	if separatorIndex != -1 {
		return config.Server[:separatorIndex]
	} else {
		return config.Server
	}
}

func connect() {
	irccon = irc.IRC(config.User, config.Nick)
	irccon.Password = config.Password
	irccon.UseTLS = config.UseTLS
	irccon.TLSConfig = &tls.Config{
		ServerName: getServerName(),
	}
	irccon.VerboseCallbackHandler = config.Debug
	err := irccon.Connect(config.Server)
	if err != nil {
		log.Fatal(err)
	}
}

func onWelcome(e *irc.Event) {
	for _, channel := range config.Channels {
		irccon.Join(channel)
	}
}

func onEndOfNames(e *irc.Event) {
	log.Printf("onEndOfNames: %v", e.Arguments)
}

func configureEvents() {
	irccon.AddCallback("001", onWelcome)
	irccon.AddCallback("366", onEndOfNames)
	irccon.AddCallback("PRIVMSG", onPRIVMSG)
}

// Run reads the Config, connect to the specified IRC server and starts the bot.
// The bot will automatically join all the channels specified in the configuration
func Run(c *Config) {
	config = c
	connect()
	configureEvents()
	irccon.Loop()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
