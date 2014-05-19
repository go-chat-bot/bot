package jira

import (
	"github.com/fabioxgn/go-bot"
	"os"
	"regexp"
	"strings"
)

const (
	pattern = "\\b[A-z]{3}-[0-9]{3,}\\b"
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

func getIssueURL(text string) (string, error) {
	issue := getIssue(text)
	if issue != "" {
		return (url + strings.ToUpper(issue)), nil
	}
	return "", nil
}

func jira(command *bot.PassiveCmd) (string, error) {
	return getIssueURL(command.Raw)
}

func init() {
	url = os.Getenv(env)
	bot.RegisterPassiveCommand(
		"jira",
		jira)
}
