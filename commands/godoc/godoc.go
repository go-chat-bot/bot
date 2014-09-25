package godoc

import (
	"fmt"
	"github.com/fabioxgn/go-bot"
	"github.com/fabioxgn/go-bot/web"
	"net/url"
)

const (
	godocSiteURL    = "http://godoc.org"
	noPackagesFound = "No packages found."
)

var (
	godocSearchURL = "http://api.godoc.org/search"
)

type godocResults struct {
	Results []struct {
		Path     string `json:"path"`
		Synopsis string `json:"synopsis"`
	} `json:"results"`
}

func search(cmd *bot.Cmd) (string, error) {
	if cmd.FullArg == "" {
		return "", nil
	}

	data := &godocResults{}

	url, _ := url.Parse(godocSearchURL)
	q := url.Query()
	q.Set("q", cmd.FullArg)
	url.RawQuery = q.Encode()

	err := web.GetJSON(url.String(), data)
	if err != nil {
		return "", err
	}

	if len(data.Results) == 0 {
		return noPackagesFound, nil
	}

	return fmt.Sprintf("%s %s/%s", data.Results[0].Synopsis, godocSiteURL, data.Results[0].Path), nil
}

func init() {
	bot.RegisterCommand(
		"godoc",
		"Searchs godoc.org and displays the first result.",
		"package name",
		search)
}
