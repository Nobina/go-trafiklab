package trafiklab

import (
	"net/url"

  "github.com/nobina/go-requester"
)

func (c *Client) Departures(req *DeparturesRequest) (*DepartureResponse, error) {
	req.key = c.apiKeys[keyDepartures]
	departuresResp := &DepartureResponse{}
	resp, err := c.client.Do(
		requester.WithPath("/api2/realtimedeparturesV4.xml"),
		requester.WithQuery(req.params()),
	)
	if err != nil {
		return nil, err
	}

	return departuresResp, resp.XML(departuresResp)
}

type DeparturesRequest struct {
	key string

	SiteID     string `json:"site_id"`
	TimeWindow string `json:"time_window"`
	Bus        bool   `json:"bus"`
	Metro      bool   `json:"metro"`
	Train      bool   `json:"train"`
	Tram       bool   `json:"tram"`
	Ship       bool   `json:"ship"`
}

func (r DeparturesRequest) params() url.Values {
	params := url.Values{}
	if r.key != "" {
		params.Set("key", r.key)
	}
	if r.SiteID != "" {
		params.Set("SiteId", r.SiteID)
	}
	if r.TimeWindow != "" {
		params.Set("TimeWindow", r.TimeWindow)
	}
	if r.Bus {
		params.Set("Bus", "true")
	}
	if r.Metro {
		params.Set("Metro", "true")
	}
	if r.Train {
		params.Set("Train", "true")
	}
	if r.Tram {
		params.Set("Tram", "true")
	}
	if r.Ship {
		params.Set("Ship", "true")
	}
	return params
}

type DepartureResponse struct {
	StatusCode    int32          `json:"status_code"`
	Message       string         `json:"message"`
	ExecutionTime int            `json:"execution_time"`
	Data          *DepartureData `xml:"ResponseData" json:"data"`
}

type DepartureData struct {
	LatestUpdate        string               `json:"latestUpdate"`
	DataAge             int32                `json:"dataAge"`
	Buses               []Departure          `json:"buses" xml:"Buses>Bus"`
	Metros              []Departure          `json:"metros" xml:"Metros>Metro"`
	Trains              []Departure          `json:"trains" xml:"Trains>Train"`
	Trams               []Departure          `json:"trams" xml:"Trams>Tram"`
	Ships               []Departure          `json:"ships" xml:"Ships>Ship"`
	StopPointDeviations []StopPointDeviation `json:"stopPointDeviations" xml:"StopPointDeviations>StopPointDeviation"`
}

type Departure struct {
	TransportMode        string               `json:"transportMode"`
	LineNumber           string               `json:"lineNumber"`
	Destination          string               `json:"destination"`
	JourneyDirection     int32                `json:"journeyDirection"`
	GroupOfLine          string               `json:"groupOfLine"`
	StopAreaName         string               `json:"stopAreaName"`
	StopAreaNumber       int32                `json:"stopAreaNumber"`
	StopPointNumber      int32                `json:"stopPointNumber"`
	StopPointDesignation string               `json:"stopPointDesignation"`
	TimeTabledDateTime   string               `json:"timeTabledDateTime"`
	ExpectedDateTime     string               `json:"expectedDateTime"`
	DisplayTime          string               `json:"displayTime"`
	JourneyNumber        int32                `json:"journeyNumber"`
	Deviations           []DepartureDeviation `json:"deviations" xml:"Deviations>Deviation"`
}

type StopPointDeviation struct {
	StopInfo struct {
		StopAreaNumber int32
		StopAreaName   string
		TransportMode  string
		GroupOfLine    string
	} `json:"stopInfo"`
	Deviation DepartureDeviation `json:"deviation"`
}

type DepartureDeviation struct {
	Consequence     string `json:"consequence"`
	ImportanceLevel int32  `json:"importanceLevel"`
	Text            string `json:"text"`
}
