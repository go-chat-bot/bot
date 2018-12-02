// Package slack implements Slack handlers for github.com/go-chat-bot/bot
package slack

import (
	"fmt"

	"github.com/go-chat-bot/bot"
	"github.com/nlopes/slack"
)

// MessageFilter allows implementing a filter function to transform the messages
// before sending to the channel, it is run before the bot sends the message to slack
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

const protocol = "slack"

func defaultMessageFilter(message string, _ *bot.User) (string, slack.PostMessageParameters) {
	return message, params
}

func responseHandler(target string, message string, sender *bot.User) {
	message, params := messageFilter(message, sender)
	api.PostMessage(target, message, params)
}

// FindUserBySlackID converts a slack.User into a bot.User struct
func FindUserBySlackID(userID string) *bot.User {
	slackUser, err := api.GetUserInfo(userID)
	if err != nil {
		fmt.Printf("Error retrieving slack user: %s\n", err)
		return &bot.User{
			ID:    userID,
			IsBot: false}
	}
	return &bot.User{
		ID:       userID,
		Nick:     slackUser.Name,
		RealName: slackUser.Profile.RealName,
		IsBot:    slackUser.IsBot}
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
	user := FindUserBySlackID(userID)
	if len(user.Nick) == 0 {
		user.IsBot = isBot
	}

	return user
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

// RunWithFilter executes the bot and sets up a message filter which will
// receive all the messages before they are sent to slack
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
	},
		&bot.Config{
			Protocol: protocol,
			Server:   teaminfo.Domain,
		},
	)

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
					go b.MessageReceived(
						&bot.ChannelData{
							Protocol:  "slack",
							Server:    teaminfo.Domain,
							Channel:   channel,
							HumanName: C.Name,
							IsPrivate: !C.IsChannel,
						},
						extractText(ev),
						extractUser(ev),
					)
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
