package trustar

import (
	"fmt"
	"net/http"
	"strings"
)

// Ping is a standard ping endpoint, used to help users easily establish that the API can be reached.
// Reference: https://docs.trustar.co/api/v13/ping.html
//
// Endpoint: GET /1.3/ping
func (c *Client) Ping() (string, error) {
	var response strings.Builder

	url := fmt.Sprintf("%s%s", c.APIBase, "ping")
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}

	if err = c.SendWithAuth(req, &response); err != nil {
		return "", err
	}

	return response.String(), nil
}

// Version returns the current stable version
// Reference: https://docs.trustar.co/api/v13/version.html
//
// Endpoint: GET /api/version
func (c *Client) Version() (string, error) {
	var response strings.Builder

	req, err := http.NewRequest("GET", "https://api.trustar.co/api/version", nil)

	if err != nil {
		return "", err
	}

	if err = c.SendWithAuth(req, &response); err != nil {
		return "", err
	}

	return response.String(), nil
}

// RequestQuotas returns the current status of the company’s request quotas. A request quota is a maximum number of requests that a company can send to the API during a given time window.
// Reference: https://docs.trustar.co/api/v13/request_quotas.html
//
// Endpoint: GET /1.3/request-quotas
func (c *Client) RequestQuotas() (RequestQuotas, error) {
	var quotas RequestQuotas

	url := fmt.Sprintf("%s%s", c.APIBase, "request-quotas")
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return quotas, err
	}

	if err = c.SendWithAuth(req, &quotas); err != nil {
		return quotas, err
	}

	return quotas, nil
}
