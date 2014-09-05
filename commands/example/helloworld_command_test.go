package example

import (
	"github.com/fabioxgn/go-bot"
	"testing"
)

func TestHelloworld(t *testing.T) {
	want := "Hello nick"
	bot := &bot.Cmd{
		Command: "helloworld",
		Nick:    "nick",
	}
	got, error := hello(bot)

	if got != want {
		t.Errorf("Expected '%v' got '%v'", want, got)
	}

	if error != nil {
		t.Errorf("Expected '%v' got '%v'", nil, error)
	}
}
