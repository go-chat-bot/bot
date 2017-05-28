package bot

import (
	"errors"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	channel string
	replies []string
	user    *User
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
}

func resetRegisteredPeriodicCommands() {
	periodicCommands = make(map[string]PeriodicConfig)
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

	// Give a second for the crons to be registered
	time.Sleep(time.Second)

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

	// Give a second for the crons to be registered
	time.Sleep(time.Second)

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

	// Give a second for the crons to be registered
	time.Sleep(time.Second)

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
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, "!cmd", &User{Nick: "user"})
	if len(replies) != 0 {
		t.Fatal("Should not execute disabled active commands")
	}

	b.Disable([]string{"passive"})
	b.MessageReceived(&ChannelData{Channel: "#go-bot"}, "regular message", &User{Nick: "user"})

	if len(replies) != 0 {
		t.Fatal("Should not execute disabled passive commands")
	}
}

func TestMessageReceived(t *testing.T) {
	Convey("Given a new message in the channel", t, func() {
		Reset(resetResponses)
		commands = make(map[string]*customCommand)
		b := New(&Handlers{
			Response: responseHandler,
		})

		Convey("When the command is not registered", func() {
			Convey("It should not post to the channel", func() {
				b.MessageReceived(&ChannelData{Channel: "#go-bot"}, "!not_a_cmd", &User{})

				So(replies, ShouldBeEmpty)
			})
		})

		Convey("When the command arguments are invalid", func() {
			Convey("It should reply with an error message", func() {
				b.MessageReceived(&ChannelData{Channel: "#go-bot"}, "!cmd \"invalid arg", &User{Nick: "user"})

				So(channel, ShouldEqual, "#go-bot")
				So(replies, ShouldHaveLength, 1)
				So(replies[0], ShouldStartWith, "Error parsing")
			})
		})

		Convey("The command can return an error", func() {
			Convey("it sould send the message with the error to the channel", func() {
				cmdError := errors.New("error")
				RegisterCommand("cmd", "", "",
					func(c *Cmd) (string, error) {
						return "", cmdError
					})

				b.MessageReceived(&ChannelData{Channel: "#go-bot"}, "!cmd", &User{Nick: "user"})

				So(channel, ShouldEqual, "#go-bot")
				So(replies, ShouldResemble,
					[]string{fmt.Sprintf(errorExecutingCommand, "cmd", cmdError.Error())})
			})
		})

		Convey("When the command is valid and registered", func() {
			commands = make(map[string]*customCommand)
			expectedMsg := "msg"
			cmd := "cmd"
			cmdDescription := "Command description"
			cmdExampleArgs := "arg1 arg2"

			RegisterCommand(cmd, cmdDescription, cmdExampleArgs,
				func(c *Cmd) (string, error) {
					return expectedMsg, nil
				})

			Convey("If it is called in the channel, reply on the channel", func() {
				b.MessageReceived(&ChannelData{Channel: "#go-bot"}, "!cmd", &User{Nick: "user"})

				So(channel, ShouldEqual, "#go-bot")
				So(replies, ShouldResemble, []string{expectedMsg})
			})

			Convey("If it is a private message, reply to the user", func() {
				user = &User{Nick: "go-bot"}
				b.MessageReceived(&ChannelData{Channel: "#go-bot"}, "!cmd", &User{Nick: "sender-nick"})
				So(user.Nick, ShouldEqual, "sender-nick")
			})

			Convey("When the command is help", func() {
				Convey("Display the available commands in the channel", func() {
					b.MessageReceived(&ChannelData{Channel: "#go-bot"}, "!help", &User{Nick: "user"})

					So(channel, ShouldEqual, "#go-bot")
					So(replies, ShouldResemble, []string{
						fmt.Sprintf(helpAboutCommand, CmdPrefix),
						fmt.Sprintf(availableCommands, "cmd"),
					})
				})

				Convey("If the command exists send a message to the channel", func() {
					b.MessageReceived(&ChannelData{Channel: "#go-bot"}, "!help cmd", &User{Nick: "user"})

					So(channel, ShouldEqual, "#go-bot")
					So(replies, ShouldResemble, []string{
						fmt.Sprintf(helpDescripton, cmdDescription),
						fmt.Sprintf(helpUsage, CmdPrefix, cmd, cmdExampleArgs),
					})
				})

				Convey("If the command does not exists, display the generic help", func() {
					b.MessageReceived(&ChannelData{Channel: "#go-bot"}, "!help not_a_command", &User{Nick: "user"})

					So(channel, ShouldEqual, "#go-bot")
					So(replies, ShouldResemble, []string{
						fmt.Sprintf(helpAboutCommand, CmdPrefix),
						fmt.Sprintf(availableCommands, "cmd"),
					})
				})

				Convey("if the help arguments are invalid, reply with an error", func() {
					b.MessageReceived(&ChannelData{Channel: "#go-bot"}, "!help cmd \"invalid arg", &User{Nick: "user"})

					So(channel, ShouldEqual, "#go-bot")
					So(replies, ShouldHaveLength, 1)
					So(replies[0], ShouldStartWith, "Error parsing")
				})
			})
		})

		Convey("When the command is V2", func() {
			Convey("it should send the message with the error to the channel", func() {
				RegisterCommandV2("cmd", "", "",
					func(c *Cmd) (CmdResult, error) {
						return CmdResult{
							Channel: "#channel",
							Message: "message"}, nil
					})

				b.MessageReceived(&ChannelData{Channel: "#go-bot"}, "!cmd", &User{Nick: "user"})

				So(channel, ShouldEqual, "#channel")
				So(replies, ShouldResemble, []string{"message"})
			})

			Convey("it should reply to the current channel if the command does not specify one", func() {
				RegisterCommandV2("cmd", "", "",
					func(c *Cmd) (CmdResult, error) {
						return CmdResult{
							Message: "message"}, nil
					})

				b.MessageReceived(&ChannelData{Channel: "#go-bot"}, "!cmd", &User{Nick: "user"})

				So(channel, ShouldEqual, "#go-bot")
				So(replies, ShouldResemble, []string{"message"})
			})
		})

		Convey("When the command is passive", func() {
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

			Convey("If it is called in the channel, reply on the channel", func() {
				b.MessageReceived(&ChannelData{Channel: "#go-bot"}, "test", &User{Nick: "user"})

				So(channel, ShouldEqual, "#go-bot")
				So(len(replies), ShouldEqual, 2)
				So(replies, ShouldContain, "test")
				So(replies, ShouldContain, "pong")
			})

			Convey("If it is a private message, reply to the user", func() {
				user = &User{Nick: "go-bot"}
				b.MessageReceived(&ChannelData{Channel: "#go-bot"}, "test", &User{Nick: "sender-nick"})

				So(user.Nick, ShouldEqual, "sender-nick")
				So(len(replies), ShouldEqual, 2)
				So(replies, ShouldContain, "test")
				So(replies, ShouldContain, "pong")
			})
		})
	})
}

func TestChannelData(t *testing.T) {
	Convey("Given a ChannelData struct", t, func() {
		Convey("Make sure ChannelData can give you the Channel URI", func() {
			cd := ChannelData{
				Protocol: "irc",
				Server:   "myserver",
				Channel:  "#mychan",
			}
			So(cd.URI(), ShouldEqual, "irc://myserver/#mychan")
		})
	})
}
