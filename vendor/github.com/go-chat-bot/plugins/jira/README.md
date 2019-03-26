### Overview

This is a plugin for [Atlassian JIRA](https://www.atlassian.com/software/jira)
issue tracking system.

### Features
* Simple authentication support for JIRA
* Outputs some of the issue details in the channels not just links
* Optional per-channel configuration of issue message format

### Setup
* Set up JIRA_BASE_URL env variable to your JIRA server URL. For example
  https://issues.jenkins-ci.org
* Set up JIRA_USER env variable to JIRA username for the bot account
* Set up JIRA_PASS env variable to JIRA password for the bot account

In addition to the above channel-specific configuration variables can be defined
in a separate JSON configuration file loaded from path specified by environment
variable `JIRA_CONFIG_FILE`. Example file can be seen in
`example_config.json`. It is an array of channel configurations with each
configuration having:
 * `channel` for which the configuration is intended
 * `template` to override default issue template (see Issue Formatting)
 * `notifyNew` is array of JIRA project keys to watch for new issues
 * `notifyResolved` is array JIRA project keys to watch for resolved issues

### Issue Formatting

By default the plugin will output issues in the following format:
```
<key> (<assignee>, <status>): <summary> - <url>
```
To see which values are available for use in templates see
[go-jira](https://github.com/andygrunwald/go-jira/blob/master/issue.go).

The format used is go template notation on the issue object. If you want to just
post URL to the issue itself you can configure it by setting the template to
`{{.Self}}` for given channel in the configuration file.

Default template looks like this:
```
{{.Key}} ({{.Fields.Assignee.Key}}, {{.Fields.Status.Name}}): {{.Fields.Summary}} - {{.Self}}
```

`JIRA_NOTIFY_INTERVAL` environment variable can be used to control how often the
notification methods will be run. It defaults to be run every minute.
