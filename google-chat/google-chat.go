package googlechat

import (
	"bytes"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/go-chat-bot/bot"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	chatAuthScope = "https://www.googleapis.com/auth/chat.bot"
	apiEndpoint   = "https://chat.googleapis.com/v1/"
)

var (
	httpChatClient *http.Client
	b              *bot.Bot
)

// Config must contain basic configuration for the bot to be able to work
type Config struct {
	PubSubProject    string
	TopicName        string
	SubscriptionName string
	Token            string
	WelcomeMessage   string
}

func responseHandler(target string, message string, sender *bot.User) {
	var space, thread string

	// this define thread in the reply if we can so we don't alwayus start new
	targets := strings.Split(target, ":")
	if len(targets) < 2 {
		space = target
	} else {
		space = targets[0]
		thread = targets[1]
	}

	reply, err := json.Marshal(&ReplyMessage{
		Text: message,
		Thread: &ReplyThread{
			Name: thread}})

	log.Printf("Replying: space: %s thread: %s message: %s\n",
		space, thread, message)
	resp, err := httpChatClient.Post(apiEndpoint+space+"/messages",
		"application/json",
		bytes.NewReader(reply))
	if err != nil {
		log.Printf("Error posting reply: %v", err)
	}

	log.Printf("Response: %s\n", resp.Status)
}

// Run reads the config, establishes OAuth connection & Pub/Sub subscription to
// the message queue
func Run(config *Config) {
	var err error
	ctx := context.Background()
	httpChatClient, err = google.DefaultClient(ctx, chatAuthScope)
	if err != nil {
		log.Printf("Error setting http client: %v\n", err)
		return
	}

	client, err := pubsub.NewClient(ctx, config.PubSubProject)
	if err != nil {
		log.Printf("Error creating client: %v\n", err)
		return
	}

	topic := client.Topic(config.TopicName)

	// Create a new subscription to the previously created topic
	// with the given name.
	sub := client.Subscription(config.SubscriptionName)
	ok, err := sub.Exists(ctx)
	if err != nil {
		log.Printf("Error getting subscription: %v\n", err)
		return
	}
	if !ok {
		// Subscription doesn't exist.
		sub, err = client.CreateSubscription(ctx, config.SubscriptionName,
			pubsub.SubscriptionConfig{
				Topic:       topic,
				AckDeadline: 10 * time.Second,
			})
		if err != nil {
			log.Printf("Error subscribing: %v\n", err)
			return
		}
	}

	b = bot.New(&bot.Handlers{
		Response: responseHandler,
	})

	err = sub.Receive(context.Background(),
		func(ctx context.Context, m *pubsub.Message) {
			var msg ChatMessage
			err = json.Unmarshal(m.Data, &msg)
			if err != nil {
				log.Printf("Failed message unmarshal(%v): %s\n", err, m.Data)
				m.Ack()
				return
			}
			if msg.Token != config.Token {
				log.Printf("Failed to verify token: %s", msg.Token)
				m.Ack()
				return
			}

			log.Printf("Space: %s (%s)\n", msg.Space.Name, msg.Space.DisplayName)
			log.Printf("Message type: %s\n", msg.Type)
			log.Printf("From: %s (%s)\n", msg.User.Name, msg.User.DisplayName)
			switch msg.Type {
			case "ADDED_TO_SPACE":
				if config.WelcomeMessage != "" {
					log.Printf("Sending welcome message to %s\n", msg.Space.Name)
					b.SendMessage(msg.Space.Name, config.WelcomeMessage, nil)
				}
			case "REMOVED_FROM_SPACE":
				break
			case "MESSAGE":
				log.Printf("Message: %s\n", msg.Message.ArgumentText)
				b.MessageReceived(
					&bot.ChannelData{
						Protocol:  "googlechat",
						Server:    "chat.google.com",
						HumanName: msg.Space.DisplayName,
						Channel:   msg.Space.Name + ":" + msg.Message.Thread.Name,
						IsPrivate: msg.Space.Type == "DM",
					},
					&bot.Message{
						Text:     msg.Message.ArgumentText,
						IsAction: false,
					},
					&bot.User{
						ID:       msg.User.Name,
						Nick:     msg.User.DisplayName,
						RealName: msg.User.DisplayName,
					})
			}

			m.Ack()
		})
	if err != nil {
		log.Printf("Error setting up receiving: %v\n", err)
		return
	}
	// Wait indefinetely
	select {}
}
