package trustar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Reference: https://docs.trustar.co/api/v13/reports/index.html

// GetReports Returns a page of incident reports matching the specified filters. All parameters are optional: if nothing is specified, the latest 25 reports accessible by the user will be returned.
//
// Endpoint: GET /1.3/reports
func (c *Client) GetReports(v url.Values) (ReportResponse, error) {
	var rr ReportResponse

	url := fmt.Sprintf("%s%s", c.APIBase, fmt.Sprintf("reports?%s", v.Encode()))
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return rr, err
	}

	if err = c.SendWithAuth(req, &rr); err != nil {
		return rr, err
	}

	return rr, nil
}

// GetReportIndicators Returns a paginated list of all indicators contained in a specified report.
//
// Endpoint: GET /1.3/reports/{id}/indicators
func (c *Client) GetReportIndicators(id string, v url.Values) (ReportIndicatorsResponse, error) {
	var rir ReportIndicatorsResponse

	url := fmt.Sprintf("%s%s", c.APIBase, fmt.Sprintf("reports/%s/indicators?%s", id, v.Encode()))
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return rir, err
	}

	if err = c.SendWithAuth(req, &rir); err != nil {
		return rir, err
	}

	return rir, nil
}

// FindCorrelatedReports Returns a paginated list of all reports that contain any of the provided indicator values.
//
// Endpoint: GET /1.3/reports/correlated
func (c *Client) FindCorrelatedReports(v url.Values) (CorrelatedReportResponse, error) {
	var crr CorrelatedReportResponse

	url := fmt.Sprintf("%s%s", c.APIBase, fmt.Sprintf("reports/correlated?%s", v.Encode()))
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return crr, err
	}

	if err = c.SendWithAuth(req, &crr); err != nil {
		return crr, err
	}

	return crr, nil
}

// SubmitReport Submit a new incident report, and receive the ID it has been assigned in TruSTAR’s system.
//
// Endpoint: POST /1.3/reports
func (c *Client) SubmitReport(report ReportSubmission) (string, error) {

	var guid strings.Builder

	i, _ := json.Marshal(report)

	url := fmt.Sprintf("%s%s", c.APIBase, "reports")
	req, err := http.NewRequest("POST", url, bytes.NewReader(i))

	if err != nil {
		return "", err
	}

	if err = c.SendWithAuth(req, &guid); err != nil {
		return "", err
	}

	return guid.String(), nil
}

// UpdateReport Update the report with the specified Trustar report ID.
//
// Endpoint: PUT /1.3/reports/{ID}
func (c *Client) UpdateReport(id string, report ReportSubmission) error {

	i, _ := json.Marshal(report)

	url := fmt.Sprintf("%s%s", c.APIBase, fmt.Sprintf("reports/%s", id))
	req, err := http.NewRequest("PUT", url, bytes.NewReader(i))

	if err != nil {
		return err
	}

	return c.SendWithAuth(req, ioutil.Discard)
}

// GetReportDetails Gets the details for a report for the specified Trustar report id
//
// Endpoint: GET /1.3/reports/{ID}
func (c *Client) GetReportDetails(id string) (ReportDetails, error) {

	var rd ReportDetails

	url := fmt.Sprintf("%s%s", c.APIBase, fmt.Sprintf("reports/%s", id))
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return rd, err
	}

	if err = c.SendWithAuth(req, &rd); err != nil {
		return rd, err
	}

	return rd, nil
}

// DeleteReport Delete a report with the specified Trustar report ID.
//
// Endpoint: DELETE /1.3/reports/{ID}
func (c *Client) DeleteReport(id string) error {

	url := fmt.Sprintf("%s%s", c.APIBase, fmt.Sprintf("reports/%s", id))
	req, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		return err
	}

	return c.SendWithAuth(req, ioutil.Discard)
}

// SearchReports Searches for all reports that contain the given search term.
//
// Endpoint: GET /1.3/reports/search
func (c *Client) SearchReports(v url.Values) (ReportResponse, error) {
	var rr ReportResponse

	url := fmt.Sprintf("%s%s", c.APIBase, fmt.Sprintf("reports/search?%s", v.Encode()))
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return rr, err
	}

	if err = c.SendWithAuth(req, &rr); err != nil {
		return rr, err
	}

	return rr, nil
}
