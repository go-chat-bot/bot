package megasena

import (
	"bytes"
	"fmt"
	"github.com/fabioxgn/go-bot"
	"math/rand"
	"sort"
	"strings"
	"time"
)

const (
	digitosJogo      = 6
	msgOpcaoInvalida = "Informe uma opção: gerar ou resultado"
)

func megasena(command *bot.Cmd) (msg string, err error) {
	if len(command.Args) == 0 {
		msg = msgOpcaoInvalida
	} else {
		switch command.Args[0] {
		case "gerar":
			msg = sortear(60)
		case "resultado":
			msg, err = resultado()
		}
	}
	msg = fmt.Sprintf("%s: %s", command.Nick, msg)
	return
}

func sortear(limit int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	numeros := make([]int, digitosJogo)
	for i := 0; i < digitosJogo; i++ {
		for {
			numero := r.Intn(limit + 1)
			if duplicado(numero, numeros) {
				continue
			}
			numeros[i] = numero
			break
		}
	}

	sort.Ints(numeros)
	return formatarJogo(numeros)
}

func formatarJogo(numeros []int) string {
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
	bot.RegisterCommand(
		"megasena",
		"Gera um jogo da megasena ou mostra os últimos números sorteados.",
		"gerar|resultado",
		megasena)
}
