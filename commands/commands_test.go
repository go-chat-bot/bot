package commands

import (
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

func TestCommandNotFound(t *testing.T) {

}
