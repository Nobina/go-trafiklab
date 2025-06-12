package journeyplanner

import (
	"errors"
	"net/http"
)

const (
	baseURL = "https://journeyplanner.integration.sl.se/v2"
)



type JourneyPlannerConfig struct {
	APIKey string // TODO(pb): Might not be needed anymore
}

func (c *JourneyPlannerConfig) Valid() error {
	if c.APIKey == "" {
		return errors.New("api key is required")
	}
	return nil
}

type Option func(*Client)

func WithDebug() Option {
	return func(c *Client) {
		c.isDebug = true
	}
}

func WithCustomBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey string
	isDebug      bool
}

func NewClient(cfg *JourneyPlannerConfig, client *http.Client, opts ...Option) *Client {
	tc := &Client{
		httpClient: client,
		apiKey:     cfg.APIKey,
		baseURL:    baseURL,
	}

	for _, opt := range opts {
		opt(tc)
	}

	return tc
}

type SystemMessage struct {
	// Type of message, e.g. "error"
	Type string `json:"type"`
	// Back-end module reporting the message
	Module string `json:"module"`
	// Internal error code
	Code int `json:"code"`
	// Description of error, if available
	Text string `json:"text"`
	// SubType of error, if available
	SubType string `json:"subType"`
}
