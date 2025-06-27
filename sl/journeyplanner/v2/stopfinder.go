package journeyplanner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	stopFinderTypeCoord = "coord"
	stopFinderTypeAny = "any"
)

// type stopFinderRequest struct {
// 	//search string
// 	// If coord, latitude and longitude in lat,lng format
// 	// SL wants a more complex format:
// 	// <x>:<y>:WGS84[dd.ddddd] E.g. "18.013809:59.335104:WGS84[dd.ddddd]"
// 	// We convert a standard lat,lng format to this format
// 	Name string
// 	// can be:
// 	// "none","stop", "poi", "suburb", "street", "address", "unknown"
// 	// must be provided at least one
// 	Filter []string
// 	// Can be: "coord", "any"
// 	// If coord, (TODO(pb): Haven't been able to search coord yet)
// 	SearchType string
// }

type TrafiklabStopRequester interface {
	toTrafiklabRequest() (*stopFinderTrafiklabRequest, error)
}

type StopFinderSearchRequest struct {
	// Search query string
	Name string
	// Filter types to include in the search. Can be:
	// "none","stop", "poi", "suburb", "street", "address", "unknown"
	// At least one must be provided
	Filter []string
}

func NewStopFinderSearchRequest(name string, filter []string) *StopFinderSearchRequest {
	return &StopFinderSearchRequest{
		Name: name,
		Filter: filter,
	}
}

func (sfr *StopFinderSearchRequest) toTrafiklabRequest() (*stopFinderTrafiklabRequest, error) {
	filter, err := StopFilterFromString(sfr.Filter...)
	if err != nil {
		return nil, err
	}
	return &stopFinderTrafiklabRequest{
		Name: sfr.Name,
		Type: stopFinderTypeAny,
		AnyObjFilter: filter,
	}, nil
}

type StopFinderPosRequest struct {
	Position LatLng
	// Filter types to include in the search. Can be:
	// "none","stop", "poi", "suburb", "street", "address", "unknown"
	// At least one must be provided
	Filter []string
}

func NewStopFinderPosRequest(position LatLng, filter []string) *StopFinderPosRequest {
	return &StopFinderPosRequest{
		Position: position,
		Filter: filter,
	}
}

func (sfr *StopFinderPosRequest) toTrafiklabRequest() (*stopFinderTrafiklabRequest, error) {
	filter, err := StopFilterFromString(sfr.Filter...)
	if err != nil {
		return nil, err
	}
	return &stopFinderTrafiklabRequest{
		Name: fmt.Sprintf("%f:%f:WGS84[dd.ddddd]", sfr.Position.Longitude, sfr.Position.Latitude),
		Type: stopFinderTypeCoord,
		AnyObjFilter: filter,
	}, nil
}


type LatLng struct {
	Latitude float64
	Longitude float64
}

func (ll *LatLng) FromString(s string) error {
	// Check if string is in Trafiklab format: "lng:lat:WGS84[dd.ddddd]"
	if strings.Contains(s, "WGS84[dd.ddddd]") {
		// Remove the WGS84[dd.ddddd] suffix
		coordPart := strings.TrimSuffix(s, ":WGS84[dd.ddddd]")
		parts := strings.Split(coordPart, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid Trafiklab format: %s", s)
		}

		var err error
		// In Trafiklab format, longitude comes first, then latitude
		ll.Longitude, err = strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return fmt.Errorf("invalid longitude in Trafiklab format: %w", err)
		}
		ll.Latitude, err = strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return fmt.Errorf("invalid latitude in Trafiklab format: %w", err)
		}
		return nil
	}

	// Fall back to comma-separated format: "lat,lng"
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return fmt.Errorf("invalid lat,lng format: %s", s)
	}

	var err error
	ll.Latitude, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return fmt.Errorf("invalid latitude: %w", err)
	}
	ll.Longitude, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return fmt.Errorf("invalid longitude: %w", err)
	}
	return nil
}

func (ll *LatLng) ToTrafiklabString() string {
	return fmt.Sprintf("%f:%f:WGS84[dd.ddddd]", ll.Longitude, ll.Latitude)
}

type stopFinderTrafiklabRequest struct {
	// name_sf
	Name string
	// type_sf
	Type string
	// any_obj_filter_sf
	AnyObjFilter StopFilter
}

func (sfr *stopFinderTrafiklabRequest) Valid() error {
	if sfr.Name == "" {
		return errors.New("name is required")
	}
	if sfr.Type != stopFinderTypeCoord && sfr.Type != stopFinderTypeAny {
		return errors.New("type must be either coord or any")
	}
	if sfr.AnyObjFilter < 0 {
		return errors.New("any_obj_filter is required")
	}
	return nil
}

func (sfr *stopFinderTrafiklabRequest) toParams() url.Values {
	params := url.Values{}
	params.Set("name_sf", sfr.Name)
	params.Set("type_sf", sfr.Type)
	params.Set("any_obj_filter_sf", strconv.Itoa(int(sfr.AnyObjFilter)))
	return params
}

type StopFinderResponse struct {
	// System messages from backend
	SystemMessages []SystemMessage `json:"systemMessages"`
	// Stop locations found
	StopLocations []StopLocation `json:"locations"`
}

type StopLocation struct {
	// ID of the stop
	ID string `json:"id"`
	// Whether this is a global ID
	IsGlobalID bool `json:"isGlobalId"`
	// Name of the stop including municipality
	Name string `json:"name"`
	// Name of the stop without municipality
	DisassembledName string `json:"disassembledName"`
	// Coordinates of the stop
	Coordinates Coordinates `json:"coord"`
	// Street name if type is street or singlehouse
	StreetName string `json:"streetName,omitempty"`
	// Building number if type is singlehouse
	BuildingNumber string `json:"buildingNumber,omitempty"`
	// Type of the result (address, stop, singlehouse, poi, street)
	Type string `json:"type"`
	// Quality of the query matching
	MatchQuality int `json:"matchQuality"`
	// Whether this is the best search match
	IsBest bool `json:"isBest"`
	// Products at this stop (0=train, 2=metro, 4=train/tram, 5=bus, 9=ship/ferry, 10=transit on demand)
	ProductClasses []int `json:"productClasses,omitempty"`
	// Parent location information
	Parent *ParentLocation `json:"parent,omitempty"`
}

type ParentLocation struct {
	// ID of the principality
	ID string `json:"id"`
	// Name of the municipality
	Name string `json:"name"`
	// Type of the principality
	Type string `json:"type"`
}

type Coordinates struct {
	// Latitude of the stop
	Latitude float64 `json:"lat"`
	// Longitude of the stop
	Longitude float64 `json:"lon"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *Coordinates) UnmarshalJSON(data []byte) error {
	var arr []float64
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) != 2 {
		return fmt.Errorf("expected 2 elements for coordinates, got %d", len(arr))
	}
	c.Latitude = arr[0]
	c.Longitude = arr[1]
	return nil
}


func (c *Client) StopFinder(ctx context.Context, sfr TrafiklabStopRequester) (*StopFinderResponse, error) {
	tfr, err := sfr.toTrafiklabRequest()
	if err != nil {
		return nil, err
	}
	if err := tfr.Valid(); err != nil {
		return nil, err
	}

	url := c.baseURL + "/stop-finder"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = tfr.toParams().Encode()
	c.addDefaultHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var res StopFinderResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) addDefaultHeaders(req *http.Request) {
	req.Header.Set("X-Correlation-ID", c.clientID)
}
