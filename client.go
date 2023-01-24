package trafiklab

import (
	"errors"
	"net/http"
	"time"

	"github.com/nobina/go-requester"
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
	client         *requester.Client
	clientOptions  []requester.ClientOption
	defaultOptions []requester.RequestOption
	apiKeys        map[string]string

	Stops         *Stops
	Travelplanner *Travelplanner
}

type ClientOption func(*Client)

func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		clientOptions: []requester.ClientOption{
			requester.WithHTTPClient(http.DefaultClient),
		},
		defaultOptions: []requester.RequestOption{
			requester.WithHost("http://api.sl.se"),
		},
		apiKeys: map[string]string{},
	}

	for _, opt := range opts {
		opt(c)
	}

	c.clientOptions = append(c.clientOptions, requester.WithDefaultOptions(c.defaultOptions...))
	c.client = requester.NewClient(c.clientOptions...)
	c.Stops = &Stops{c}
	c.Travelplanner = &Travelplanner{c}

	return c
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) { c.clientOptions = append(c.clientOptions, requester.WithHTTPClient(httpClient)) }
}

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) { c.defaultOptions = append(c.defaultOptions, requester.WithHost(baseURL)) }
}

func WithDeparturesAPIKey(apiKey string) ClientOption {
	return func(c *Client) { c.apiKeys[keyDepartures] = apiKey }
}

func WithDeviationsAPIKey(apiKey string) ClientOption {
	return func(c *Client) { c.apiKeys[keyDeviations] = apiKey }
}

func WithStopsNearbyAPIKey(apiKey string) ClientOption {
	return func(c *Client) { c.apiKeys[keyStopsNearby] = apiKey }
}

func WithStopsQueryAPIKey(apiKey string) ClientOption {
	return func(c *Client) { c.apiKeys[keyStopsQuery] = apiKey }
}

func WithTrafficStatusAPIKey(apiKey string) ClientOption {
	return func(c *Client) { c.apiKeys[keyTrafficStatus] = apiKey }
}

func WithTravelplannerAPIKey(apiKey string) ClientOption {
	return func(c *Client) { c.apiKeys[keyTravelplanner] = apiKey }
}
