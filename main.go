package main

import (
	"fmt"
	"github.com/fabioxgn/go-bot/commands"
	"github.com/thoj/go-ircevent"
	"github.com/yvasiyarov/gorelic"
	"log"
	"os"
)

const (
	configFile = "config.json"
)

var (
	irccon *irc.Connection
	config = &Config{}
)

func isPrivateMsg(channel string) bool {
	return channel == config.Nick
}

func onPRIVMSG(e *irc.Event) {
	channel := e.Arguments[0]

	if isPrivateMsg(channel) {
		channel = e.Nick //e.Nick is who sent the pvt message
	}

	cmd := commands.Parse(e.Message, config.CmdPrefix, channel, e.Nick)
	if cmd.IsCommand {
		commands.HandleCmd(cmd, irccon)
	} else {
		// It's not a command
		// TODO: Test for passive commands (parse url, etc) ?
	}
}

func connect() {
	irccon = irc.IRC(config.User, config.Nick)
	irccon.UseTLS = config.UseTLS
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

func readConfig() {
	file, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	config.Read(file)
	fmt.Printf("%v\n", config)
}

func startMetrics() {
	newRelicAPIKey := os.Getenv("NEW_RELIC_KEY")

	if newRelicAPIKey == "" {
		return
	}

	agent := gorelic.NewAgent()
	agent.Verbose = false
	agent.NewrelicLicense = newRelicAPIKey
	agent.Run()
}

func main() {
	startMetrics()
	readConfig()
	connect()
	configureEvents()
	irccon.Loop()
}
