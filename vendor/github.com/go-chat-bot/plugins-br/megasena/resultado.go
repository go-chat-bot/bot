package megasena

import (
	"fmt"

	"github.com/go-chat-bot/plugins/web"
)

var (
	url = "http://developers.agenciaideias.com.br/loterias/megasena/json"
)

type retorno struct {
	Concurso struct {
		Numero           string   `json:"numero"`
		Data             string   `json:"data"`
		Cidade           string   `json:"cidade"`
		Local            string   `json:"local"`
		ValorAcumulado   string   `json:"valor_acumulado"`
		NumerosSorteados []string `json:"numeros_sorteados"`
		Premiacao        struct {
			Sena struct {
				Ganhadores string `json:"ganhadores"`
				ValorPago  string `json:"valor_pago"`
			} `json:"sena"`
			Quina struct {
				Ganhadores string `json:"ganhadores"`
				ValorPago  string `json:"valor_pago"`
			} `json:"quina"`
			Quadra struct {
				Ganhadores string `json:"ganhadores"`
				ValorPago  string `json:"valor_pago"`
			} `json:"quadra"`
		} `json:"premiacao"`
		ArrecadacaoTotal string `json:"arrecadacao_total"`
	} `json:"concurso"`
	ProximoConcurso struct {
		Data          string `json:"data"`
		ValorEstimado string `json:"valor_estimado"`
	} `json:"proximo_concurso"`
	ConcursoFinalZero struct {
		Numero         string `json:"numero"`
		ValorAcumulado string `json:"valor_acumulado"`
	} `json:"concurso_final_zero"`
	MegaViradaValorAcumulado string `json:"mega_virada_valor_acumulado"`
}

func resultado() (string, error) {
	data := &retorno{}
	err := web.GetJSON(url, data)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Sorteio %s de %s: %s - %s premiado(s) R$ %s.",
		data.Concurso.Numero,
		data.Concurso.Data,
		data.Concurso.NumerosSorteados,
		data.Concurso.Premiacao.Sena.Ganhadores,
		data.Concurso.Premiacao.Sena.ValorPago), nil
}
