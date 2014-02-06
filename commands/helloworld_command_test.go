package commands

import (
	"testing"
)

func TestHelloworld(t *testing.T) {
	want := "Hello world!"
	got := "Hello world!"

	if got != want {
		t.Errorf("Expected '%v' got '%v'", want, got)
	}
}
