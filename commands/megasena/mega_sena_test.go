package megasena

import (
	"github.com/fabioxgn/go-bot/bot"
	"regexp"
	"testing"
)

func TestSortear(t *testing.T) {
	want := "01 02 03 04 05 06"
	got := sortear(6)

	if want != got {
		t.Errorf("Expected %v got %v", want, got)
	}
}

func TestMegaSena(t *testing.T) {
	cmd := &bot.Cmd{
		Command: "megasena",
		Nick:    "nick",
		Args:    []string{"gerar"},
	}
	got, err := megasena(cmd)

	if err != nil {
		t.Errorf("Expected '%v' got '%v'", nil, err)
	}

	match, err := regexp.MatchString("nick: (\\d{2} {1}){5}\\d{2}", got)
	if !match {
		t.Errorf("got %v", got)
	}
	if err != nil {
		t.Errorf(err.Error())
	}
}
