package main

import (
	"encoding/json"
	"io"
)

// Struct mapping the config
type Config struct {
	Server   string
	Channels []string
	User     string
	Nick     string
	Cmd      string
	UseTLS   bool
}

// Read the configuration from a JSON file
func (c *Config) Read(value io.Reader) {
	decoder := json.NewDecoder(value)
	err := decoder.Decode(&c)
	if err != nil {
		panic(err)
	}
}
