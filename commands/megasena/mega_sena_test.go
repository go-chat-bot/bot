package megasena

import (
	"fmt"
	"github.com/fabioxgn/go-bot"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

const (
	retornoJSON = `{"concurso":{
    "numero":"1636",
    "data":"17\/09\/2014",
    "cidade":"OSASCO-SP",
    "local":"Caminh\u00e3o da Sorte",
    "valor_acumulado":"29.530.043,53",
    "numeros_sorteados":[
      "19",
      "26",
      "33",
      "35",
      "51",
      "52"
    ],
    "premiacao":{
      "sena":{
        "ganhadores":"0",
        "valor_pago":"0,00"
      },
      "quina":{
        "ganhadores":"90",
        "valor_pago":"38.637,27"
      },
      "quadra":{
        "ganhadores":"8.474",
        "valor_pago":"586,22"
      }
    },
    "arrecadacao_total":"59.395.800,00"
  },
  "proximo_concurso":{
    "data":"20\/09\/2014",
    "valor_estimado":"37.000.000,00"
  },
  "concurso_final_zero":{
    "numero":"1640",
    "valor_acumulado":"7.271.924,76"
  },
  "mega_virada_valor_acumulado":"54.516.366,32"}`
)

func TestSortear(t *testing.T) {
	Convey("Sortear", t, func() {
		So(sortear(6), ShouldEqual, "01 02 03 04 05 06")
	})
}

func TestMegaSena(t *testing.T) {
	Convey("Megasena", t, func() {

		cmd := &bot.Cmd{
			Command: "megasena",
			Nick:    "nick",
		}

		Convey("Quando não é passado argumento", func() {
			cmd.Args = []string{}
			got, err := megasena(cmd)

			So(err, ShouldBeNil)
			So(got, ShouldEqual, fmt.Sprintf("%s: %s", cmd.Nick, msgOpcaoInvalida))
		})

		Convey("Quando o argumento for gerar", func() {
			cmd.Args = []string{"gerar"}
			got, err := megasena(cmd)

			So(err, ShouldBeNil)

			match, err := regexp.MatchString("nick: (\\d{2} {1}){5}\\d{2}", got)

			So(err, ShouldBeNil)
			So(match, ShouldBeTrue)
		})

		Convey("Quando o argumento for resultado", func() {
			ts := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintln(w, retornoJSON)
				}))
			defer ts.Close()

			url = ts.URL

			cmd.Args = []string{"resultado"}
			got, err := megasena(cmd)

			So(err, ShouldBeNil)
			So(got, ShouldEqual, "nick: Sorteio 1636 de 17/09/2014: [19 26 33 35 51 52] - 0 premiado(s) R$ 0,00.")
		})

		Convey("Quando o argumento for resultado e o retorno for inválido", func() {
			ts := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintln(w, "invalid")
				}))
			defer ts.Close()

			url = ts.URL

			cmd.Args = []string{"resultado"}
			_, err := megasena(cmd)

			So(err, ShouldNotBeNil)
		})
	})
}
