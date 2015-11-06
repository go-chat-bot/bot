package slack

import (
	"fmt"

	"github.com/go-chat-bot/bot"
	"github.com/nlopes/slack"
)

var (
	rtm *slack.RTM
)

func responseHandler(target string, message string, sender *bot.User) {
	rtm.SendMessage(rtm.NewOutgoingMessage(message, target))
}

// Run connects to slack RTM API using the provided token
func Run(token string) {
	api := slack.New(token)
	rtm = api.NewRTM()

	bot.New(&bot.Handlers{
		Response: responseHandler,
	})

	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				bot.MessageReceived(ev.Channel, ev.Text, &bot.User{Nick: ev.User})

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop
			}
		}
	}
}
