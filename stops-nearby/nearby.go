package stopsnearby

import (
	"net/url"
)

func (c *Client) StopsNearby(nearbyReq *StopsNearbyRequest) (*LocationList, error) {
	path := "/api2/nearbystopsv2.xml?" + nearbyReq.params().Encode()
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	nearbyResp := &LocationList{}
	_, err = c.Do(req, nearbyResp)
	if err != nil {
		return nil, err
	}

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
	Data      []StopLocation `xml:"StopLocation"`
}

type StopLocation struct {
	IDX      int     `xml:"idx,attr"`
	Name     string  `xml:"name,attr"`
	ID       string  `xml:"id,attr"`
	Lat      float64 `xml:"lat,attr"`
	Lon      float64 `xml:"lon,attr"`
	Distance int     `xml:"dist,attr"`
}
