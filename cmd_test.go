package bot

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	channel     string
	replies     chan string
	cmdError    chan string
	user        *User
	msgs        []string
	errs        []string
	protoParams interface{}
)

const (
	expectedMsg    = "msg"
	cmd            = "cmd"
	cmdDescription = "Command description"
	cmdExampleArgs = "arg1 arg2"
)

func waitMessages(t *testing.T, count int, errorCount int) {
	for {
		select {
		case reply := <-replies:
			msgs = append(msgs, reply)
		case err := <-cmdError:
			errs = append(errs, err)
		case <-time.After(1 * time.Second):
			t.Error("Timeout waiting for messages")
			t.Errorf("msgs received: %v", msgs)
			t.Errorf("errs received: %v", errs)
			t.Fatal()
		}
		if len(msgs) == count && len(errs) == errorCount {
			return
		}
	}
}

func responseHandler(target string, message string, sender *User) {
	channel = target
	user = sender
	replies <- message
}

func responseHandlerV2(om OutgoingMessage) {
	channel = om.Target
	user = om.Sender
	protoParams = om.ProtoParams
	replies <- om.Message
}

func errorHandler(msg string, err error) {
	cmdError <- fmt.Sprintf("%s: %s", msg, err)
}

func reset() {
	channel = ""
	user = &User{Nick: ""}
	replies = make(chan string, 10)
	cmdError = make(chan string, 10)
	msgs = []string{}
	errs = []string{}
	protoParams = nil
	commands = make(map[string]*customCommand)
	periodicCommands = make(map[string]PeriodicConfig)
	passiveCommands = make(map[string]*customCommand)
	filterCommands = make(map[string]*customCommand)
}

func newBot() *Bot {
	return New(&Handlers{
		Response: responseHandler,
		Errored:  errorHandler,
	},
		&Config{
			Protocol: "test",
			Server:   "test",
		},
	)
}

func newBotV2() *Bot {
	return New(&Handlers{
		Response:   responseHandler,
		ResponseV2: responseHandlerV2,
		Errored:    errorHandler,
	},
		&Config{
			Protocol: "test",
			Server:   "test",
		},
	)
}

func registerValidCommand() {
	RegisterCommand(cmd, cmdDescription, cmdExampleArgs,
		func(c *Cmd) (string, error) {
			return expectedMsg, nil
		})
}

func TestPeriodicCommands(t *testing.T) {
	reset()
	RegisterPeriodicCommand("morning",
		PeriodicConfig{
			CronSpec: "0 08 * * mon-fri",
			Channels: []string{"#channel"},
			CmdFunc:  func(channel string) (string, error) { return "ok " + channel, nil },
		})
	b := New(
		&Handlers{Response: responseHandler},
		&Config{Protocol: "test", Server: "test"},
	)
	defer b.Close()

	entries := b.cron.Entries()
	if len(entries) != 1 {
		t.Fatal("Should have one cron job entry")
	}
	if entries[0].Next.Hour() != 8 {
		t.Fatal("Cron job should be scheduled to 8am")
	}

	entries[0].Job.Run()

	waitMessages(t, 1, 0)

	if msgs[0] != "ok #channel" {
		t.Fatal("Invalid reply")
	}
}
func TestMultiplePeriodicCommands(t *testing.T) {
	reset()
	RegisterPeriodicCommand("morning",
		PeriodicConfig{
			CronSpec: "0 08 * * mon-fri",
			Channels: []string{"#channel"},
			CmdFunc:  func(channel string) (string, error) { return "ok_morning " + channel, nil },
		})
	RegisterPeriodicCommand("afternoon",
		PeriodicConfig{
			CronSpec: "0 12 * * mon-fri",
			Channels: []string{"#channel"},
			CmdFunc:  func(channel string) (string, error) { return "ok_afternoon " + channel, nil },
		})
	b := New(
		&Handlers{Response: responseHandler},
		&Config{Protocol: "test", Server: "test"},
	)
	defer b.Close()

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

	waitMessages(t, 2, 0)

	if len(msgs) != 2 {
		t.Fatal("Should have two replies in the channel")
	}
	sort.Strings(msgs)
	if msgs[0] != "ok_afternoon #channel" {
		t.Fatal("Invalid reply in afternoon cron job")
	}
	if msgs[1] != "ok_morning #channel" {
		t.Fatalf("Invalid reply in morning cron job.")
	}
}

func TestErroredPeriodicCommand(t *testing.T) {
	reset()
	RegisterPeriodicCommand("bugged",
		PeriodicConfig{
			CronSpec: "0 08 * * mon-fri",
			Channels: []string{"#channel"},
			CmdFunc:  func(channel string) (string, error) { return "bug", errors.New("error") },
		})
	b := newBot()
	defer b.Close()

	entries := b.cron.Entries()

	if len(entries) != 1 {
		t.Fatal("Should have one cron job entry")
	}

	entries[0].Job.Run()
	waitMessages(t, 0, 1)

	if len(msgs) != 0 {
		t.Error("Should not have a reply in the channel")
	}
	if len(errs) != 1 {
		t.Error("Expected 1 error")
	}
}

func TestPeriodicCommandsV2(t *testing.T) {
	reset()
	RegisterPeriodicCommandV2("morning",
		PeriodicConfig{
			CronSpec: "0 08 * * mon-fri",
			CmdFuncV2: func() ([]CmdResult, error) {
				ret := []CmdResult{
					{Message: "message 1", Channel: "#channel1"},
					{Message: "message 2", Channel: "#channel2"}}
				return ret, nil
			},
		})
	channels := make([]string, 0, 2)
	b := New(&Handlers{Response: func(target string, message string, sender *User) {
		channels = append(channels, target)
		channel = target
		user = sender
		replies <- message
	}},
		&Config{Protocol: "test", Server: "test"})

	defer b.Close()

	entries := b.cron.Entries()
	if len(entries) != 1 {
		t.Fatal("Should have one cron job entry")
	}
	if entries[0].Next.Hour() != 8 {
		t.Fatal("Cron job should be scheduled to 8am")
	}

	entries[0].Job.Run()

	waitMessages(t, 2, 0)
	if len(channels) != 2 {
		t.Fatal("Should have 2 destinations channels", len(channels))
	}

	if msgs[0] != "message 1" {
		t.Fatal("Invalid first reply")
	}

	if channels[0] != "#channel1" {
		t.Fatal("Invalid channel for first message", channels[0])
	}

	if msgs[1] != "message 2" {
		t.Fatal("Invalid second reply")
	}

	if channels[1] != "#channel2" {
		t.Fatal("Invalid channel for second message", channels[1])
	}
}

func TestErroredPeriodicCommandsV2(t *testing.T) {
	reset()
	RegisterPeriodicCommandV2("morning",
		PeriodicConfig{
			CronSpec: "0 08 * * mon-fri",
			CmdFuncV2: func() ([]CmdResult, error) {
				return nil, errors.New("error")
			},
		})
	b := newBot()
	defer b.Close()

	entries := b.cron.Entries()
	if len(entries) != 1 {
		t.Fatal("Should have one cron job entry")
	}
	if entries[0].Next.Hour() != 8 {
		t.Fatal("Cron job should be scheduled to 8am")
	}

	entries[0].Job.Run()

	waitMessages(t, 0, 1)
	if len(msgs) != 0 {
		t.Error("Should not have a reply in the channel")
	}
	if len(errs) != 1 {
		t.Error("Expected 1 error")
	}
}

func TestDisabledCommands(t *testing.T) {
	reset()
	commands = make(map[string]*customCommand)
	b := newBot()

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

	time.Sleep(100)
	if len(msgs) != 0 {
		t.Fatal("Should not execute disabled active commands")
	}

	b.Disable([]string{"passive"})
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "regular message"}, &User{Nick: "user"})

	time.Sleep(100)
	if len(msgs) != 0 {
		t.Fatal("Should not execute disabled passive commands")
	}
}

func TestCommandNotRegistered(t *testing.T) {
	reset()
	b := newBot()
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!not_a_cmd"}, &User{})

	time.Sleep(100)

	if len(msgs) != 0 {
		t.Fatal("Should not reply if a command is not found")
	}
}

func TestInvalidCmdArgs(t *testing.T) {
	reset()
	registerValidCommand()

	b := newBot()
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!cmd \"invalid arg"}, &User{Nick: "user"})

	waitMessages(t, 1, 0)

	if channel != "#go-bot" {
		t.Error("Should reply to #go-bot channel")
	}
	if !strings.HasPrefix(msgs[0], "Error parsing") {
		t.Fatal("Should reply with an error message")
	}
}

func TestErroredCmd(t *testing.T) {
	reset()
	cmdError := errors.New("error")
	RegisterCommand("cmd", "", "",
		func(c *Cmd) (string, error) {
			return "", cmdError
		})

	b := newBot()
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!cmd"}, &User{Nick: "user"})

	waitMessages(t, 1, 1)

	if channel != "#go-bot" {
		t.Fatal("Invalid channel")
	}
	if msgs[0] != fmt.Sprintf(errorExecutingCommand, "cmd", cmdError.Error()) {
		t.Fatal("Reply should contain the error message")
	}
	if len(errs) != 1 {
		t.Error("Expected the command to return an error")
	}
}

func TestValidCmdOnChannel(t *testing.T) {
	reset()
	registerValidCommand()

	b := newBot()
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!cmd"}, &User{Nick: "user"})

	waitMessages(t, 1, 0)

	if channel != "#go-bot" {
		t.Fatal("Command called on channel should reply to channel")
	}
	if msgs[0] != expectedMsg {
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
	reset()
	registerValidCommand()
	b := newBot()
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!help"}, &User{Nick: "user"})

	waitMessages(t, 2, 0)

	expectedReply := []string{
		fmt.Sprintf(helpAboutCommand, CmdPrefix),
		fmt.Sprintf(availableCommands, "cmd"),
	}

	if !reflect.DeepEqual(msgs, expectedReply) {
		t.Fatalf("Invalid reply. Expected %v got %v", expectedReply, msgs)
	}
}

func TestDisableHelp(t *testing.T) {
	reset()
	registerValidCommand()
	b := newBot()
	b.Disable([]string{"help"})
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!help"}, &User{Nick: "user"})

	time.Sleep(100)

	if len(replies) > 0 {
		t.Fatalf("Should not execute help after disabling it")
	}
}

func TestHelpForACommand(t *testing.T) {
	reset()
	registerValidCommand()
	b := newBot()
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!help cmd"}, &User{Nick: "user"})

	waitMessages(t, 2, 0)

	expectedReply := []string{
		fmt.Sprintf(helpDescripton, cmdDescription),
		fmt.Sprintf(helpUsage, CmdPrefix, cmd, cmdExampleArgs),
	}

	if !reflect.DeepEqual(msgs, expectedReply) {
		t.Fatalf("Invalid reply. Expected %v got %v", expectedReply, msgs)
	}
}

func TestHelpWithNonExistingCommand(t *testing.T) {
	reset()
	registerValidCommand()
	b := newBot()
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!help not_a_cmd"}, &User{Nick: "user"})

	expectedReply := []string{
		fmt.Sprintf(helpAboutCommand, CmdPrefix),
		fmt.Sprintf(availableCommands, "cmd"),
	}

	waitMessages(t, 2, 0)

	if !reflect.DeepEqual(msgs, expectedReply) {
		t.Fatalf("Invalid reply. Expected %v got %v", expectedReply, msgs)
	}
}

func TestHelpWithInvalidArgs(t *testing.T) {
	reset()
	registerValidCommand()
	b := newBot()
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!help cmd \"invalid arg"}, &User{Nick: "user"})

	waitMessages(t, 1, 0)

	if !strings.HasPrefix(msgs[0], "Error parsing") {
		t.Fatal("Should reply with an error message")
	}
}

func TestCmdV2(t *testing.T) {
	reset()
	RegisterCommandV2("cmd", "", "",
		func(c *Cmd) (CmdResult, error) {
			return CmdResult{
				Channel: "#channel",
				Message: "message"}, nil
		})

	b := newBot()
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!cmd"}, &User{Nick: "user"})

	waitMessages(t, 1, 0)

	if channel != "#channel" {
		t.Error("Wrong channel")
	}
	if !reflect.DeepEqual([]string{"message"}, msgs) {
		t.Error("Invalid reply")
	}
}

func TestCmdV2WithProtoParams(t *testing.T) {
	reset()
	RegisterCommandV2("cmd", "", "",
		func(c *Cmd) (CmdResult, error) {
			return CmdResult{
				Channel:     "#channel",
				Message:     "message",
				ProtoParams: &CmdResult{Message: "Nested!"},
			}, nil
		})

	b := newBotV2()
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!cmd"}, &User{Nick: "user"})

	waitMessages(t, 1, 0)

	if channel != "#channel" {
		t.Error("Wrong channel")
	}
	if !reflect.DeepEqual([]string{"message"}, msgs) {
		t.Error("Invalid reply")
	}
	if pa, ok := protoParams.(*CmdResult); ok {
		if pa.Message != "Nested!" {
			t.Error("Information lost in copying.")
		}
	} else {
		t.Error("Failed to pass proto args through.")
	}
}

func TestCmdV2WithoutSpecifyingChannel(t *testing.T) {
	reset()
	RegisterCommandV2("cmd", "", "",
		func(c *Cmd) (CmdResult, error) {
			return CmdResult{Message: "message"}, nil
		})

	b := newBot()
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!cmd"}, &User{Nick: "user"})

	waitMessages(t, 1, 0)

	if channel != "#go-bot" {
		t.Error("Should reply to original channel if no channel is returned")
	}
}

func TestPassiveCommand(t *testing.T) {
	reset()
	passiveCommands = make(map[string]*customCommand)
	echo := func(cmd *PassiveCmd) (string, error) { return cmd.Raw, nil }
	ping := func(cmd *PassiveCmd) (string, error) { return "pong", nil }
	errored := func(cmd *PassiveCmd) (string, error) { return "", errors.New("error") }

	RegisterPassiveCommand("echo", echo)
	RegisterPassiveCommand("ping", ping)
	RegisterPassiveCommand("errored", errored)

	b := newBot()
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "test"}, &User{Nick: "user"})

	waitMessages(t, 2, 1)

	if channel != "#go-bot" {
		t.Error("Invalid channel")
	}
	if len(msgs) != 2 {
		t.Fatal("Invalid reply")
	}
	if len(errs) != 1 {
		t.Error("Expected 1 error")
	}

	sort.Strings(msgs)
	if msgs[0] != "pong" {
		t.Error("ping command not executed")
	}
	if msgs[1] != "test" {
		t.Error("echo command not executed")
	}
}

func TestPassiveCommandV2(t *testing.T) {
	reset()
	result := CmdResultV3{
		Channel: "#channel",
		Message: make(chan string),
		Done:    make(chan bool)}

	ping := func(cmd *PassiveCmd) (CmdResultV3, error) { return result, nil }
	errored := func(cmd *PassiveCmd) (CmdResultV3, error) { return CmdResultV3{}, errors.New("error") }

	RegisterPassiveCommandV2("ping", ping)
	RegisterPassiveCommandV2("errored", errored)

	b := newBot()
	go b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "test"}, &User{Nick: "user"})
	result.Message <- "pong"
	result.Done <- true

	waitMessages(t, 1, 1)

	if channel != "#channel" {
		t.Error("Invalid channel")
	}
	if len(msgs) != 1 {
		t.Fatal("Invalid reply")
	}

	if msgs[0] != "pong" {
		t.Error("ping command not executed")
	}
}

func TestCmdV3(t *testing.T) {
	reset()
	result := CmdResultV3{
		Channel: "#channel",
		Message: make(chan string),
		Done:    make(chan bool)}
	RegisterCommandV3("cmd", "", "",
		func(c *Cmd) (CmdResultV3, error) {
			return result, nil
		})

	b := newBot()
	go b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!cmd"}, &User{Nick: "user"})
	result.Message <- "message"
	result.Done <- true

	waitMessages(t, 1, 0)

	if channel != "#channel" {
		t.Error("Wrong channel")
	}
	if !reflect.DeepEqual([]string{"message"}, msgs) {
		t.Error("Invalid reply")
	}
}

func TestCmdV3WithoutSpecifyingChannel(t *testing.T) {
	reset()
	result := CmdResultV3{
		Message: make(chan string),
		Done:    make(chan bool)}
	RegisterCommandV3("cmd", "", "",
		func(c *Cmd) (CmdResultV3, error) {
			return result, nil
		})

	b := newBot()
	go b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "!cmd"}, &User{Nick: "user"})
	result.Message <- "message"
	result.Done <- true

	waitMessages(t, 1, 0)

	if channel != "#go-bot" {
		t.Error("Should reply to original channel if no channel is returned")
	}
	if !reflect.DeepEqual([]string{"message"}, msgs) {
		t.Error("Invalid reply")
	}
}

func TestFilterCommand(t *testing.T) {
	reset()
	passiveCommands = make(map[string]*customCommand)
	ping := func(cmd *PassiveCmd) (string, error) { return "pong", nil }
	modified := func(cmd *FilterCmd) (string, error) { return "PONG!", nil }
	errored := func(cmd *FilterCmd) (string, error) { return "", errors.New("error") }

	RegisterPassiveCommand("ping", ping)
	RegisterFilterCommand("modified", modified)
	RegisterFilterCommand("errored", errored)

	b := newBot()
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "test"}, &User{Nick: "user"})

	waitMessages(t, 1, 1)

	if channel != "#go-bot" {
		t.Error("Invalid channel")
	}
	if len(msgs) != 1 {
		t.Fatal("Invalid reply")
	}
	if len(errs) != 1 {
		t.Error("Expected 1 error")
	}

	sort.Strings(msgs)
	if msgs[0] != "PONG!" {
		t.Error("filter command not working")
	}
}

func TestFilterCommandSilence(t *testing.T) {
	reset()
	passiveCommands = make(map[string]*customCommand)
	ping := func(cmd *PassiveCmd) (string, error) { return "pong", nil }
	silenced := func(cmd *FilterCmd) (string, error) { return "", nil }
	errored := func(cmd *FilterCmd) (string, error) { return "Ignored", errors.New("error") }

	RegisterPassiveCommand("ping", ping)
	RegisterFilterCommand("silenced", silenced)
	RegisterFilterCommand("errored", errored)

	b := newBot()
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "test"}, &User{Nick: "user"})

	waitMessages(t, 0, 1)

	if len(msgs) != 0 {
		t.Fatal("Expected no messages!")
	}
	if len(errs) != 1 {
		t.Error("Expected 1 error")
	}
}

func TestFilterCommandSilenceSendV2(t *testing.T) {
	reset()
	passiveCommands = make(map[string]*customCommand)
	ping := func(cmd *PassiveCmd) (string, error) { return "pong", nil }
	silenced := func(cmd *FilterCmd) (string, error) { return "", nil }
	errored := func(cmd *FilterCmd) (string, error) { return "Ignored", errors.New("error") }

	RegisterPassiveCommand("ping", ping)
	RegisterFilterCommand("silenced", silenced)
	RegisterFilterCommand("errored", errored)

	b := newBotV2()
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, &Message{Text: "test"}, &User{Nick: "user"})

	waitMessages(t, 0, 1)

	if len(msgs) != 0 {
		t.Fatal("Expected no messages!")
	}
	if len(errs) != 1 {
		t.Error("Expected 1 error")
	}
}

// how to test channels..
// https://www.hugopicado.com/2016/10/01/testing-over-golang-channels.html

func TestMessageStreams(t *testing.T) {
	var mutex = &sync.Mutex{}
	reset()

	var msSender1 *MessageStream
	var msSender2 *MessageStream

	RegisterMessageStream("streamOne", func(ms1 *MessageStream) error {
		mutex.Lock()
		msSender1 = ms1
		mutex.Unlock()
		return nil
	})
	RegisterMessageStream("streamTwo", func(ms2 *MessageStream) error {
		mutex.Lock()
		msSender2 = ms2
		mutex.Unlock()
		return nil
	})

	b1 := New(&Handlers{Response: responseHandler, Errored: errorHandler}, &Config{Protocol: "protoA", Server: "test"})
	b2 := New(&Handlers{Response: responseHandler, Errored: errorHandler}, &Config{Protocol: "protoB", Server: "test"})

	msmB1 := MessageStreamMessage{
		Message:     "hello botOne",
		ChannelData: &ChannelData{Server: b1.Server, Protocol: b1.Protocol, Channel: "#go-bot"},
	}
	msmB2 := MessageStreamMessage{
		Message:     "hello botTwo",
		ChannelData: &ChannelData{Server: b2.Server, Protocol: b2.Protocol, Channel: "#go-bot"},
	}

	// give New() a second to make() the chans and setup the objects
	time.Sleep(2 * time.Second)

	// when you send a message destined for b1 #go-bot, even if you send it to b2, it should arrive at b1
	mutex.Lock()
	msSender1.Data <- msmB1
	if "hello botOne" != <-replies {
		t.Fatal("message not Recieved at Channel")
	}

	msSender2.Data <- msmB1
	if "hello botOne" != <-replies {
		t.Fatal("message not Recieved at Channel")
	}

	// and vice-versa
	// when you send a message destined for b2 #go-bots, even if you send it to b1, it should arrive at b2
	msSender1.Data <- msmB2
	if "hello botTwo" != <-replies {
		t.Fatal("message not Recieved at Channel")
	}

	msSender2.Data <- msmB2
	if "hello botTwo" != <-replies {
		t.Fatal("message not Recieved at Channel")
	}
	mutex.Unlock()

}
