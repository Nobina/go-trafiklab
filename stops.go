package trafiklab

import (
	"io"
	"net/url"
)

var (
	stopsQueryAPI = &apiConfig{
		host:   "http://api.sl.se",
		path:   "/api2/typeahead.xml",
		method: "GET",
	}
	stopsNearbyAPI = &apiConfig{
		host:   "http://api.sl.se",
		path:   "/api2/nearbystopsv2.xml",
		method: "GET",
	}
)

type Stops struct {
	common *Client
}

func (c *Stops) Query(queryReq *StopsQueryRequest) (*TypeaheadResponse, error) {
	queryReq.key = c.common.apiKeys[keyStopsQuery]
	queryResp := &TypeaheadResponse{}
	if _, err := c.common.doXML(stopsQueryAPI, queryReq, queryResp); err != nil {
		return nil, err
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

func (r StopsQueryRequest) body() (io.Reader, error) { return nil, nil }
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

func (c *Stops) Nearby(nearbyReq *StopsNearbyRequest) (*LocationList, error) {
	nearbyReq.key = c.common.apiKeys[keyStopsNearby]
	nearbyResp := &LocationList{}
	if _, err := c.common.doXML(stopsNearbyAPI, nearbyReq, nearbyResp); err != nil {
		return nil, err
	}

	return nearbyResp, nil
}

type StopsNearbyRequest struct {
	key string

	OriginCoordLat  string
	OriginCoordLong string
	MaxNo           string
	Radius          string
	Type            string
}

func (r StopsNearbyRequest) body() (io.Reader, error) { return nil, nil }
func (r StopsNearbyRequest) params() url.Values {
	params := url.Values{}
	if r.key != "" {
		params.Set("key", r.key)
	}
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
	Data      []StopLocation `xml:"StopLocation"`
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
