package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/go-chat-bot/bot"
)

const (
	invalidAmountOfParams = "Invalid amount of parameters"
	invalidParams         = "Invalid parameters"
)

func crypto(command *bot.Cmd) (string, error) {

	if len(command.Args) < 2 {
		return invalidAmountOfParams, nil
	}

	inputData := []byte(strings.Join(command.Args[1:], " "))
	switch strings.ToUpper(command.Args[0]) {
	case "MD5":
		return encryptMD5(inputData), nil
	case "SHA1", "SHA-1":
		return encryptSHA1(inputData), nil
	default:
		return invalidParams, nil
	}
}

func encryptMD5(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

func encryptSHA1(data []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(data))
}

func init() {
	bot.RegisterCommand(
		"crypto",
		"Encrypts the input data from its hash value",
		"md5|sha-1 enter here text to encrypt",
		crypto)
}
