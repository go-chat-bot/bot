package cnpj

import (
	"fmt"
	"strconv"

	"github.com/go-chat-bot/bot"
	cnpjHelper "github.com/martinusso/go-docs/cnpj"
)

const (
	msgParametroInvalido            = "Parâmetro inválido."
	msgQuantidadeParametrosInvalida = "Quantidade de parâmetros inválida."
	msgFmtCnpjValido                = "CNPJ %s é válido."
	msgFmtCnpjInvalido              = "CNPJ %s é inválido."
)

func cnpj(command *bot.Cmd) (string, error) {

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
			return fmt.Sprintf(msgFmtCnpjValido, command.Args[0]), nil
		}
		return fmt.Sprintf(msgFmtCnpjInvalido, command.Args[0]), nil
	}

	qtCNPJ, err := strconv.Atoi(param)
	if err != nil {
		return msgParametroInvalido, nil
	}

	var cnpj string
	for i := 0; i < qtCNPJ; i++ {
		cnpj += gerarCNPJ() + " "
	}
	return cnpj, nil
}

func gerarCNPJ() string {
	return cnpjHelper.Generate()
}

func valid(cnpj string) bool {
	return cnpjHelper.Valid(cnpj)
}

func init() {
	bot.RegisterCommand(
		"cnpj",
		"Gerador/Validador de CNPJ.",
		"n para gerar n CNPJ e !cnpj 99999999000191 para validar um CNPJ",
		cnpj)
}
