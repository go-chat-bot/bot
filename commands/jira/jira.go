package jira

import (
	"github.com/fabioxgn/go-bot"
	"os"
	"regexp"
	"strings"
)

const (
	pattern = "(^|\\s)+[A-z]{3}-[0-9]+\\b"
	env     = "JIRA_ISSUES_URL"
)

var (
	url string
	re  = regexp.MustCompile(pattern)
)

func getIssue(text string) string {
	if re.MatchString(text) {
		return re.FindString(text)
	}
	return ""
}

func getIssueURL(text string) (string, error) {
	issue := getIssue(text)
	if issue != "" {
		return url + strings.ToUpper(strings.TrimSpace(issue)), nil
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
