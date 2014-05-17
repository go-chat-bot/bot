package bot

import (
	"github.com/thoj/go-ircevent"
	"log"
)

const (
	CmdPrefix = "!"
)

// TODO
type Config struct {
	Server   string
	Channels []string
	User     string
	Nick     string
	Password string
	UseTLS   bool
	Debug    bool
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

func Run(c *Config) {
	config = c
	connect()
	configureEvents()
	irccon.Loop()
}
