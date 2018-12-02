package rocket

import (
	"fmt"
	"log"

	bot "github.com/bnfinet/go-chat-bot"
	"github.com/pyinx/gorocket/api"
	"github.com/pyinx/gorocket/rest"
)

var (
	client *rest.Client
	config *Config
)

const (
	protocol = "rocket"
)

// Config must contain the necessary data to connect to an rocket.chat server
type Config struct {
	Server   string
	Port     string
	User     string
	Email    string
	Password string
	UseTLS   bool
	Debug    bool
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

// Run reads the Config, connect to the specified rocket.chat server and starts the bot.
func Run(c *Config) {
	config = c
	client = rest.NewClient(config.Server, config.Port, config.UseTLS, config.Debug)
	err := client.Login(api.UserCredentials{Email: config.Email, Name: config.User, Password: config.Password})
	if err != nil {
		log.Fatalf("login err: %s\n", err)
	}

	b := bot.New(&bot.Handlers{
		Response: responseHandler,
	},
		&bot.Config{
			Protocol: protocol,
			Server:   config.Server,
		},
	)

	b.Disable([]string{"url"})

	msgChan := client.GetAllMessages()
	for {
		select {
		case msgs := <-msgChan:
			for _, msg := range msgs {
				if !ownMessage(c, msg) {
					b.MessageReceived(
						&bot.ChannelData{
							Protocol:  protocol,
							Server:    "",
							Channel:   msg.ChannelId,
							IsPrivate: false,
						},
						&bot.Message{Text: msg.Text},
						&bot.User{ID: msg.User.Id, RealName: msg.User.UserName, Nick: msg.User.UserName, IsBot: false})
				}

			}
		}
	}
}
