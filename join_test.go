package bot

import (
	"testing"
)

func TestJoin(t *testing.T) {
	channel1 := "#channel1"
	channel2 := "#channel2"
	command := Cmd{
		Command: "join",
		Args:    []string{channel1, channel2},
	}

	joinedChannels := []string{}
	joinFunc := func(channel string) {
		joinedChannels = append(joinedChannels, channel)
	}

	conn := &ConnectionMock{
		JoinFunc: joinFunc,
	}

	Join(command, conn)

	if len(joinedChannels) != 2 {
		t.Fail()
	}

	if joinedChannels[0] != channel1 {
		t.Errorf("Channel1 not joined. Expected %v got %v", channel1, joinedChannels[0])
	}

	if joinedChannels[1] != channel2 {
		t.Errorf("Channel2 not joined. Expected %v got %v", channel2, joinedChannels[1])
	}
}
