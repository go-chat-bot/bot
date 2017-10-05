// Package slack implements Slack handlers for github.com/go-chat-bot/bot
package slack

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-chat-bot/bot"
	"github.com/nlopes/slack"
)

type MessageFilter func(string, *bot.User) (string, slack.PostMessageParameters)

var (
	rtm      *slack.RTM
	api      *slack.Client
	teaminfo *slack.TeamInfo

	channelList                 = map[string]slack.Channel{}
	params                      = slack.PostMessageParameters{AsUser: true}
	messageFilter MessageFilter = defaultMessageFilter
	botUserID                   = ""
)

func defaultMessageFilter(message string, sender *bot.User) (string, slack.PostMessageParameters) {
	return message, params
}

func responseHandler(target string, message string, sender *bot.User) {
	message, params := messageFilter(message, sender)
	api.PostMessage(target, message, params)
}

// Extracts user information from slack API
func extractUser(event *slack.MessageEvent) *bot.User {
	var isBot bool
	var userID string
	if len(event.User) == 0 {
		userID = event.BotID
		isBot = true
	} else {
		userID = event.User
		isBot = false
	}
	slackUser, err := api.GetUserInfo(userID)
	if err != nil {
		fmt.Printf("Error retrieving slack user: %s\n", err)
		return &bot.User{
			ID:    userID,
			IsBot: isBot}
	}
	return &bot.User{
		ID:       userID,
		Nick:     slackUser.Name,
		RealName: slackUser.Profile.RealName,
		IsBot:    isBot}
}

func parse(msgtxt string) string {
	regx := regexp.MustCompile("<[^<>]+>")
	matches := regx.FindAllString(msgtxt, -1)
	for _, match := range matches {
		clean := strings.Split(match, "|")[1]
		clean = clean[:len(clean)-1]
		msgtxt = strings.Replace(msgtxt, match, clean, -1)
	}
	return msgtxt
}

func extractText(event *slack.MessageEvent) *bot.Message {
	msg := &bot.Message{}
	if len(event.Text) != 0 {
		msg.Text = event.Text
		if event.SubType == "me_message" {
			msg.IsAction = true
		}
	} else {
		attachments := event.Attachments
		if len(attachments) > 0 {
			msg.Text = attachments[0].Fallback
		}
	}
	msg.Text = parse(msg.Text)
	return msg
}

func readBotInfo(api *slack.Client) {
	info, err := api.AuthTest()
	if err != nil {
		fmt.Printf("Error calling AuthTest: %s\n", err)
		return
	}
	botUserID = info.UserID
}

func readChannelData(api *slack.Client) {
	channels, err := api.GetChannels(true)
	if err != nil {
		fmt.Printf("Error getting Channels: %s\n", err)
		return
	}
	for _, channel := range channels {
		channelList[channel.ID] = channel
	}
}

func ownMessage(UserID string) bool {
	return botUserID == UserID
}

func RunWithFilter(token string, customMessageFilter MessageFilter) {
	if customMessageFilter == nil {
		panic("A valid message filter must be provided.")
	}
	messageFilter = customMessageFilter
	Run(token)
}

// Run connects to slack RTM API using the provided token
func Run(token string) {
	api = slack.New(token)
	rtm = api.NewRTM()
	teaminfo, _ = api.GetTeamInfo()

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
				readChannelData(api)
			case *slack.ChannelCreatedEvent:
				readChannelData(api)
			case *slack.ChannelRenameEvent:
				readChannelData(api)

			case *slack.MessageEvent:
				if !ev.Hidden && !ownMessage(ev.User) {
					C := channelList[ev.Channel]
					var channel = ev.Channel
					if C.IsChannel {
						channel = fmt.Sprintf("#%s", C.Name)
					}
					b.MessageReceived(
						&bot.ChannelData{
							Protocol:  "slack",
							Server:    teaminfo.Domain,
							Channel:   channel,
							IsPrivate: !C.IsChannel,
						},
						extractText(ev),
						extractUser(ev))
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
