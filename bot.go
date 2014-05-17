package bot

import (
	"github.com/thoj/go-ircevent"
	"log"
)

const (
	CmdPrefix = "!" // Prefix used to identify a command. This is identified as a command: !hello
)

// Config must contain the necessary data to connect to an IRC server
type Config struct {
	Server   string   // IRC server:port. Ex: irc.freenode.org:7000
	Channels []string // Channels to connect. Ex: []string{"#go-bot", "#channel mypassword"}
	User     string   // The IRC username the bot will use
	Nick     string   // The nick the bot will use
	Password string   // Server password
	UseTLS   bool     // Should connect using TLS?
	Debug    bool     // This will log all IRC communication to standad output
}

type ircConnection interface {
	Privmsg(target, message string)
	GetNick() string
}

var (
	irccon *irc.Connection
	config *Config
)

func onPRIVMSG(e *irc.Event) {
	messageReceived(e.Arguments[0], e.Message(), e.Nick, irccon)
}

func connect() {
	irccon = irc.IRC(config.User, config.Nick)
	irccon.Password = config.Password
	irccon.UseTLS = config.UseTLS
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
	log.Println("onEndOfNames: %v", e.Arguments)
	irccon.Privmsg(e.Arguments[1], "Hi there.\n")
}

func configureEvents() {
	irccon.AddCallback("001", onWelcome)
	irccon.AddCallback("366", onEndOfNames)
	irccon.AddCallback("PRIVMSG", onPRIVMSG)
}

// Run reads the Config, connect to IRC and starts the bot
func Run(c *Config) {
	config = c
	connect()
	configureEvents()
	irccon.Loop()
}
