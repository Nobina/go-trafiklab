package deviations

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/nobina/go-trafiklab/requests"
)

type Config struct {
	BaseURL string
}

func (cfg *Config) Valid() error {
	if cfg.BaseURL == "" {
		return fmt.Errorf("missing base url")
	}
	return nil
}

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(cfg *Config, client *http.Client) *Client {
	return &Client{
		httpClient: client,
		baseURL:    cfg.BaseURL,
	}
}

func (c *Client) Deviations(ctx context.Context, payload *DeviationsRequest) ([]*DeviationsResponse, error) {
	url := c.baseURL + "/v1/messages"

	req, err := requests.JSON(ctx, http.MethodGet, url, nil)
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
		log.Printf("unexpected status code: %d", res.StatusCode)
		log.Printf("url: %s\n", url+"?"+req.URL.RawQuery)
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	deviationsResp := []*DeviationsResponse{}

	err = json.NewDecoder(res.Body).Decode(&deviationsResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return deviationsResp, nil
}

type DeviationsRequest struct {
	Future             bool     `json:"future"`
	TransportAuthority int      `json:"transport_authority"`
	LineNumbers        []int    `json:"line_number"`
	TransportModes     []string `json:"transport_mode"`
	SiteIDs            []int    `json:"site_id"`
}

func (r DeviationsRequest) params() url.Values {
	params := url.Values{}
	if len(r.TransportModes) > 0 {
		for _, v := range r.TransportModes {
			params.Add("transport_mode", v)
		}
	}
	if len(r.LineNumbers) > 0 {
		for _, v := range r.LineNumbers {
			params.Add("line", strconv.Itoa(v))
		}
	}
	if len(r.SiteIDs) > 0 {
		for _, v := range r.SiteIDs {
			params.Add("site", strconv.Itoa(v))
		}
	}
	if r.Future {
		params.Set("future", "true")
	}
	if r.TransportAuthority != 0 {
		params.Set("transport_authority", strconv.Itoa(r.TransportAuthority))
	}
	return params
}

type DeviationsResponse struct {
	Version         int               `json:"version"`
	Created         time.Time         `json:"created"`
	Modified        time.Time         `json:"modified"`
	DeviationCaseID int               `json:"deviation_case_id"`
	Publish         Publish           `json:"publish"`
	Priority        Priority          `json:"priority"`
	MessageVariants []MessageVariants `json:"message_variants"`
	Scope           Scope             `json:"scope"`
}
type Publish struct {
	From time.Time `json:"from"`
	Upto time.Time `json:"upto"`
}
type Priority struct {
	ImportanceLevel int `json:"importance_level"`
	InfluenceLevel  int `json:"influence_level"`
	UrgencyLevel    int `json:"urgency_level"`
}
type MessageVariants struct {
	Header     string `json:"header"`
	Details    string `json:"details"`
	ScopeAlias string `json:"scope_alias"`
	Weblink    string `json:"weblink"`
	Language   string `json:"language"`
}
type StopPoints struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type StopAreas struct {
	ID                 int          `json:"id"`
	TransportAuthority int          `json:"transport_authority"`
	Name               string       `json:"name"`
	Type               string       `json:"type"`
	StopPoints         []StopPoints `json:"stop_points"`
}
type Lines struct {
	ID                 int    `json:"id"`
	TransportAuthority int    `json:"transport_authority"`
	Designation        string `json:"designation"`
	TransportMode      string `json:"transport_mode"`
	Name               string `json:"name"`
	GroupOfLines       string `json:"group_of_lines"`
}
type Scope struct {
	StopAreas []StopAreas `json:"stop_areas"`
	Lines     []Lines     `json:"lines"`
}
