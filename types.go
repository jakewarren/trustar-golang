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
		EnclaveIds    []string       `json:"enclaveIds"`
		FirstSeen     int64          `json:"firstSeen"`
		GUID          string         `json:"guid"`
		IndicatorType string         `json:"indicatorType"`
		LastSeen      int64          `json:"lastSeen"`
		NoteCount     int64          `json:"noteCount"`
		Notes         []string       `json:"notes"`
		PriorityLevel string         `json:"priorityLevel"`
		Sightings     int64          `json:"sightings"`
		Tags          []IndicatorTag `json:"tags"`
		Value         string         `json:"value"`
	}

	// Indicator hold the indicator metadata search queries
	Indicator struct {
		GUID          string `json:"guid,omitempty"`
		IndicatorType string `json:"indicatorType,omitempty"`
		PriorityLevel string `json:"priorityLevel,omitempty"`
		Value         string `json:"value"`
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
		Guid      string `json:"guid"`
		Name      string `json:"name"`
		EnclaveID string `json:"enclaveId"`
	}

	// ReportSubmission is used for submitting a new report
	ReportSubmission struct {
		DistributionType   string   `json:"distributionType"`
		EnclaveIds         []string `json:"enclaveIds"`
		ExternalTrackingID string   `json:"externalTrackingId,omitempty"`
		ExternalURL        string   `json:"externalUrl,omitempty"`
		ReportBody         string   `json:"reportBody"`
		TimeBegan          string   `json:"timeBegan,omitempty"`
		Title              string   `json:"title"`
	}

	// ReportDetails contains the details for a specific report
	ReportDetails struct {
		Created          int64    `json:"created"`
		DistributionType string   `json:"distributionType"`
		EnclaveIds       []string `json:"enclaveIds"`
		ExternalID       string   `json:"externalId"`
		ID               string   `json:"id"`
		ReportBody       string   `json:"reportBody"`
		Sector           Sector   `json:"sector"`
		TimeBegan        int64    `json:"timeBegan"`
		Title            string   `json:"title"`
		Updated          int64    `json:"updated"`
	}

	// Sector records sector information for reports and indicators
	Sector struct {
		ID    int64  `json:"id"`
		Label string `json:"label"`
		Name  string `json:"name"`
	}
)

// Error method implementation for ErrorResponse struct
func (r *ErrorResponse) Error() string {
	respBody, _ := ioutil.ReadAll(r.Response.Body)
	defer r.Response.Body.Close()
	return fmt.Sprintf("%v %v: %d %s", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, respBody)
}
