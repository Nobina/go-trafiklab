package stops

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"gopkg.in/DataDog/dd-trace-go.v1/internal/log"
)

type Config struct {
	APIKey  string
	BaseURL string
}

func (cfg *Config) Valid() error {
	if cfg.APIKey == "" {
		return errors.New("missing api key")
	}
	if cfg.BaseURL == "" {
		return errors.New("missing base url")
	}
	return nil
}

type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func NewClient(cfg *Config, client *http.Client) *Client {
	return &Client{
		httpClient: client,
		apiKey:     cfg.APIKey,
		baseURL:    cfg.BaseURL,
	}
}

func (c *Client) Query(ctx context.Context, payload *StopsQueryRequest) (*TypeaheadResponse, error) {
	payload.key = c.apiKey
	url := c.baseURL + "/v1/typeahead.xml"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	q := payload.params()
	req.URL.RawQuery = q.Encode()

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Errorf("unexpected status code: %d", res.StatusCode)
		log.Errorf("response: %v", res)
		return nil, fmt.Errorf("unexpected status code: %d, response: %v, for url: %s", res.StatusCode, res, url+req.URL.RawQuery)
	}

	queryResp := &TypeaheadResponse{}
	err = xml.NewDecoder(res.Body).Decode(queryResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return queryResp, nil
}

type StopsQueryRequest struct {
	key string

	SearchString string
	StationsOnly bool
	MaxResults   string
	Type         string
}

func (r StopsQueryRequest) params() url.Values {
	params := url.Values{}
	if r.key != "" {
		params.Set("key", r.key)
	}
	if r.SearchString != "" {
		params.Set("SearchString", r.SearchString)
	}
	if r.StationsOnly {
		params.Set("StationsOnly", "True")
	} else {
		params.Set("StationsOnly", "False")
	}
	if r.MaxResults != "" {
		params.Set("MaxResults", r.MaxResults)
	}
	if r.Type != "" {
		params.Set("type", r.Type)
	}
	return params
}

type TypeaheadResponse struct {
	StatusCode    int32           `json:"statusCode"`
	Message       string          `json:"message"`
	ExecutionTime int64           `json:"executionTime"`
	Data          []TypeaheadStop `json:"stops" xml:"ResponseData>Site"`
}

type TypeaheadStop struct {
	Name   string `json:"name"`
	SiteID string `xml:"SiteId" json:"siteId"`
	Type   string `json:"type"`
	X      string `json:"x"`
	Y      string `json:"y"`
}
