package trustar

import (
	"fmt"
	"net/http"
)

// Reference: https://docs.trustar.co/api/v13/enclaves/index.html

// GetEnclaves Returns the list of all enclaves that the user has access to, as well as whether they can read, create, and update reports in that enclave.
//
// Endpoint: GET /1.3/enclaves
func (c *Client) GetEnclaves() ([]Enclave, error) {
	var enclaves []Enclave

	url := fmt.Sprintf("%s%s", c.APIBase, "enclaves")
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return enclaves, err
	}

	if err = c.SendWithAuth(req, &enclaves); err != nil {
		return enclaves, err
	}

	return enclaves, nil
}
