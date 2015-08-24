package bot

import (
	"fmt"

	"github.com/nlopes/slack"
)

// SlackConnection implements the basic ircConnection interface
type SlackConnection struct {
	rtm *slack.RTM
}

func (c *SlackConnection) GetNick() string {
	return ""
}

func (c SlackConnection) Join(channel string) {}

func (c SlackConnection) Part(channel string) {}

func (c SlackConnection) Privmsg(target, message string) {
	c.rtm.SendMessage(c.rtm.NewOutgoingMessage(message, target))
}

// RunSlack connects to slack RTM API using the provided token
func RunSlack(token string) {
	api := slack.New(token)
	api.SetDebug(true)

	conn := new(SlackConnection)
	conn.rtm = api.NewRTM()
	go conn.rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-conn.rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {

			case *slack.ConnectedEvent:
				fmt.Println("Infos:", ev.Info)
				fmt.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)
				messageReceived(ev.Channel, ev.Text, ev.User, conn)

			case *slack.PresenceChangeEvent:
				fmt.Printf("Presence Change: %v\n", ev)

			case *slack.LatencyReport:
				fmt.Printf("Current latency: %v\n", ev.Value)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}
