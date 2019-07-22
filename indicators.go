package trustar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Reference: https://docs.trustar.co/api/v13/indicators/index.html

// SearchIndicators Searches for all indicators that contain the given search term. Also allows filtering by date, enclave, and tags.
//
// Endpoint: GET /1.3/indicators/search
func (c *Client) SearchIndicators(v url.Values) (SearchIndicatorReponse, error) {
	var sir SearchIndicatorReponse

	url := fmt.Sprintf("%s%s", c.APIBase, fmt.Sprintf("indicators/search?%s", v.Encode()))
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return sir, err
	}

	if err = c.SendWithAuth(req, &sir); err != nil {
		return sir, err
	}

	return sir, nil
}

// FindRelatedIndicators Search all TruSTAR incident reports for provided indicators and return all correlated indicators from search results. Two indicators are considered “correlated” if they can be found in a common report.
//
// Endpoint: GET /1.3/indicators/related
func (c *Client) FindRelatedIndicators(v url.Values) (RelatedIndicatorsResponse, error) {
	var rir RelatedIndicatorsResponse

	url := fmt.Sprintf("%s%s", c.APIBase, fmt.Sprintf("indicators/related?%s", v.Encode()))
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return rir, err
	}

	if err = c.SendWithAuth(req, &rir); err != nil {
		return rir, err
	}

	return rir, nil
}

// WhitelistIndicators Whitelist a list of indicator values for the user’s company.
//
// Endpoint: POST /1.3/whitelist
func (c *Client) WhitelistIndicators(indicators []string) error {

	var wr interface{}

	i, _ := json.Marshal(indicators)

	url := fmt.Sprintf("%s%s", c.APIBase, "whitelist")
	req, err := http.NewRequest("POST", url, bytes.NewReader(i))

	if err != nil {
		return err
	}

	if err = c.SendWithAuth(req, &wr); err != nil {
		return err
	}

	return nil
}

// GetWhitelist Get a paginated list of the indicators that have been whitelisted by the user’s company.
//
// Endpoint: GET /1.3/whitelist
func (c *Client) GetWhitelist(v url.Values) (WhitelistIndicatorsResponse, error) {
	var wir WhitelistIndicatorsResponse

	url := fmt.Sprintf("%s%s", c.APIBase, fmt.Sprintf("whitelist?%s", v.Encode()))
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return wir, err
	}

	if err = c.SendWithAuth(req, &wir); err != nil {
		return wir, err
	}

	return wir, nil
}

// DeleteFromWhitelist Delete an indicator from the user’s company whitelist.
//
// Endpoint: DELETE /1.3/whitelist
func (c *Client) DeleteFromWhitelist(v url.Values) error {
	var wr interface{}

	url := fmt.Sprintf("%s%s", c.APIBase, fmt.Sprintf("whitelist?%s", v.Encode()))
	req, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		return err
	}

	if err = c.SendWithAuth(req, &wr); err != nil {
		// this endpoint returns no content so ignore EOF errors
		if err.Error() != "EOF" {
			return err
		}
	}

	return nil
}

// GetIndicatorMetadata Provide metadata associated with an indicator
//
// Endpoint: POST /1.3/indicators/metadata
func (c *Client) GetIndicatorMetadata(indicators []Indicator) (IndicatorMetadataResponse, error) {

	var imr IndicatorMetadataResponse

	i, _ := json.Marshal(indicators)

	url := fmt.Sprintf("%s%s", c.APIBase, "indicators/metadata")
	req, err := http.NewRequest("POST", url, bytes.NewReader(i))

	if err != nil {
		return nil, err
	}

	if err = c.SendWithAuth(req, &imr); err != nil {
		return nil, err
	}

	return imr, nil
}

// GetTrendingIndicators Returns the 10 indicators that have recently appeared in the most community reports.
//
// Endpoint: GET /1.3/indicators/community-trending
func (c *Client) GetTrendingIndicators(v url.Values) (TrendingIndicators, error) {
	var ti TrendingIndicators

	url := fmt.Sprintf("%s%s", c.APIBase, fmt.Sprintf("indicators/community-trending?%s", v.Encode()))
	req, err := http.NewRequest("GET", url, nil)

	fmt.Println(req.URL)

	if err != nil {
		return nil, err
	}

	if err = c.SendWithAuth(req, &ti); err != nil {
		// this endpoint returns no content so ignore EOF errors
		if err.Error() != "EOF" {
			return nil, err
		}
	}

	return ti, nil
}

// SubmitIndicators Submit Indicators
//
// Endpoint: POST /1.3/indicators
func (c *Client) SubmitIndicators(indicators IndicatorSubmission) error {

	var imr interface{}

	i, _ := json.Marshal(indicators)

	url := fmt.Sprintf("%s%s", c.APIBase, "indicators")
	req, err := http.NewRequest("POST", url, bytes.NewReader(i))

	if err != nil {
		return err
	}

	if err = c.SendWithAuth(req, &imr); err != nil {
		return err
	}

	return nil
}
