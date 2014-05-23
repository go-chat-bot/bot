package jira

import (
	"github.com/fabioxgn/go-bot"
	"os"
	"regexp"
	"strings"
)

const (
	pattern = "\\b[A-z]{3}-[0-9]+\\b"
	env     = "JIRA_ISSUES_URL"
)

var (
	url string
)

func getIssue(text string) string {
	if match, _ := regexp.MatchString(pattern, text); match {
		return regexp.MustCompile(pattern).FindString(text)
	}
	return ""
}

func getIssueURL(nick, text string) (string, error) {
	if !strings.Contains(nick, "bot") {
		issue := getIssue(text)
		if issue != "" {
			return (url + strings.ToUpper(issue)), nil
		}
	}
	return "", nil
}

func jira(command *bot.PassiveCmd) (string, error) {
	return getIssueURL(command.Nick, command.Raw)
}

func init() {
	url = os.Getenv(env)
	bot.RegisterPassiveCommand(
		"jira",
		jira)
}
