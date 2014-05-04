package megasena

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Retorno struct {
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

func Resultado() string {
	url := "http://developers.agenciaideias.com.br/loterias/megasena/json"
	res, err := http.Get(url)
	if err != nil {
		return err.Error()
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err.Error()
	}
	data := &Retorno{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("Sorteio %s de %s: %s - %s premiado(s) R$ %s.",
		data.Concurso.Numero,
		data.Concurso.Data,
		data.Concurso.NumerosSorteados,
		data.Concurso.Premiacao.Sena.Ganhadores,
		data.Concurso.Premiacao.Sena.ValorPago)
}
