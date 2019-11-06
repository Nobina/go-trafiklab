package trafiklab

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	keyDepartures    = "departures"
	keyDeviations    = "deviations"
	keyStopsNearby   = "stops_nearby"
	keyStopsQuery    = "stops_query"
	keyTrafficStatus = "traffic_status"
	keyTravelplanner = "travelplanner"
)

var (
	ErrMissingAPIKey = errors.New("missing api key")
)

var (
	LocationEuropeStockholm, _ = time.LoadLocation("Europe/Stockholm")
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKeys    map[string]string

	Stops         *Stops
	Travelplanner *Travelplanner
}

type ClientOption func(*Client) error

type apiConfig struct {
	host   string
	path   string
	method string
	header map[string]string
}

type apiRequest interface {
	params() url.Values
	body() (io.Reader, error)
}

func (c *Client) do(config *apiConfig, apiReq apiRequest) (*http.Response, error) {
	host := config.host
	if c.baseURL != "" {
		host = c.baseURL
	}

	var body io.Reader
	q := url.Values{}

	if apiReq != nil {
		q = apiReq.params()
		if b, err := apiReq.body(); err != nil {
			return nil, err
		} else {
			body = b
		}
	}

	req, err := http.NewRequest(config.method, host+config.path, body)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = q.Encode()

	if config.header != nil {
		for k, v := range config.header {
			req.Header.Set(k, v)
		}
	}

	return c.httpClient.Do(req)
}

func (c *Client) doJSON(config *apiConfig, apiReq apiRequest, v interface{}) (*http.Response, error) {
	httpResp, err := c.do(config, apiReq)
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	return httpResp, json.NewDecoder(httpResp.Body).Decode(v)
}

func (c *Client) doXML(config *apiConfig, apiReq apiRequest, v interface{}) (*http.Response, error) {
	httpResp, err := c.do(config, apiReq)
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	return httpResp, xml.NewDecoder(httpResp.Body).Decode(v)
}

func NewClient(options ...ClientOption) (*Client, error) {
	c := &Client{
		apiKeys: map[string]string{},
	}

	if options != nil {
		for _, option := range options {
			if err := option(c); err != nil {
				return nil, err
			}
		}
	}

	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	c.Stops = &Stops{c}
	c.Travelplanner = &Travelplanner{c}

	return c, nil
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) error {
		c.httpClient = httpClient
		return nil
	}
}

func WithDeparturesAPIKey(apiKey string) ClientOption {
	return func(c *Client) error {
		c.apiKeys[keyDepartures] = apiKey
		return nil
	}
}

func WithDeviationsAPIKey(apiKey string) ClientOption {
	return func(c *Client) error {
		c.apiKeys[keyDeviations] = apiKey
		return nil
	}
}

func WithStopsNearbyAPIKey(apiKey string) ClientOption {
	return func(c *Client) error {
		c.apiKeys[keyStopsNearby] = apiKey
		return nil
	}
}

func WithStopsQueryAPIKey(apiKey string) ClientOption {
	return func(c *Client) error {
		c.apiKeys[keyStopsQuery] = apiKey
		return nil
	}
}

func WithTrafficStatusAPIKey(apiKey string) ClientOption {
	return func(c *Client) error {
		c.apiKeys[keyTrafficStatus] = apiKey
		return nil
	}
}

func WithTravelplannerAPIKey(apiKey string) ClientOption {
	return func(c *Client) error {
		c.apiKeys[keyTravelplanner] = apiKey
		return nil
	}
}
