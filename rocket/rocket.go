package rocket

import (
	"fmt"
	"github.com/go-chat-bot/bot"
	"github.com/pyinx/gorocket/api"
	"github.com/pyinx/gorocket/rest"
	"log"
)

var (
	client *rest.Client
	config *Config
)

type Config struct {
	Server   string
	Port     string
	User     string
	Email    string
	Password string
	// Channels []string
	UseTLS bool
	Debug  bool
}

func responseHandler(target string, message string, sender *bot.User) {
	atUser := sender.RealName
	if message == "" {
		return
	}
	message = fmt.Sprintf("@%s %s", atUser, message)
	channelInfo, _ := client.GetChannelInfoById(target)
	err := client.Send(channelInfo, message)
	if err != nil {
		if config.Debug {
			log.Printf("send message err: %s\n", err)
		}
	}
}

func ownMessage(c *Config, msg api.Message) bool {
	return c.User == msg.User.UserName
}

func Run(c *Config) {
	config = c
	client = rest.NewClient(config.Server, config.Port, config.UseTLS, config.Debug)
	err := client.Login(api.UserCredentials{Email: config.Email, Name: config.User, Password: config.Password})
	if err != nil {
		log.Fatal("login err: %s\n", err)
	}

	b := bot.New(&bot.Handlers{
		Response: responseHandler,
	})
	b.Disable([]string{"url"})

	msgChan := client.GetAllMessages()
	for {
		select {
		case msgs := <-msgChan:
			for _, msg := range msgs {
				if !ownMessage(c, msg) {
					b.MessageReceived(
						&bot.ChannelData{
							Protocol:  "rocket",
							Server:    "",
							Channel:   msg.ChannelId,
							IsPrivate: false,
						},
						&bot.Message{Text: msg.Text, IsAction: true},
						&bot.User{ID: msg.User.Id, RealName: msg.User.UserName, Nick: msg.User.UserName, IsBot: false})
				}

			}
		}
	}
}
