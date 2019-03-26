package gag

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/go-chat-bot/bot"
)

const (
	randomURL = "http://9gag.com/random"
)

func randomPage(command *bot.Cmd) (string, error) {
	var redirectNotAllowed = errors.New("redirect")
	redirectedURL := ""

	client := http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		redirectedURL = req.URL.String()
		return redirectNotAllowed
	}

	_, err := client.Get(randomURL)
	if urlError, ok := err.(*url.Error); ok && urlError.Err == redirectNotAllowed {
		return redirectedURL, nil
	}
	return "", err
}

func init() {
	bot.RegisterCommand(
		"9gag",
		"Returns a random 9gag page.",
		"",
		randomPage)
}
