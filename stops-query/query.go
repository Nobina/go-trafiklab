package stopsquery

import (
	"net/url"
)

func (c *Client) StopsQuery(nearbyReq *StopsQueryRequest) (*TypeaheadResponse, error) {
	path := "/api2/typeahead.xml?" + nearbyReq.params().Encode()
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	queryResp := &TypeaheadResponse{}
	_, err = c.Do(req, queryResp)
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

type StopsQueryRequest struct {
	SearchString string
	StationsOnly bool
	MaxResults   string
	Type         string
}

func (r StopsQueryRequest) params() url.Values {
	params := url.Values{}
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
