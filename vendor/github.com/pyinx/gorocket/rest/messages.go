package rest

import (
	"bytes"
	"fmt"
	"github.com/pyinx/gorocket/api"
	"html"
	"log"
	"net/http"
	"strings"
	"time"
)

type messagesResponse struct {
	statusResponse
	ChannelName string        `json:"channel"`
	Messages    []api.Message `json:"messages"`
}

type messageResponse struct {
	statusResponse
	ChannelName string      `json:"channel"`
	Message     api.Message `json:"message"`
}

// Sends a message to a channel. The name of the channel has to be not nil.
// The message will be html escaped.
//
// https://rocket.chat/docs/developer-guides/rest-api/chat/postmessage
func (c *Client) Send(channel *api.Channel, msg string) error {
	msg = strings.Replace(msg, "\n", "\\n", -1)
	body := fmt.Sprintf(`{ "channel": "%s", "text": "%s"}`, channel.Name, html.EscapeString(msg))
	request, _ := http.NewRequest("POST", c.getUrl()+"/api/v1/chat.postMessage", bytes.NewBufferString(body))

	response := new(messageResponse)

	return c.doRequest(request, response)
}

// Get messages from a channel. The channel id has to be not nil. Optionally a
// count can be specified to limit the size of the returned messages.
//
// https://rocket.chat/docs/developer-guides/rest-api/channels/history
func (c *Client) GetMessagesOnce(channel *api.Channel, lastTime string) ([]api.Message, error) {
	u := fmt.Sprintf("%s/api/v1/channels.history?roomId=%s", c.getUrl(), channel.Id)

	if lastTime != "" {
		u = fmt.Sprintf("%s&oldest=%s", u, lastTime)
	}

	request, _ := http.NewRequest("GET", u, nil)
	response := new(messagesResponse)

	if err := c.doRequest(request, response); err != nil {
		return nil, err
	}

	return response.Messages, nil
}

func (c *Client) GetMessages(channel *api.Channel, lastTime string, msgChan chan []api.Message) {
	for {
		u := fmt.Sprintf("%s/api/v1/channels.history?roomId=%s", c.getUrl(), channel.Id)

		if lastTime != "" {
			u = fmt.Sprintf("%s&oldest=%s", u, lastTime)
		}
		request, _ := http.NewRequest("GET", u, nil)
		response := new(messagesResponse)

		if err := c.doRequest(request, response); err != nil {
			log.Printf("ERROR: get message from channel err: %s\n", err)
			time.Sleep(200 * time.Microsecond)
			continue
		}
		if len(response.Messages) != 0 {
			lastTime = response.Messages[0].Timestamp
			msgChan <- response.Messages
		}
		time.Sleep(200 * time.Microsecond)
	}
}

func (c *Client) GetMuliMessages(channels []api.Channel) chan []api.Message {
	msgChan := make(chan []api.Message, len(channels))
	go func() {
		msgMap := make(map[string]string)
		for {
			for _, channel := range channels {
				lastTime := msgMap[channel.Name]
				msg, err := c.GetMessagesOnce(&channel, lastTime)
				if err != nil {
					log.Printf("ERROR: get message from channel %s err: %s\n", channel.Name, err)
				} else {
					if len(msg) != 0 {
						msgMap[channel.Name] = msg[0].Timestamp
						msgChan <- msg
					}
				}
			}
			time.Sleep(200 * time.Microsecond)
		}
	}()
	return msgChan
}

func (c *Client) GetAllMessages() chan []api.Message {
	channels, _ := c.GetJoinedChannels()
	msgChan := make(chan []api.Message, 1024)
	go func() {
		msgMap := make(map[string]string)
		for {
			for _, channel := range channels {
				lastTime := msgMap[channel.Name]
				msg, err := c.GetMessagesOnce(&channel, lastTime)
				if err != nil {
					log.Printf("ERROR: get message from channel %s err: %s\n", channel.Name, err)
				} else {
					if len(msg) != 0 {
						msgMap[channel.Name] = msg[0].Timestamp
						msgChan <- msg
					}
				}
			}
			time.Sleep(200 * time.Microsecond)
			channels, _ = c.GetJoinedChannels()
		}
	}()
	return msgChan
}
