package commands

import (
	"fmt"
	"testing"
)

func TestRegisterCommand(t *testing.T) {
	command := "teste"
	ok := false
	fn := func(args []string) string { ok = true; return "" }
	RegisterCommand(command, fn)

	cmd := commands[command]
	cmd([]string{""})

	if !ok {
		t.Fail()
	}
}

func TestHandleExistingCommand(t *testing.T) {

}

func TestHandleCommandNotFound(t *testing.T) {
	commands = make(map[string]CommandFunc)

	channel := ""
	msg := []string{}
	fn := func(c string, m string) {
		channel = c
		msg = append(msg, m)
	}

	expectedChannel := "#go-bot"

	cmd := &Command{}
	cmd.Command = "allyourbase"
	HandleCmd(cmd, expectedChannel, fn)

	if channel != expectedChannel {
		t.Errorf("Invalid channel. Expected '%v' got '%v'", expectedChannel, channel)
	}

	if len(msg) != 2 {
		t.Errorf("Invalid msg length. Expected 2 got %v", len(msg))
	}

	expectedMsg := fmt.Sprintf(commandNotAvailable, cmd.Command)
	if msg[0] != expectedMsg {
		t.Errorf("Invalid msg. Expected '%v' got '%v'", expectedMsg, msg)
	}
}
