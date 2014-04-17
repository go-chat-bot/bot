package main

import (
	"github.com/fabioxgn/go-bot/cmd"
	"github.com/fabioxgn/go-bot/irc"
	"testing"
)

func TestPartNoArgs(t *testing.T) {
	channel := "#go-bot"
	command := &cmd.Cmd{
		Command: "part",
		Args:    []string{},
		Channel: channel,
	}

	partedChannel := ""
	partFunc := func(channel string) {
		partedChannel = channel
	}

	conn := &irc.ConnectionMock{
		PartFunc: partFunc,
	}

	Part(command, conn)

	if partedChannel != channel {
		t.Errorf("Expected %v got %v", channel, partedChannel)
	}

}

func TestPartWithArgs(t *testing.T) {
	channel1 := "#go-bot"
	channel2 := "#lightirc"
	command := &cmd.Cmd{
		Command: "part",
		Args:    []string{channel1, channel2},
	}

	partedChannels := []string{}
	partFunc := func(channel string) {
		partedChannels = append(partedChannels, channel)
	}

	conn := &irc.ConnectionMock{
		PartFunc: partFunc,
	}

	Part(command, conn)

	want := 2
	got := len(partedChannels)
	if got != want {
		t.Errorf("Expected to part %v channels got %v", want, got)
	}

	if partedChannels[0] != channel1 {
		t.Errorf("Channel1 not parted. Expected %v got %v", channel1, partedChannels[0])
	}

	if partedChannels[1] != channel2 {
		t.Errorf("Channel2 not parted. Expected %v got %v", channel2, partedChannels[1])
	}
}
