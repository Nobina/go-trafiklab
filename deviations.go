package trafiklab

import (
	"net/url"
	"time"

	"github.com/nobina/go-requester"
)

func (c *Client) Deviations(req *DeviationsRequest) (*DeviationsResponse, error) {
	req.key = c.apiKeys[keyDeviations]
	if req.key == "" {
		return nil, ErrMissingAPIKey
	}
	deviationsResp := &DeviationsResponse{}
	resp, err := c.client.Do(
		requester.WithPath("/api2/deviationsrawdata.json"),
		requester.WithQuery(req.params()),
	)
	if err != nil {
		return nil, err
	}

	return deviationsResp, resp.JSON(deviationsResp)
}

type DeviationsRequest struct {
	key string

	TransportMode string `json:"transport_mode"`
	LineNumber    string `json:"line_number"`
	SiteID        string `json:"site_id"`
	FromDate      string `json:"from_date"`
	ToDate        string `json:"to_date"`
}

func (r DeviationsRequest) params() url.Values {
	params := url.Values{}
	if r.key != "" {
		params.Set("key", r.key)
	}
	if r.TransportMode != "" {
		params.Set("TransportMode", r.TransportMode)
	}
	if r.LineNumber != "" {
		params.Set("LineNumber", r.LineNumber)
	}
	if r.SiteID != "" {
		params.Set("SiteId", r.SiteID)
	}
	if r.FromDate != "" {
		params.Set("FromDate", r.FromDate)
	}
	if r.ToDate != "" {
		params.Set("ToDate", r.ToDate)
	}
	return params
}

type DeviationsResponse struct {
	StatusCode    int32       `json:"status_code"`
	Message       string      `json:"message"`
	ExecutionTime int         `json:"execution_time"`
	Data          []Deviation `json:"ResponseData"`
}

type Deviation struct {
	Priority                int    `json:"Priority"`
	SiteID                  string `json:"SiteId"`
	LineNumber              string `json:"LineNumber"`
	TransportMode           string `json:"TransportMode"`
	Created                 string `json:"Created"`
	MainNews                bool   `json:"MainNews"`
	SortOrder               int    `json:"SortOrder"`
	Header                  string `json:"Header"`
	Details                 string `json:"Details"`
	Scope                   string `json:"Scope"`
	DevCaseGid              int64  `json:DevCaseGid`
	DevMessageVersionNumber int    `json:"DevMessageVersionNumber"`
	ScopeElements           string `json:"ScopeElements"`
	FromDateTime            string `json:"FromDateTime"`
	UpToDateTime            string `json:"UpToDateTime"`
	Updated                 string `json:"Updated"`
}

func (d *Deviation) FromDate() (time.Time, error) {
	return time.ParseInLocation(time.RFC3339Nano, d.FromDateTime, LocationEuropeStockholm)
}

func (d *Deviation) ToDate() (time.Time, error) {
	return time.ParseInLocation(time.RFC3339Nano, d.UpToDateTime, LocationEuropeStockholm)
}

func (d *Deviation) UpdatedDate() (time.Time, error) {
	return time.ParseInLocation(time.RFC3339Nano, d.Updated, LocationEuropeStockholm)
}
