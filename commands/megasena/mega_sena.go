package megasena

import (
	"bytes"
	"fmt"
	"github.com/fabioxgn/go-bot/cmd"
	"math/rand"
	"sort"
	"strings"
	"time"
)

const (
	DigitosJogo = 6 //TODO: Suportar até 15 números
)

func megasena(command *cmd.Cmd) (msg string, err error) {
	msg = fmt.Sprintf("%s, %s", command.Nick, sortear(60))
	return
}

func sortear(limit int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	numeros := make([]int, DigitosJogo)
	for i := 0; i < DigitosJogo; i++ {
		for {
			numero := r.Intn(limit + 1)
			if duplicado(numero, numeros) {
				continue
			}
			numeros[i] = numero
			break
		}
	}

	return formatarJogo(numeros)
}

func formatarJogo(numeros []int) string {
	sort.Ints(numeros)

	var jogo bytes.Buffer
	for _, numero := range numeros {
		jogo.WriteString(fmt.Sprintf(" %0.2d", numero))
	}

	return strings.TrimSpace(jogo.String())
}

func duplicado(numero int, numeros []int) bool {
	for _, i := range numeros {
		if i == numero {
			return true
		}
	}
	return false
}

func init() {
	cmd.RegisterCommand(&cmd.CustomCommand{
		Cmd:         "megasena",
		CmdFunc:     megasena,
		Description: "Gera os 6 números que te deixarão rico ;) ou não.",
	})
}
