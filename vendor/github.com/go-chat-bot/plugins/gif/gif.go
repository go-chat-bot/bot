package gif

import (
	"fmt"
	"github.com/go-chat-bot/bot"
	"github.com/go-chat-bot/plugins/web"
	"math/rand"
	"net/url"
	"time"
)

const (
	giphyURL = "http://api.giphy.com/v1/gifs/search?q=%s&api_key=dc6zaTOxFJmzC&limit=50"
)

type giphy struct {
	Data []struct {
		BitlyURL string `json:"bitly_url"`
		Images   struct {
			FixedHeight struct {
				Height string `json:"height"`
				URL    string `json:"url"`
				Width  string `json:"width"`
			} `json:"fixed_height"`
			FixedHeightDownsampled struct {
				Height string `json:"height"`
				URL    string `json:"url"`
				Width  string `json:"width"`
			} `json:"fixed_height_downsampled"`
			FixedHeightStill struct {
				Height string `json:"height"`
				URL    string `json:"url"`
				Width  string `json:"width"`
			} `json:"fixed_height_still"`
			FixedWidth struct {
				Height string `json:"height"`
				URL    string `json:"url"`
				Width  string `json:"width"`
			} `json:"fixed_width"`
			FixedWidthDownsampled struct {
				Height string `json:"height"`
				URL    string `json:"url"`
				Width  string `json:"width"`
			} `json:"fixed_width_downsampled"`
			FixedWidthStill struct {
				Height string `json:"height"`
				URL    string `json:"url"`
				Width  string `json:"width"`
			} `json:"fixed_width_still"`
			Original struct {
				Frames string `json:"frames"`
				Height string `json:"height"`
				Size   string `json:"size"`
				URL    string `json:"url"`
				Width  string `json:"width"`
			} `json:"original"`
		} `json:"images"`
		Type        string `json:"type"`
		Username    string `json:"username"`
		BitlyGifURL string `json:"bitly_gif_url"`
		EmbedURL    string `json:"embed_url"`
		ID          string `json:"id"`
		Rating      string `json:"rating"`
		Source      string `json:"source"`
		URL         string `json:"url"`
	} `json:"data"`
	Meta struct {
		Msg    string `json:"msg"`
		Status int64  `json:"status"`
	} `json:"meta"`
	Pagination struct {
		Count      int64 `json:"count"`
		Offset     int64 `json:"offset"`
		TotalCount int64 `json:"total_count"`
	} `json:"pagination"`
}

func gif(command *bot.Cmd) (msg string, err error) {
	data := &giphy{}
	err = web.GetJSON(fmt.Sprintf(giphyURL, url.QueryEscape(command.RawArgs)), data)
	if err != nil {
		return "", err
	}

	if len(data.Data) == 0 {
		return "No gifs found. try: !gif cat", nil
	}

	index := rand.Intn(len(data.Data))
	return fmt.Sprintf(data.Data[index].Images.FixedHeight.URL), nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
	bot.RegisterCommand(
		"gif",
		"Searchs and posts a random gif url from Giphy.",
		"cat",
		gif)
}
