package cpf

import (
	"fmt"
	"strconv"

	"github.com/go-chat-bot/bot"
	cpfHelper "github.com/martinusso/go-docs/cpf"
)

const (
	msgParametroInvalido            = "Parâmetro inválido."
	msgQuantidadeParametrosInvalida = "Quantidade de parâmetros inválida."
	msgFmtCpfValido                 = "CPF %s é válido."
	msgFmtCpfInvalido               = "CPF %s é inválido."
)

func cpf(command *bot.Cmd) (string, error) {

	var param string
	switch len(command.Args) {
	case 0:
		param = "1"
	case 1:
		param = command.Args[0]
	default:
		return msgQuantidadeParametrosInvalida, nil

	}

	if len(param) > 2 {
		if valid(param) {
			return fmt.Sprintf(msgFmtCpfValido, command.Args[0]), nil
		}
		return fmt.Sprintf(msgFmtCpfInvalido, command.Args[0]), nil
	}

	qtCPF, err := strconv.Atoi(param)
	if err != nil {
		return msgParametroInvalido, nil
	}

	var cpf string
	for i := 0; i < qtCPF; i++ {
		cpf += gerarCPF() + " "
	}
	return cpf, nil
}

func gerarCPF() string {
	return cpfHelper.Generate()
}

func valid(cpf string) bool {
	return cpfHelper.Valid(cpf)
}

func init() {
	bot.RegisterCommand(
		"cpf",
		"Gerador/Validador de CPF.",
		"n para gerar n CPF e !cpf 12345678909 para validar um CPF",
		cpf)
}
