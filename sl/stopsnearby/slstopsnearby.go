package stopsnearby

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"

	"github.com/nobina/go-trafiklab/slidentifiers"
)

const (
	EFAPrefix = "909100100"
)

type Config struct {
	APIKey  string
	BaseURL string
}

type StopsNearbyClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

// Deprecated
func NewStopsNearbyClientDeprecated(cfg *Config, client *http.Client) *StopsNearbyClient {
	return &StopsNearbyClient{
		httpClient: client,
		apiKey:     cfg.APIKey,
		baseURL:    cfg.BaseURL,
	}
}

// Deprecated
func (c *StopsNearbyClient) Nearby(ctx context.Context, body *StopsNearbyRequest) (*LocationList, error) {
	url := c.baseURL + "/nearbystopsv2.xml"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	q := body.params()
	q.Add("key", c.apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed request: %w", err)
	}
	defer resp.Body.Close()

	nearbyResp := &LocationList{}
	err = xml.NewDecoder(resp.Body).Decode(nearbyResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	filtered := []*StopLocation{}
	for _, stop := range nearbyResp.Data {
		// convert stop id to global id
		id := stop.MainMastExtID
		id, err = slidentifiers.ConvertIDToHafas(id)
		if err != nil {
			continue
		}
		id, err = slidentifiers.ConvertHafasToEFA(id, EFAPrefix)
		if err != nil {
			continue
		}
		stop.MainMastExtID = id
		filtered = append(filtered, stop)
	}
	nearbyResp.Data = filtered
	return nearbyResp, nil
}

type StopsNearbyRequest struct {
	OriginCoordLat  string
	OriginCoordLong string
	MaxNo           string
	Radius          string
	Type            string
}

func (r StopsNearbyRequest) params() url.Values {
	params := url.Values{}
	if r.OriginCoordLat != "" {
		params.Set("originCoordLat", r.OriginCoordLat)
	}
	if r.OriginCoordLong != "" {
		params.Set("originCoordLong", r.OriginCoordLong)
	}
	if r.MaxNo != "" {
		params.Set("maxNo", r.MaxNo)
	}
	if r.Radius != "" {
		params.Set("r", r.Radius)
	}
	if r.Type != "" {
		params.Set("type", r.Type)
	}
	return params
}

type LocationList struct {
	ErrorCode string         `xml:"errorCode,attr"`
	Data      []*StopLocation `xml:"StopLocation"`
}

type StopLocation struct {
	Name          string  `xml:"name,attr"`
	ID            string  `xml:"id,attr"`
	ExtID         string  `xml:"extId,attr"`
	MainMastExtID string  `xml:"mainMastExtId,attr"`
	Lat           float64 `xml:"lat,attr"`
	Lon           float64 `xml:"lon,attr"`
	Distance      int     `xml:"dist,attr"`
}
