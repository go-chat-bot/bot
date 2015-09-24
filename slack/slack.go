package slack

import (
	"fmt"

	"github.com/go-chat-bot/bot"
	"github.com/nlopes/slack"
)

type slackConnection struct {
	rtm *slack.RTM
}

func (c slackConnection) MessageReceived(target, message, sender string) {
	c.rtm.SendMessage(c.rtm.NewOutgoingMessage(message, target))
}

// RunSlack connects to slack RTM API using the provided token
func Run(token string) {
	api := slack.New(token)

	conn := new(slackConnection)
	conn.rtm = api.NewRTM()

	gobot := bot.NewBot(conn.MessageReceived, []string{})

	go conn.rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-conn.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				gobot.MessageReceived(ev.Channel, ev.Text, ev.User)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop
			}
		}
	}
}
