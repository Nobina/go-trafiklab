package trafiklab

import (
	"io"
	"net/url"
)

var (
	trafficStatusAPI = &apiConfig{
		host:   "http://api.sl.se",
		path:   "/api2/trafficsituation.xml",
		method: "GET",
	}
)

type EventIcon string

var (
	EventIconPlanned = EventIcon("EventPlanned")
	EventIconGood    = EventIcon("EventGood")
	EventIconMinor   = EventIcon("EventMinor")
	EventIconMajor   = EventIcon("EventMajor")
)

func (c *Client) TrafficStatus() (*TrafficStatusResponse, error) {
	trafficStatusResp := &TrafficStatusResponse{}
	if _, err := c.doXML(trafficStatusAPI, trafficStatusRequest{
		key: c.apiKeys[keyTrafficStatus],
	}, trafficStatusResp); err != nil {
		return nil, err
	}

	return trafficStatusResp, nil
}

type trafficStatusRequest struct {
	key string
}

func (r trafficStatusRequest) body() (io.Reader, error) { return nil, nil }
func (r trafficStatusRequest) params() url.Values {
	params := url.Values{}
	if r.key != "" {
		params.Set("key", r.key)
	}
	return params
}

type TrafficStatusResponse struct {
	StatusCode    int32    `json:"status_code"`
	Message       string   `json:"message"`
	ExecutionTime int      `json:"execution_time"`
	Data          []Status `xml:"ResponseData>TrafficTypes>TrafficType" json:"data"`
}

type Status struct {
	Name            string  `xml:"Name,attr" json:"name"`
	Type            string  `xml:"Type,attr" json:"type"`
	StatusIcon      string  `xml:"StatusIcon,attr" json:"status_icon"`
	Expanded        bool    `xml:"Expanded,attr" json:"expanded"`
	HasPlannedEvent bool    `xml:"HasPlannedEvent,attr" json:"has_planned_event"`
	Events          []Event `xml:"Events>TrafficEvent" json:"events"`
}

type Event struct {
	EventID      int       `xml:"EventId" json:"event_id"`
	Message      string    `json:"message"`
	LineNumbers  string    `json:"line_numbers"`
	Expanded     bool      `json:"expanded"`
	Planned      bool      `json:"planned"`
	SortIndex    int       `json:"sort_index"`
	TrafficLine  string    `json:"traffic_line"`
	EventInfoURL string    `json:"event_info_url"`
	StatusIcon   EventIcon `json:"status_icon"`
}
