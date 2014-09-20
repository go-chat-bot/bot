package cotacao

import (
	"fmt"
	"github.com/fabioxgn/go-bot"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	expectedJSON = `{
      "bovespa":{
        "cotacao":"60800",
        "variacao":"-1.68"
      },
      "dolar":{
        "cotacao":"2.2430",
        "variacao":"+0.36"
      },
      "euro":{
        "cotacao":"2.9018",
        "variacao":"-1.21"
      },
      "atualizacao":"04\/09\/14   -18:13"
    }`
)

func TestCotacao(t *testing.T) {

	Convey("Ao executar o comando cotação", t, func() {
		cmd := &bot.Cmd{}

		Convey("Deve responder com a cotação do dólar e euro", func() {
			ts := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintln(w, expectedJSON)
				}))
			defer ts.Close()

			url = ts.URL

			c, err := cotacao(cmd)

			So(err, ShouldBeNil)
			So(c, ShouldEqual, "Dólar: 2.2430 (+0.36), Euro: 2.9018 (-1.21)")
		})

		Convey("Quando o webservice retornar algo inválido deve retornar erro", func() {
			ts := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintln(w, "invalid")
				}))
			defer ts.Close()

			url = ts.URL

			_, err := cotacao(cmd)

			So(err, ShouldNotBeNil)
		})
	})
}
