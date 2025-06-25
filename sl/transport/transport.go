package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/nobina/go-trafiklab/requests"
	"github.com/nobina/go-trafiklab/slidentifiers"
)

const (
	TransportModeBus   = "BUS"
	TransportModeTram  = "TRAM"
	TransportModeMetro = "METRO"
	TransportModeTrain = "TRAIN"
	TransportModeFerry = "FERRY"
	TransportModeShip  = "SHIP"
	TransportModeTaxi  = "TAXI"
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
	isDebug    bool
}

func NewClient(cfg *Config, client *http.Client, options ...Option) *Client {
	c := &Client{
		httpClient: client,
		baseURL:    cfg.BaseURL,
	}

	for _, opt := range options {
		opt(c)
	}
	return c
}

type Option func(*Client)

func WithDebug() Option {
	return func(c *Client) {
		c.isDebug = true
	}
}

func (c *Client) Departures(ctx context.Context, payload *DeparturesRequest) (*DepartureResponse, error) {

	var err error
	if len(payload.SiteID) == 16 {
		payload.SiteID, err = slidentifiers.ConvertEFAtoSiteID(payload.SiteID)
		if err != nil {
			return nil, fmt.Errorf("failed to convert site id to old id: %w", err)
		}
	}

	url := fmt.Sprintf("%s/v1/sites/%s/departures", c.baseURL, payload.SiteID)
	req, err := requests.JSON(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	q := payload.params()
	req.URL.RawQuery = q.Encode()

	if c.isDebug {
		log.Printf("url: %s\n", url+req.URL.RawQuery)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed request: %w", err)
	}
	defer resp.Body.Close()

	if c.isDebug {
		res, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Printf("failed to dump response: %v", err)
		} else {
			log.Printf("response: %s\n", res)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, response: %v, for url: %s", resp.StatusCode, resp, url+req.URL.RawQuery)
	}

	departuresResp := &DepartureResponse{}
	err = json.NewDecoder(resp.Body).Decode(departuresResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w, for url: %s", err, url+req.URL.RawQuery)
	}

	departuresResp = filterTransportTypes(departuresResp, payload.Bus, payload.Metro, payload.Train, payload.Tram, payload.Ship)
	return departuresResp, nil
}

// The new API for SL doesn't support multiple filters so we will have to do it ourselves...
func filterTransportTypes(res *DepartureResponse, bus, metro, train, tram, ship bool) *DepartureResponse {
	if bus && metro && train && tram && ship {
		return res
	}
	var departures []*Departure
	for _, departure := range res.Departures {
		transportMode := departure.Line.TransportMode
		if bus && transportMode == TransportModeBus {
			departures = append(departures, departure)
		}
		if metro && transportMode == TransportModeMetro {
			departures = append(departures, departure)
		}
		if train && transportMode == TransportModeTrain {
			departures = append(departures, departure)
		}
		if tram && transportMode == TransportModeTram {
			departures = append(departures, departure)
		}
		if ship && transportMode == TransportModeShip {
			departures = append(departures, departure)
		}
	}

	res.Departures = departures
	return res
}

type DeparturesRequest struct {
	SiteID   string `json:"site_id"`
	Forecast int    `json:"time_window"`
	Bus      bool   `json:"bus"`
	Metro    bool   `json:"metro"`
	Train    bool   `json:"train"`
	Tram     bool   `json:"tram"`
	Ship     bool   `json:"ship"`
}

func (r DeparturesRequest) params() url.Values {
	params := url.Values{}
	if r.Forecast != 0 {
		params.Set("forecast", strconv.Itoa(r.Forecast))
	}
	return params
}

type DepartureResponse struct {
	Departures     []*Departure      `json:"departures"`
	StopDeviations []*StopDeviations `json:"stop_deviations"`
}
type Journey struct {
	ID              int64  `json:"id"`
	State           string `json:"state"`
	PredictionState string `json:"prediction_state"`
	PassengerLevel  string `json:"passenger_level"`
}
type StopArea struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Sname string `json:"sname"`
	Type  string `json:"type"`
}
type StopPoint struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Designation string `json:"designation"`
}
type Line struct {
	ID            int    `json:"id"`
	Designation   string `json:"designation"`
	TransportMode string `json:"transport_mode"`
	GroupOfLines  string `json:"group_of_lines"`
}
type Departure struct {
	Direction     string               `json:"direction"`
	DirectionCode int                  `json:"direction_code"`
	Via           string               `json:"via"`
	Destination   string               `json:"destination"`
	State         string               `json:"state"`
	Scheduled     string               `json:"scheduled"`
	Expected      string               `json:"expected"`
	Display       string               `json:"display"`
	Journey       Journey              `json:"journey"`
	StopArea      StopArea             `json:"stop_area"`
	StopPoint     StopPoint            `json:"stop_point"`
	Line          Line                 `json:"line"`
	Deviations    []DepartureDeviation `json:"deviations"`
}
type StopDeviations struct {
	Importance  int    `json:"importance"`
	Consequence string `json:"consequence"`
	Message     string `json:"message"`
}

type DepartureDeviation struct {
	Consequence     string `json:"consequence"`
	ImportanceLevel int    `json:"importance_level"`
	Message         string `json:"message"`
}
