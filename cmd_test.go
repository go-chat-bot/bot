package bot

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

var (
	channel string
	replies []string
	user    *User
)

const (
	expectedMsg    = "msg"
	cmd            = "cmd"
	cmdDescription = "Command description"
	cmdExampleArgs = "arg1 arg2"
)

func responseHandler(target string, message string, sender *User) {
	channel = target
	user = sender
	replies = append(replies, message)
}

func resetResponses() {
	channel = ""
	user = &User{Nick: ""}
	replies = []string{}
	commands = make(map[string]*customCommand)
}

func newBot() *Bot {
	return New(&Handlers{
		Response: responseHandler,
	})
}

func resetRegisteredPeriodicCommands() {
	periodicCommands = make(map[string]PeriodicConfig)
}

func registerValidCommand() {
	RegisterCommand(cmd, cmdDescription, cmdExampleArgs,
		func(c *Cmd) (string, error) {
			return expectedMsg, nil
		})
}

func TestPeriodicCommands(t *testing.T) {
	resetResponses()
	resetRegisteredPeriodicCommands()
	RegisterPeriodicCommand("morning",
		PeriodicConfig{
			CronSpec: "0 0 08 * * mon-fri",
			Channels: []string{"#channel"},
			CmdFunc:  func(channel string) (string, error) { return "ok " + channel, nil },
		})
	b := New(&Handlers{Response: responseHandler})

	entries := b.cron.Entries()
	if len(entries) != 1 {
		t.Fatal("Should have one cron job entry")
	}
	if entries[0].Next.Hour() != 8 {
		t.Fatal("Cron job should be scheduled to 8am")
	}

	entries[0].Job.Run()

	if len(replies) != 1 {
		t.Fatal("Should have one reply in the channel")
	}
	if replies[0] != "ok #channel" {
		t.Fatal("Invalid reply")
	}
}

func TestMultiplePeriodicCommands(t *testing.T) {
	resetResponses()
	resetRegisteredPeriodicCommands()
	RegisterPeriodicCommand("morning",
		PeriodicConfig{
			CronSpec: "0 0 08 * * mon-fri",
			Channels: []string{"#channel"},
			CmdFunc:  func(channel string) (string, error) { return "ok_morning " + channel, nil },
		})
	RegisterPeriodicCommand("afternoon",
		PeriodicConfig{
			CronSpec: "0 0 12 * * mon-fri",
			Channels: []string{"#channel"},
			CmdFunc:  func(channel string) (string, error) { return "ok_afternoon " + channel, nil },
		})
	b := New(&Handlers{Response: responseHandler})

	entries := b.cron.Entries()
	if len(entries) != 2 {
		t.Fatal("Should have 2 cron job entries")
	}
	if entries[0].Next.Hour() != 8 {
		t.Fatal("First cron job should be scheduled for 8am")
	}
	if entries[1].Next.Hour() != 12 {
		t.Fatal("Second cron job should be schedule for 12am")
	}

	entries[0].Job.Run()
	entries[1].Job.Run()

	if len(replies) != 2 {
		t.Fatal("Should have two replies in the channel")
	}
	if replies[0] != "ok_morning #channel" {
		t.Fatal("Invalid reply in first cron job")
	}
	if replies[1] != "ok_afternoon #channel" {
		t.Fatal("Invalid reply in second cron job")
	}
}

func TestErroredPeriodicCommand(t *testing.T) {
	resetResponses()
	resetRegisteredPeriodicCommands()
	RegisterPeriodicCommand("bugged",
		PeriodicConfig{
			CronSpec: "0 0 08 * * mon-fri",
			Channels: []string{"#channel"},
			CmdFunc:  func(channel string) (string, error) { return "bug", errors.New("error") },
		})
	b := New(&Handlers{Response: responseHandler})

	entries := b.cron.Entries()

	if len(entries) != 1 {
		t.Fatal("Should have one cron job entry")
	}

	entries[0].Job.Run()

	if len(replies) != 0 {
		t.Fatal("Should not have a reply in the channel")
	}
}

func TestDisabledCommands(t *testing.T) {
	resetResponses()
	commands = make(map[string]*customCommand)
	b := New(&Handlers{
		Response: responseHandler,
	})

	RegisterCommand("cmd", "", "",
		func(c *Cmd) (string, error) {
			return "active", nil
		})

	RegisterPassiveCommand("passive",
		func(cmd *PassiveCmd) (string, error) {
			return "passive", nil
		})

	b.Disable([]string{"cmd"})
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!cmd"}, &User{Nick: "user"})
	if len(replies) != 0 {
		t.Fatal("Should not execute disabled active commands")
	}

	b.Disable([]string{"passive"})
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "regular message"}, &User{Nick: "user"})

	if len(replies) != 0 {
		t.Fatal("Should not execute disabled passive commands")
	}
}

func TestCommandNotRegistered(t *testing.T) {
	resetResponses()

	newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!not_a_cmd"}, &User{})

	if len(replies) != 0 {
		t.Fatal("Should not reply if a command is not found")
	}
}

func TestInvalidCmdArgs(t *testing.T) {
	resetResponses()
	registerValidCommand()

	newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!cmd \"invalid arg"}, &User{Nick: "user"})

	if channel != "#go-bot" {
		t.Error("Should reply to #go-bot channel")
	}
	if len(replies) != 1 {
		t.Fatal("Invalid reply")
	}
	if !strings.HasPrefix(replies[0], "Error parsing") {
		t.Fatal("Should reply with an error message")
	}
}

func TestErroredCmd(t *testing.T) {
	resetResponses()
	cmdError := errors.New("error")
	RegisterCommand("cmd", "", "",
		func(c *Cmd) (string, error) {
			return "", cmdError
		})

	newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!cmd"}, &User{Nick: "user"})

	if channel != "#go-bot" {
		t.Fatal("Invalid channel")
	}
	if len(replies) != 1 {
		t.Fatal("Invalid reply")
	}
	if replies[0] != fmt.Sprintf(errorExecutingCommand, "cmd", cmdError.Error()) {
		t.Fatal("Reply should contain the error message")
	}
}

func TestValidCmdOnChannel(t *testing.T) {
	resetResponses()
	registerValidCommand()

	newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!cmd"}, &User{Nick: "user"})

	if channel != "#go-bot" {
		t.Fatal("Command called on channel should reply to channel")
	}
	if len(replies) != 1 {
		t.Fatal("Should have one reply on channel")
	}
	if replies[0] != expectedMsg {
		t.Fatal("Invalid command reply")
	}
}

func TestChannelData(t *testing.T) {
	cd := ChannelData{
		Protocol: "irc",
		Server:   "myserver",
		Channel:  "#mychan",
	}
	if cd.URI() != "irc://myserver/#mychan" {
		t.Fatal("URI should return a valid IRC URI")
	}
}

func TestHelpWithNoArgs(t *testing.T) {
	resetResponses()
	registerValidCommand()
	newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!help"}, &User{Nick: "user"})

	expectedReply := []string{
		fmt.Sprintf(helpAboutCommand, CmdPrefix),
		fmt.Sprintf(availableCommands, "cmd"),
	}

	if !reflect.DeepEqual(replies, expectedReply) {
		t.Fatalf("Invalid reply. Expected %v got %v", expectedReply, replies)
	}
}

func TestDisableHelp(t *testing.T) {
	resetResponses()
	registerValidCommand()
	b := newBot()
	b.Disable([]string{"help"})
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!help"}, &User{Nick: "user"})

	if len(replies) > 0 {
		t.Fatalf("Should not execute help after disabling it")
	}
}

func TestHelpForACommand(t *testing.T) {
	resetResponses()
	registerValidCommand()
	newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!help cmd"}, &User{Nick: "user"})

	expectedReply := []string{
		fmt.Sprintf(helpDescripton, cmdDescription),
		fmt.Sprintf(helpUsage, CmdPrefix, cmd, cmdExampleArgs),
	}

	if !reflect.DeepEqual(replies, expectedReply) {
		t.Fatalf("Invalid reply. Expected %v got %v", expectedReply, replies)
	}
}

func TestHelpWithNonExistingCommand(t *testing.T) {
	resetResponses()
	registerValidCommand()
	newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!help not_a_cmd"}, &User{Nick: "user"})

	expectedReply := []string{
		fmt.Sprintf(helpAboutCommand, CmdPrefix),
		fmt.Sprintf(availableCommands, "cmd"),
	}

	if !reflect.DeepEqual(replies, expectedReply) {
		t.Fatalf("Invalid reply. Expected %v got %v", expectedReply, replies)
	}
}

func TestHelpWithInvalidArgs(t *testing.T) {
	resetResponses()
	registerValidCommand()
	newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!help cmd \"invalid arg"}, &User{Nick: "user"})

	if len(replies) != 1 {
		t.Fatal("Invalid reply")
	}
	if !strings.HasPrefix(replies[0], "Error parsing") {
		t.Fatal("Should reply with an error message")
	}
}

func TestCmdV2(t *testing.T) {
	resetResponses()
	RegisterCommandV2("cmd", "", "",
		func(c *Cmd) (CmdResult, error) {
			return CmdResult{
				Channel: "#channel",
				Message: "message"}, nil
		})

	newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!cmd"}, &User{Nick: "user"})

	if channel != "#channel" {
		t.Error("Wrong channel")
	}
	if !reflect.DeepEqual([]string{"message"}, replies) {
		t.Error("Invalid reply")
	}
}

func TestCmdV2WithoutSpecifyingChannel(t *testing.T) {
	resetResponses()
	RegisterCommandV2("cmd", "", "",
		func(c *Cmd) (CmdResult, error) {
			return CmdResult{Message: "message"}, nil
		})

	newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!cmd"}, &User{Nick: "user"})

	if channel != "#go-bot" {
		t.Error("Should reply to original channel if no channel is returned")
	}
}

func TestPassiveCommand(t *testing.T) {
	resetResponses()

	passiveCommands = make(map[string]passiveCmdFunc)

	echo := func(cmd *PassiveCmd) (string, error) {
		return cmd.Raw, nil
	}
	ping := func(cmd *PassiveCmd) (string, error) {
		return "pong", nil
	}
	errored := func(cmd *PassiveCmd) (string, error) {
		return "", errors.New("error")
	}

	RegisterPassiveCommand("echo", echo)
	RegisterPassiveCommand("ping", ping)
	RegisterPassiveCommand("errored", errored)

	newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "test"}, &User{Nick: "user"})

	if channel != "#go-bot" {
		t.Error("Invalid channel")
	}
	if len(replies) != 2 {
		t.Fatal("Invalid reply")
	}

	sort.Strings(replies)
	if replies[0] != "pong" {
		t.Error("ping command not executed")
	}
	if replies[1] != "test" {
		t.Error("echo command not executed")
	}
}

func TestCmdV3(t *testing.T) {
	resetResponses()
	result := CmdResultV3{
		Channel: "#channel",
		Message: make(chan string),
		Done:    make(chan bool)}
	RegisterCommandV3("cmd", "", "",
		func(c *Cmd) (CmdResultV3, error) {
			return result, nil
		})

	go newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!cmd"}, &User{Nick: "user"})
	result.Message <- "message"
	result.Done <- true

	if channel != "#channel" {
		t.Error("Wrong channel")
	}
	if !reflect.DeepEqual([]string{"message"}, replies) {
		t.Error("Invalid reply")
	}
}

func TestCmdV3WithoutSpecifyingChannel(t *testing.T) {
	resetResponses()
	result := CmdResultV3{
		Message: make(chan string),
		Done:    make(chan bool)}
	RegisterCommandV3("cmd", "", "",
		func(c *Cmd) (CmdResultV3, error) {
			return result, nil
		})

	go newBot().MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!cmd"}, &User{Nick: "user"})
	result.Message <- "message"
	result.Done <- true

	if channel != "#go-bot" {
		t.Error("Should reply to original channel if no channel is returned")
	}
	if !reflect.DeepEqual([]string{"message"}, replies) {
		t.Error("Invalid reply")
	}
}
