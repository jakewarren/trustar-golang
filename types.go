package trustar

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	// APIBaseLive points to the live version of the API and the current stable version that's supported
	APIBaseLive = "https://api.trustar.co/api/1.3/"

	// RequestNewTokenBeforeExpiresIn is used by SendWithAuth and try to get new Token when it's about to expire
	RequestNewTokenBeforeExpiresIn = time.Duration(60) * time.Second
)

type (
	expirationTime int64

	// TokenResponse is for API response for the /oauth2/token endpoint
	TokenResponse struct {
		Token     string         `json:"access_token"`
		Scope     string         `json:"scope"`
		TokenType string         `json:"token_type"`
		ExpiresIn expirationTime `json:"expires_in"`
	}

	// Client represents a TruSTAR REST API Client
	Client struct {
		sync.Mutex
		Client         *http.Client
		ClientID       string
		Secret         string
		APIBase        string
		Log            io.Writer // If user set log file name all requests will be logged there
		Token          *TokenResponse
		tokenExpiresAt time.Time
	}

	// ErrorResponse holds the response if an error occurs
	ErrorResponse struct {
		Response *http.Response
	}

	// Enclave represents enclave data
	Enclave struct {
		Create bool   `json:"create"`
		ID     string `json:"id"`
		Name   string `json:"name"`
		Read   bool   `json:"read"`
		Type   string `json:"type"`
		Update bool   `json:"update"`
	}

	// ReportResponse represents the list of reports matching the user's query
	ReportResponse struct {
		Empty      bool            `json:"empty,omitempty"`
		HasNext    bool            `json:"hasNext"`
		Reports    []ReportDetails `json:"items"`
		PageNumber int64           `json:"pageNumber"`
		PageSize   int64           `json:"pageSize"`
	}

	// CorrelatedReportResponse is the reponse object from the Find Correlated Reports endpoint
	CorrelatedReportResponse struct {
		Empty      bool            `json:"empty"`
		HasNext    bool            `json:"hasNext"`
		Items      []ReportDetails `json:"items"`
		PageNumber int64           `json:"pageNumber"`
		PageSize   int64           `json:"pageSize"`
	}

	// SearchIndicatorReponse contains the search results when searching indicators
	SearchIndicatorReponse struct {
		Empty         bool        `json:"empty"`
		HasNext       bool        `json:"hasNext"`
		Items         []Indicator `json:"items"`
		PageNumber    int64       `json:"pageNumber"`
		PageSize      int64       `json:"pageSize"`
		TotalElements int64       `json:"totalElements"`
		TotalPages    int64       `json:"totalPages"`
	}

	// ReportIndicatorsResponse is the response object from Get Indicators for Report endpoint
	ReportIndicatorsResponse struct {
		Empty      bool        `json:"empty"`
		HasNext    bool        `json:"hasNext"`
		Items      []Indicator `json:"items"`
		PageNumber int64       `json:"pageNumber"`
		PageSize   int64       `json:"pageSize"`
	}

	// RelatedIndicatorsResponse is the response object from Find Related Indicators endpoint
	RelatedIndicatorsResponse struct {
		Empty      bool        `json:"empty"`
		HasNext    bool        `json:"hasNext"`
		Items      []Indicator `json:"items"`
		PageNumber int64       `json:"pageNumber"`
		PageSize   int64       `json:"pageSize"`
	}

	// RequestQuotas represents the current status of the company’s request quotas.
	RequestQuotas []struct {
		GUID          string `json:"guid"`
		LastResetTime int64  `json:"lastResetTime"`
		MaxRequests   int64  `json:"maxRequests"`
		NextResetTime int64  `json:"nextResetTime"`
		TimeWindow    int64  `json:"timeWindow"`
		UsedRequests  int64  `json:"usedRequests"`
	}

	// WhitelistIndicatorsResponse represents the indicators that have been whitelisted by the user's company
	WhitelistIndicatorsResponse struct {
		Empty      bool        `json:"empty"`
		HasNext    bool        `json:"hasNext"`
		Items      []Indicator `json:"items"`
		PageNumber int64       `json:"pageNumber"`
		PageSize   int64       `json:"pageSize"`
	}

	// IndicatorMetadataResponse is a metadata object containing the metadata for the requested indicator(s).
	IndicatorMetadataResponse []struct {
		EnclaveIds    []string       `json:"enclaveIds"`    // the enclaves (of those the user has access to) that the indicator has appeared in a report or indicator submission to
		FirstSeen     int64          `json:"firstSeen"`     // the time (in milliseconds since epoch) that the indicator first appeared in a report or indicator submission to any enclaves the user has access to
		GUID          string         `json:"guid"`          // unique id of the indicator
		IndicatorType string         `json:"indicatorType"` // the type of indicator (IP, URL, EMAIL_ADDRESS, etc.)
		LastSeen      int64          `json:"lastSeen"`      // the time (in milliseconds since epoch) that the indicator last appeared in a report or indicator submission to any enclaves the user has access to
		NoteCount     int64          `json:"noteCount"`     // the number of notes associated with this indicators
		Notes         []string       `json:"notes"`         // the notes associated with the indicator
		PriorityLevel string         `json:"priorityLevel"` // LOW, MEDIUM, or HIGH. NOT_FOUND if no score has been computed for this indicator.
		Sightings     int64          `json:"sightings"`     // the number of times the indicator has appeared in a report or indicator submission to any enclaves the user has access to
		Tags          []IndicatorTag `json:"tags"`          // the set of Tag objects that the indicator has been tagged with
		Value         string         `json:"value"`         // indicator value
	}

	// Indicator hold the indicator metadata search queries
	Indicator struct {
		GUID          string `json:"guid,omitempty"`
		IndicatorType string `json:"indicatorType,omitempty"` // the type of indicator (IP, URL, EMAIL_ADDRESS, etc.)
		PriorityLevel string `json:"priorityLevel,omitempty"` // LOW, MEDIUM, or HIGH. NOT_FOUND if no score has been computed for this indicator.
		Value         string `json:"value"`                   // the indicator’s value
		Weight        int    `json:"weight"`                  // Possible values are 0 and 1. A value of 0 indicates that, although the term fits the technical requirements to be considered an indicator, our machine learning model has determined that it is likely not an indicator of compromise when considered in the context of a specific report.
		Reason        string `json:"reason"`                  // the reason the indicator has a weight of 0 (not present if weight is 1)
		Whitelisted   string `json:"whitelisted"`             // whether the indicator has been whitelisted by the requesting company
	}

	// TrendingIndicators holds a list of Indicator objects that are trending in the community
	TrendingIndicators []struct {
		CorrelationCount int64 `json:"correlationCount"`
		Indicator
	}

	// IndicatorSubmission records the information needed to submit a new indicator to TruSTAR
	IndicatorSubmission struct {
		EnclaveIDS []string           `json:"enclaveIds"`
		Content    []IndicatorContent `json:"content"`
		Tags       []IndicatorTag     `json:"tags,omitempty"`
	}

	// IndicatorContent details an indicator object when submitting a new indicator
	IndicatorContent struct {
		Value     string         `json:"value"`
		FirstSeen int64          `json:"firstSeen,omitempty"`
		LastSeen  int64          `json:"lastSeen,omitempty"`
		Sightings int64          `json:"sightings,omitempty"`
		Source    string         `json:"source,omitempty"`
		Notes     string         `json:"notes,omitempty"`
		Tags      []IndicatorTag `json:"tags,omitempty"`
	}

	// IndicatorTag contains the tag information
	IndicatorTag struct {
		Guid      string `json:"guid"`      // the ID of the tag
		Name      string `json:"name"`      // the name of the tag (i.e. the actual string value of the tag)
		EnclaveID string `json:"enclaveId"` // the ID of the enclave of the tag
	}

	// ReportSubmission is used for submitting a new report
	ReportSubmission struct {
		DistributionType   string   `json:"distributionType"`             // [required] COMMUNITY (will disregard any enclaveIds) or ENCLAVE (must include enclaveIds)
		EnclaveIds         []string `json:"enclaveIds"`                   // Non-empty array of TruSTAR-generated enclave ids (available on Station under settings or through the GET /enclaves endpoint). Use the enclave ID, NOT the enclave name.
		ExternalTrackingID string   `json:"externalTrackingId,omitempty"` // External tracking ID provided by user. Must be unique across all reports for a given company.
		ExternalURL        string   `json:"externalUrl,omitempty"`        // URL for the external report that this originated from, if one exists. Limit 500 alphanumeric characters.
		ReportBody         string   `json:"reportBody"`                   // [required] Text content of report
		TimeBegan          string   `json:"timeBegan,omitempty"`          // SO-8601 formatted incident time with timezone, e.g. 2016-09-22T11:38:35+00:00
		Title              string   `json:"title"`                        // [required] Title of the report
	}

	// ReportDetails contains the details for a specific report
	ReportDetails struct {
		Created          int64    `json:"created"`          // the time of creation, in milliseconds since epoch
		DistributionType string   `json:"distributionType"` // ENCLAVE or COMMUNITY - if COMMUNITY, the report is open to the community. This field is deprecated, but is retained for backwards compatibility. The Community has been transitioned to an enclave, so all reports have a distributionType of ENCLAVE.
		EnclaveIds       []string `json:"enclaveIds"`       // the list of IDs of the enclaves that the report has been submitted to
		ExternalID       string   `json:"externalId"`       // the external ID of the report (any string, user-defined)
		ID               string   `json:"id"`               // the internal ID of the report (a GUID)
		ReportBody       string   `json:"reportBody"`       // the body of the report
		Sector           Sector   `json:"sector"`           // the company's sector
		TimeBegan        int64    `json:"timeBegan"`        // the user-defined time when the incident began, in milliseconds since epoch
		Title            string   `json:"title"`            // the report title
		Updated          int64    `json:"updated"`          // the time of the last update, in milliseconds since epoch
	}

	// Sector records sector information for reports and indicators
	Sector struct {
		ID    int64  `json:"id"`    // the ID of the company’s sector
		Label string `json:"label"` // the name of the company’s sector
		Name  string `json:"name"`  // the label of the company’s sector
	}
)

// Error method implementation for ErrorResponse struct
func (r *ErrorResponse) Error() string {
	respBody, _ := ioutil.ReadAll(r.Response.Body)
	defer r.Response.Body.Close()
	return fmt.Sprintf("%v %v: %d %s", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, respBody)
}
