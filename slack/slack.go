// Package slack implements Slack handlers for github.com/go-chat-bot/bot
package slack

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/go-chat-bot/bot"
	"github.com/slack-go/slack"
)

// MessageFilter allows implementing a filter function to transform the messages
// before sending to the channel, it is run before the bot sends the message to slack
type MessageFilter func(string, *bot.User) (string, slack.PostMessageParameters)

// Strip Slack's mrkdwn URL
var mrkdwnURLRegexp = regexp.MustCompile(`\<http:\/\/.*\|(.*)\>`)

var (
	rtm      *slack.RTM
	api      *slack.Client
	teaminfo *slack.TeamInfo
	verbose  bool

	channelList                 = map[string]slack.Channel{}
	params                      = slack.PostMessageParameters{AsUser: true}
	messageFilter MessageFilter = defaultMessageFilter
	botUserID                   = ""
)

const protocol = "slack"

// Slack channel names have very specific prefixes:
// # = channel name, C = channel ID, G = group message, D = direct message
var slack_channel_regexp = regexp.MustCompile("^[#CGD]")

func defaultMessageFilter(message string, _ *bot.User) (string, slack.PostMessageParameters) {
	return message, params
}

func responseHandler(target string, message string, sender *bot.User) {
	if !(slack_channel_regexp.MatchString(target)) {
		// log.Printf("Slack message target does not appear to be a Slack channel.")
		return
	}

	message, params := messageFilter(message, sender)

	if verbose {
		nick := "botpassivecmd"
		if sender != nil {
			nick = sender.Nick
		}
		log.Printf("[%s] (outgoing) channel: %s reply-to: %s message: %s", protocol, whereMessage(target), nick, message)
	}
	_, _, err := api.PostMessage(
		target,
		slack.MsgOptionPostMessageParameters(params),
		slack.MsgOptionText(message, false),
	)
	if err != nil {
		fmt.Printf("Error sending a slack message: %s\n", err.Error())
	}
}

func responseHandlerV2(om bot.OutgoingMessage) {
	if !(slack_channel_regexp.MatchString(om.Target)) {
		// log.Printf("Slack messageV2 target does not appear to be a Slack channel.")
		return
	}

	message, params := messageFilter(om.Message, om.Sender)
	if pmp, ok := om.ProtoParams.(*slack.PostMessageParameters); ok {
		params = *pmp
	}
	if verbose {
		nick := "botpassivecmd"
		if om.Sender != nil {
			nick = om.Sender.Nick
		}
		log.Printf("[%s] (outgoing) channel: %s reply-to: %s message: %s", protocol, whereMessage(om.Target), nick, message)
	}
	_, _, err := api.PostMessage(
		om.Target,
		slack.MsgOptionPostMessageParameters(params),
		slack.MsgOptionText(message, false),
	)
	if err != nil {
		fmt.Printf("Error sending a slack message: %s\n", err.Error())
	}
}

// AddReactionToMessage allows you to add a reaction, to a message.
func AddReactionToMessage(msgid, channel string, reaction string) error {
	toReact := slack.ItemRef{
		Timestamp: msgid,
		Channel:   channel,
	}

	return api.AddReaction(reaction, toReact)
}

// RemoveReactionFromMessage allows you to remove a reaction, from a message.
func RemoveReactionFromMessage(msgid, channel string, reaction string) error {
	reactionRef := slack.ItemRef{
		Timestamp: msgid,
		Channel:   channel,
	}

	return api.RemoveReaction(reaction, reactionRef)
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
	msg := &bot.Message{
		ProtoMsg: event,
	}
	if len(event.Text) != 0 {
		msg.Text = mrkdwnURLRegexp.ReplaceAllString(event.Text, "$1")
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
	params := slack.GetConversationsParameters{}
	channels, _, err := api.GetConversations(&params)
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

func whereMessage(channel string) string {
	if strings.HasPrefix(channel, "#") {
		// this appears to be a full channel name already, no modifications
		return channel
	}
	// this most likelys is a channel ID
	C, err := api.GetConversationInfo(channel, false)
	if err == nil && C != nil {
		if C.IsIM {
			return "privatemsg"
		}
		return "#" + C.Name
	}
	return "not-a-slack-channel"
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
	_, verbose = os.LookupEnv("SLACK_VERBOSE")
	_, apiDebug := os.LookupEnv("SLACK_DEBUG")
	api = slack.New(token, slack.OptionDebug(apiDebug))
	rtm = api.NewRTM()
	teaminfo, _ = api.GetTeamInfo()

	b := bot.New(&bot.Handlers{
		Response:   responseHandler,
		ResponseV2: responseHandlerV2,
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
				if ev.Hidden || ownMessage(ev.User) {
					continue
				}

				var channel = ev.Channel
				C, _ := api.GetConversationInfo(channel, false)
				if C.IsChannel {
					channel = fmt.Sprintf("#%s", C.Name)
				}
				text := extractText(ev)
				user := extractUser(ev)
				if verbose && strings.HasPrefix(text.Text, bot.CmdPrefix) {
					// logs incoming message only if verbose is ON and the message starts with the bot's command prefix
					log.Printf("[%s] (incoming) channel: %s from: %s message: %s", protocol, whereMessage(C.ID), user.Nick, text.Text)
				}
				go b.MessageReceived(
					&bot.ChannelData{
						Protocol:  protocol,
						Server:    teaminfo.Domain,
						Channel:   channel,
						HumanName: C.Name,
						IsPrivate: !C.IsChannel,
					},
					text,
					user,
				)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop
			}
		}
	}
}
