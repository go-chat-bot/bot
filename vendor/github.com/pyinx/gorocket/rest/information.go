package rest

import (
	"github.com/pyinx/gorocket/api"
	"net/http"
)

type infoResponse struct {
	Info api.Info `json:"info"`
}

// Get information about the server.
// This function does not need a logged in user.
//
// https://rocket.chat/docs/developer-guides/rest-api/miscellaneous/info
func (c *Client) GetServerInfo() (*api.Info, error) {
	request, _ := http.NewRequest("GET", c.getUrl()+"/api/v1/info", nil)

	response := new(infoResponse)

	if err := c.doRequest(request, response); err != nil {
		return nil, err
	}

	return &response.Info, nil
}
