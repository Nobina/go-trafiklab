package travelplanner

import (
	"net/url"
)

func (c *Client) Reconstruction(ctx string) (*TripResp, error) {
	path := "/api2/travelplannerv3_1/reconstruction.xml?ctx=" + url.QueryEscape(ctx)
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	tripResp := &TripResp{}
	_, err = c.Do(req, tripResp)
	if err != nil {
		return nil, err
	}

	return tripResp, nil
}
