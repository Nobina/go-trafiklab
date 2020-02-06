package trafiklab

import (
	"net/url"

	"github.com/nobina/go-requester"
)

type Stops struct {
	common *Client
}

func (c *Stops) Query(req *StopsQueryRequest) (*TypeaheadResponse, error) {
	req.key = c.common.apiKeys[keyStopsQuery]
	queryResp := &TypeaheadResponse{}
	resp, err := c.common.client.Do(
		requester.WithPath("/api2/typeahead.xml"),
		requester.WithQuery(req.params()),
	)
	if err != nil {
		return nil, err
	}

	return queryResp, resp.XML(queryResp)
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

func (c *Stops) Nearby(req *StopsNearbyRequest) (*LocationList, error) {
	req.key = c.common.apiKeys[keyStopsNearby]
	nearbyResp := &LocationList{}
	resp, err := c.common.client.Do(
		requester.WithPath("/api2/nearbystopsv2.xml"),
		requester.WithQuery(req.params()),
	)
	if err != nil {
		return nil, err
	}

	return nearbyResp, resp.XML(nearbyResp)
}

type StopsNearbyRequest struct {
	key string

	OriginCoordLat  string
	OriginCoordLong string
	MaxNo           string
	Radius          string
	Type            string
}

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
