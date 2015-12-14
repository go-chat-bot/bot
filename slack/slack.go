// Package slack implements Slack handlers for github.com/go-chat-bot/bot
package slack

import (
	"fmt"

	"github.com/go-chat-bot/bot"
	"github.com/nlopes/slack"
)

var (
	rtm *slack.RTM
	api *slack.Client
)

func responseHandler(target string, message string, sender *bot.User) {
	rtm.SendMessage(rtm.NewOutgoingMessage(message, target))
}

// Extracts user information from slack API
func extractUser(userID string) *bot.User {
	slackUser, err := api.GetUserInfo(userID)
	if err != nil {
		fmt.Printf("Error retrieving slack user: %s\n", err)
		return &bot.User{Nick: userID}
	}
	return &bot.User{Nick: slackUser.Name, RealName: slackUser.Profile.RealName}
}

// Run connects to slack RTM API using the provided token
func Run(token string) {
	api = slack.New(token)
	rtm = api.NewRTM()

	b := bot.New(&bot.Handlers{
		Response: responseHandler,
	})

	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				b.MessageReceived(ev.Channel, ev.Text, extractUser(ev.User))

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop
			}
		}
	}
}
