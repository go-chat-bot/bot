package rest

import (
	"bytes"
	"errors"
	"github.com/pyinx/gorocket/api"
	"net/http"
	"net/url"
)

type logoutResponse struct {
	statusResponse
	data struct {
		message string `json:"message"`
	} `json:"data"`
}

type logonResponse struct {
	statusResponse
	Data struct {
		Token  string `json:"authToken"`
		UserId string `json:"userId"`
	} `json:"data"`
}

// Login a user. The Email and the Password are mandatory. The auth token of the user is stored in the Client instance.
//
// https://rocket.chat/docs/developer-guides/rest-api/authentication/login
func (c *Client) Login(credentials api.UserCredentials) error {
	data := url.Values{"user": {credentials.Email}, "password": {credentials.Password}}
	request, _ := http.NewRequest("POST", c.getUrl()+"/api/v1/login", bytes.NewBufferString(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response := new(logonResponse)

	if err := c.doRequest(request, response); err != nil {
		return err
	}

	if response.Status == "success" {
		c.auth = &authInfo{id: response.Data.UserId, token: response.Data.Token}
		return nil
	} else {
		return errors.New("Response status: " + response.Status)
	}
}

// Logout a user. The function returns the response message of the server.
//
// https://rocket.chat/docs/developer-guides/rest-api/authentication/logout
func (c *Client) Logout() (string, error) {

	if c.auth == nil {
		return "Was not logged in", nil
	}

	request, _ := http.NewRequest("POST", c.getUrl()+"/api/v1/logout", nil)

	response := new(logoutResponse)

	if err := c.doRequest(request, response); err != nil {
		return "", err
	}

	if response.Status == "success" {
		return response.data.message, nil
	} else {
		return "", errors.New("Response status: " + response.Status)
	}
}
