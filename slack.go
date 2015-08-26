package bot

import (
	"fmt"

	"github.com/nlopes/slack"
)

type slackConnection struct {
	rtm *slack.RTM
}

func (c *slackConnection) GetNick() string {
	return ""
}
func (c slackConnection) Join(channel string) {}
func (c slackConnection) Part(channel string) {}

func (c slackConnection) Privmsg(target, message string) {
	c.rtm.SendMessage(c.rtm.NewOutgoingMessage(message, target))
}

// RunSlack connects to slack RTM API using the provided token
func RunSlack(token string) {
	api := slack.New(token)

	conn := new(slackConnection)
	conn.rtm = api.NewRTM()
	go conn.rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-conn.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				messageReceived(ev.Channel, ev.Text, ev.User, conn)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop
			}
		}
	}
}
