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

	params    = slack.PostMessageParameters{AsUser: true}
	botUserID = ""
)

func responseHandler(target string, message string, sender *bot.User) {
	api.PostMessage(target, message, params)
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

func readBotInfo(api *slack.Client) {
	info, err := api.AuthTest()
	if err != nil {
		fmt.Printf("Error calling AuthTest: %s\n", err)
		return
	}
	botUserID = info.UserID
}

func ownMessage(UserID string) bool {
	return botUserID == UserID
}

// Run connects to slack RTM API using the provided token
func Run(token string) {
	api = slack.New(token)
	rtm = api.NewRTM()

	b := bot.New(&bot.Handlers{
		Response: responseHandler,
	})
	b.Disable([]string{"url"})

	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				readBotInfo(api)
			case *slack.MessageEvent:
				if !ev.Hidden && !ownMessage(ev.User) {
					b.MessageReceived(ev.Channel, ev.Text, extractUser(ev.User))
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop
			}
		}
	}
}
