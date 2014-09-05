package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type GetBodyFunc func(string) ([]byte, error)
type GetJSONFunc func(string, interface{}) error

func GetBody(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func GetJSON(url string, v interface{}) error {
	body, err := GetBody(url)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

func GetJSONFromString(text string, setURL func(string)) func(string, interface{}) error {
	return func(url string, v interface{}) error {
		setURL(url)
		return json.Unmarshal([]byte(text), v)
	}
}
