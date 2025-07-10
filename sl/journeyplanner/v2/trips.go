package journeyplanner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/nobina/go-trafiklab/slidentifiers"
)

const (
	SearchTypeCoord = "coord"
	SearchTypeAny = "any"

	RouteTypeLeastInterchange = "leastinterchange"
	RouteTypeLeastTime = "leasttime"
	RouteTypeLeastWalking = "leastwalking"

	MotFlagCommuterTrain = "commuter_train"
	MotFlagMetro = "metro"
	MotFlagTramsTrains = "trams_trains"
	MotFlagBus = "bus"
	MotFlagShipFerry = "ship_ferry"
	MotFlagTransitOnDemand = "transit_on_demand"
	MotFlagNationalTrain = "national_train"
	MotFlagAccessibleBus = "accessible_bus"
)

type TripsRequest struct {
	// Date and time of the trip
	At time.Time
	// Number of trips to return (required, min: 0, max: 3)
	NumTrips int
	// List of flags to enable. Available flags:
	// - compute_monomodal_trip_pedestrian: Enables the calculation of an additional walk only trip
	// - compute_monomodal_trip_bicycle: Enables the calculation of an additional bike only trip
	// - calc_one_direction: If enabled, prevents calculating one trip before the requested departure time
	// - must_excl: Enables definitive exclusion of Operators and Lines
	// - sel_op: Enables operators filter
	// - no_alt: Supress alternative/additional trips
	// - prefer_incl: Enables preferred inclusion of Operators and Lines
	// - use_prox_foot_search: Enables walk to alternative stops, nearby the selected stop

	// - use_only: Enables definitive inclusion of Operators
	// - prefer_excl: Enables preferred exclusion of Operators and Lines
	// - gen_c: Enables output of coordinate sequences for trip legs
	Flags []string

	// - mot_0: Include commuter train (pendeltåg) in trip calculation
	// - mot_2: Include metro (tunnelbana) in trip calculation
	// - mot_4: Include local train/tram (lokaltåg/spårväg) in trip calculation
	// - mot_5: Include bus (buss) in trip calculation
	// - mot_9: Include ship and ferry (båttrafik) in trip calculation
	// - mot_10: Include transit on demand area service (anropsstyrd områdestrafik) in trip calculation
	// - mot_14: Include national train (fjärrtåg) in trip calculation
	// - mot_19: Include accessible bus (närtrafik) in trip calculation
	// Deprecated: Use IncludeMotFlags instead
	AvoidMotFlags []string

	// IncludeMotFlags is a list of mot flags to include in the trip calculation
	// - mot_0: Include commuter train (pendeltåg) in trip calculation
	// - mot_2: Include metro (tunnelbana) in trip calculation
	// - mot_4: Include local train/tram (lokaltåg/spårväg) in trip calculation
	// - mot_5: Include bus (buss) in trip calculation
	// - mot_9: Include ship and ferry (båttrafik) in trip calculation
	// - mot_10: Include transit on demand area service (anropsstyrd områdestrafik) in trip calculation
	// - mot_14: Include national train (fjärrtåg) in trip calculation
	// - mot_19: Include accessible bus (närtrafik) in trip calculation
	IncludeMotFlags []string

	// Language of the session (sv or en)
	Language string
	// Type of the origin (coord or any) (required)
	TypeOrigin string

	// NameOrigin lat,lng if TypeOrigin is coord, else search string
	NameOrigin string
	// Type of the destination (coord or any) (required)
	TypeDestination string
	// NameDestination lat,lng if TypeDestination is coord, else search string
	NameDestination string
	// Type of via point (always any)
	TypeVia string
	// ID of a stop that the trip should go via
	ViaID string
	// Type of not-via point (always any)
	TypeNotVia string
	// ID of a stop that the trip should not go via
	NotViaID string
	// Maximum number of interchanges (min: 0, max: 9)
	MaxChanges int
	// Maximum time in minutes for walk only trip (min: 0, max: 120)
	MaxTimePedestrian int
	// Maximum time in minutes for bike only trip (min: 0, max: 120)
	MaxTimeBicycle int
	// Maximum distance in meters for footpath sections (min: 0, max: 1000)
	MaxLengthPedestrian int
	// Minimum distance in meters for footpath sections (min: 0, max: 1000)
	MinLengthPedestrian int
	// Maximum distance in meters for bike sections (min: 0, max: 1000)
	MaxLengthBicycle int
	// Minimum distance in meters for bike sections (min: 0, max: 1000)
	MinLengthBicycle int
	// Sets the speed in percentage for interchanges and walk/bike paths (min: 25, max: 400)
	ChangeSpeed int
	// Calculate trips with least interchanges, fastest connections or shortest footpaths
	RouteType string
	// Search for departure time or arrival time (dep or arr)
	TripDateTimeDepArr string
	// Extra waiting time at a via stop (HHMM)
	DwellTime string
	// Lines to be excluded
	MustExclLines []string
	// Lines to be preferred excluded
	PreferExclLines []string
	// Lines to be preferred included
	PreferInclLines []string
	// Operators to be used exclusively
	UseOnlyOperators []string
	// Operators to be excluded
	MustExclOperators []string
	// Operators to be preferred excluded
	PreferExclOperators []string
	// Operators to be preferred included
	PreferInclOperators []string
}

var validFlags = map[string]struct{}{
	"compute_monomodal_trip_pedestrian": {},
	"compute_monomodal_trip_bicycle":    {},
	"calc_one_direction":                {},
	"must_excl":                         {},
	"sel_op":                           {},
	"sel_line":                         {},
	"no_alt":                           {},
	"prefer_incl":                      {},
	"use_prox_foot_search":             {},
	"use_only":                         {},
	"prefer_excl":                      {},
	"gen_c":                            {},
}

var avoidSaneFlagsToMotFlags = map[string]string{
	MotFlagCommuterTrain:	"incl_mot_0",
	MotFlagMetro: 		  	"incl_mot_2",
	MotFlagTramsTrains:		"incl_mot_4",
	MotFlagBus:				"incl_mot_5",
	MotFlagShipFerry: 		"incl_mot_9",
	MotFlagTransitOnDemand:	"incl_mot_10",
	MotFlagNationalTrain:	"incl_mot_14",
	MotFlagAccessibleBus:	"incl_mot_19",
}

func (tr *TripsRequest) Valid() error {
	// Check required fields
	if tr.At.IsZero() {
		return errors.New("date and time is required")
	}
	if tr.NumTrips < 1 || tr.NumTrips > 3 {
		return errors.New("number of trips must be between 0 and 3")
	}
	if tr.NameOrigin == ""{
		return errors.New("name_origin is required")
	}
	if tr.TypeOrigin == "" {
		return errors.New("type_origin is required")
	} else if tr.TypeOrigin == "coord" {
		latLng := &LatLng{}
		err := latLng.FromString(tr.NameOrigin)
		if err != nil {
			return fmt.Errorf("invalid lat,lng format: %w", err)
		}
		tr.NameOrigin = latLng.ToTrafiklabString()
	}
	if tr.TypeDestination == "" {
		return errors.New("type_destination is required")
	} else if tr.TypeDestination == "coord" {
		latLng := &LatLng{}
		err := latLng.FromString(tr.NameDestination)
		if err != nil {
			return fmt.Errorf("invalid lat,lng format: %w", err)
		}
		tr.NameDestination = latLng.ToTrafiklabString()
	}

	// Validate language
	if tr.Language != "" && tr.Language != "sv" && tr.Language != "en" {
		return errors.New("language must be either 'sv' or 'en'")
	}

	// Validate type_origin
	if tr.TypeOrigin != "coord" && tr.TypeOrigin != "any" {
		return errors.New("type_origin must be either 'coord' or 'any'")
	}

	// Validate type_destination
	if tr.TypeDestination != "coord" && tr.TypeDestination != "any" {
		return errors.New("type_destination must be either 'coord' or 'any'")
	}

	// Validate type_via
	if tr.TypeVia != "" && tr.TypeVia != "any" {
		return errors.New("type_via must be 'any'")
	}

	// Validate type_not_via
	if tr.TypeNotVia != "" && tr.TypeNotVia != "any" {
		return errors.New("type_not_via must be 'any'")
	}

	// Validate max_changes
	if tr.MaxChanges < 0 || tr.MaxChanges > 9 {
		return errors.New("max_changes must be between 0 and 9")
	}

	// Validate max_time_pedestrian
	if tr.MaxTimePedestrian < 0 || tr.MaxTimePedestrian > 120 {
		return errors.New("max_time_pedestrian must be between 0 and 120")
	}

	// Validate max_time_bicycle
	if tr.MaxTimeBicycle < 0 || tr.MaxTimeBicycle > 120 {
		return errors.New("max_time_bicycle must be between 0 and 120")
	}

	// Validate max_length_pedestrian
	if tr.MaxLengthPedestrian < 0 || tr.MaxLengthPedestrian > 1000 {
		return errors.New("max_length_pedestrian must be between 0 and 1000")
	}

	// Validate min_length_pedestrian
	if tr.MinLengthPedestrian < 0 || tr.MinLengthPedestrian > 1000 {
		return errors.New("min_length_pedestrian must be between 0 and 1000")
	}

	// Validate max_length_bicycle
	if tr.MaxLengthBicycle < 0 || tr.MaxLengthBicycle > 1000 {
		return errors.New("max_length_bicycle must be between 0 and 1000")
	}

	// Validate min_length_bicycle
	if tr.MinLengthBicycle < 0 || tr.MinLengthBicycle > 1000 {
		return errors.New("min_length_bicycle must be between 0 and 1000")
	}

	// Validate route_type
	validRouteTypes := []string{"leastinterchange", "leasttime", "leastwalking"}
	if tr.RouteType != "" && !slices.Contains(validRouteTypes, tr.RouteType) {
		return errors.New("route_type must be one of: " + strings.Join(validRouteTypes, ", "))
	}

	// Validate trip_date_time_dep_arr
	if tr.TripDateTimeDepArr != "" && tr.TripDateTimeDepArr != "dep" && tr.TripDateTimeDepArr != "arr" {
		return errors.New("trip_date_time_dep_arr must be either 'dep' or 'arr'")
	}

	// Validate flags
	for _, flag := range tr.Flags {
		if _, ok := validFlags[flag]; !ok {
			return fmt.Errorf("invalid flag: %s", flag)
		}
	}
	if len(tr.AvoidMotFlags) > 0 && len(tr.IncludeMotFlags) > 0 {
		return errors.New("avoid_mot_flags and include_mot_flags cannot be used together")
	}
	// Validate avoid_mot_flags
	if len(tr.AvoidMotFlags) > 0 {
		for _, flag := range tr.AvoidMotFlags {
			if _, ok := avoidSaneFlagsToMotFlags[flag]; !ok {
				return fmt.Errorf("invalid avoid_mot_flag: %s", flag)
			}
		}
	}

	// Validate include_mot_flags
	if len(tr.IncludeMotFlags) > 0 {
		for _, flag := range tr.IncludeMotFlags {
			if _, ok := avoidSaneFlagsToMotFlags[flag]; !ok {
				return fmt.Errorf("invalid include_mot_flag: %s", flag)
			}
		}
	}

	// If no mot flags are set, set all flags to true
	if len(tr.IncludeMotFlags) == 0 && len(tr.AvoidMotFlags) == 0 {
		tr.IncludeMotFlags = []string{
			MotFlagCommuterTrain,
			MotFlagNationalTrain,
			MotFlagTramsTrains,
			MotFlagMetro,
			MotFlagBus,
			MotFlagShipFerry,
			MotFlagAccessibleBus,
			MotFlagTransitOnDemand,
		}
	}

	return nil
}

func (tr *TripsRequest) toParams() url.Values {
	params := url.Values{}

	params.Set("itd_date", tr.At.Format("20060102")) // YYYYMMDD format
	params.Set("itd_time", tr.At.Format("1504"))     // HHMM format

	params.Set("name_origin", tr.NameOrigin)
	params.Set("name_destination", tr.NameDestination)
	params.Set("type_origin", tr.TypeOrigin)
	params.Set("type_destination", tr.TypeDestination)
	params.Set("calc_number_of_trips", strconv.Itoa(tr.NumTrips))


	// Set all flags to 1
	for _, flag := range tr.Flags {
		params.Set(flag, "true")
	}

	// Handle mot flags - either avoid or include logic
	if len(tr.AvoidMotFlags) > 0 {
		// Set all avoid flags to false
		for _, flag := range tr.AvoidMotFlags {
			params.Set(avoidSaneFlagsToMotFlags[flag], "false")
		}
		// Set all flags not in the avoid list to true
		for flag, motFlag := range avoidSaneFlagsToMotFlags {
			if !slices.Contains(tr.AvoidMotFlags, flag) {
				params.Set(motFlag, "true")
			}
		}
	} else if len(tr.IncludeMotFlags) > 0 {
		// Set all flags to false first
		for _, motFlag := range avoidSaneFlagsToMotFlags {
			params.Set(motFlag, "false")
		}
		// Then set only the included flags to true
		for _, flag := range tr.IncludeMotFlags {
			params.Set(avoidSaneFlagsToMotFlags[flag], "true")
		}
	}

	// Set optional parameters if they are not empty
	if tr.Language != "" {
		params.Set("language", tr.Language)
	}

	if tr.NameOrigin != "" {
		params.Set("name_origin", tr.NameOrigin)
	}
	if tr.NameDestination != "" {
		params.Set("name_destination", tr.NameDestination)
	}
	if tr.TypeVia != "" {
		params.Set("type_via", tr.TypeVia)
	}
	if tr.ViaID != "" {
		params.Set("name_via", tr.ViaID)
	}
	if tr.TypeNotVia != "" {
		params.Set("type_not_via", tr.TypeNotVia)
	}
	if tr.NotViaID != "" {
		params.Set("name_not_via", tr.NotViaID)
	}
	if tr.MaxChanges > 0 {
		params.Set("max_changes", strconv.Itoa(tr.MaxChanges))
	}
	if tr.MaxTimePedestrian > 0 {
		params.Set("max_time_pedestrian", strconv.Itoa(tr.MaxTimePedestrian))
	}
	if tr.MaxTimeBicycle > 0 {
		params.Set("max_time_bicycle", strconv.Itoa(tr.MaxTimeBicycle))
	}
	if tr.MaxLengthPedestrian > 0 {
		params.Set("max_length_pedestrian", strconv.Itoa(tr.MaxLengthPedestrian))
	}
	if tr.MinLengthPedestrian > 0 {
		params.Set("min_length_pedestrian", strconv.Itoa(tr.MinLengthPedestrian))
	}
	if tr.MaxLengthBicycle > 0 {
		params.Set("max_length_bicycle", strconv.Itoa(tr.MaxLengthBicycle))
	}
	if tr.MinLengthBicycle > 0 {
		params.Set("min_length_bicycle", strconv.Itoa(tr.MinLengthBicycle))
	}
	if tr.ChangeSpeed > 0 {
		params.Set("change_speed", strconv.Itoa(tr.ChangeSpeed))
	}
	if tr.RouteType != "" {
		params.Set("route_type", tr.RouteType)
	}
	if tr.TripDateTimeDepArr != "" {
		params.Set("itd_trip_date_time_dep_arr", tr.TripDateTimeDepArr)
	}
	if tr.DwellTime != "" {
		params.Set("dwell_time", tr.DwellTime)
	}

	// Add line and operator filters
	for _, line := range tr.MustExclLines {
		params.Add("must_excl_line", line)
	}
	for _, line := range tr.PreferExclLines {
		params.Add("prefer_excl_line", line)
	}
	for _, line := range tr.PreferInclLines {
		params.Add("prefer_incl_line", line)
	}
	for _, op := range tr.UseOnlyOperators {
		params.Add("use_only_op", op)
	}
	for _, op := range tr.MustExclOperators {
		params.Add("must_excl_op", op)
	}
	for _, op := range tr.PreferExclOperators {
		params.Add("prefer_excl_op", op)
	}
	for _, op := range tr.PreferInclOperators {
		params.Add("prefer_incl_op", op)
	}

	return params
}

// Convenience function to get the raw query for the trips request
func (tr *TripsRequest) RawQuery() string {
	return tr.toParams().Encode()
}

type TripsResponse struct {
	// System messages from backend
	SystemMessages []SystemMessage `json:"systemmessages"`
	// List of journeys found
	Journeys []Journey `json:"journeys"`
}

type Journey struct {
	// ID of the journey
	ID string `json:"tripId"`
	// Duration of the journey in seconds
	Duration int `json:"tripDuration"`
	// Real-time duration of the journey in seconds
	RealTimeDuration int `json:"tripRtDuration"`
	// Rating of the journey
	Rating int `json:"rating"`
	// Whether this is an additional journey
	IsAdditional bool `json:"isAdditional"`
	// Number of interchanges
	Interchanges int `json:"interchanges"`
	// Whether real-time is only informative
	IsRealTimeOnlyInformative bool `json:"isRealtimeOnlyInformative"`
	// Real-time explanation index
	RealTimeExplanationIdx string `json:"realtimeExplanationIdx"`
	// List of legs in the journey
	Legs []JourneyLeg `json:"legs"`
	// Days of service
	DaysOfService JourneyDaysOfService `json:"daysOfService"`
	// If trip is cancelled
	TripImpossible bool `json:"tripImpossible"`
}

type JourneyLeg struct {
	// List of information messages
	Infos []JourneyInfo `json:"infos"`
	// List of hints for the journey
	Hints []LegHint `json:"hints"`
	// Distance of the leg
	Distance int `json:"distance"`
	// Duration of the leg
	Duration int `json:"duration"`
	// Foot path information
	FootPathInfo []FootPathInfo `json:"footPathInfo"`
	// Origin of the leg
	Origin JourneyStop `json:"origin"`
	// Destination of the leg
	Destination JourneyStop `json:"destination"`
	// Transportation information
	Transportation JourneyTransportation `json:"transportation"`
	// Sequence of stops
	StopSequence []JourneyStop `json:"stopSequence"`
	// Coordinates for the leg
	Coordinates [][]float64 `json:"coords"`
	// Properties of the leg
	Properties map[string]any `json:"properties"`
}

type JourneyStop struct {
	// Whether this is a global ID
	IsGlobalID bool `json:"isGlobalId"`
	// ID of the stop
	ID string `json:"id"`
	// Name of the stop with locality
	Name string `json:"name"`
	// Name without locality
	DisassembledName string `json:"disassembledName"`
	// Type of the stop
	Type string `json:"type"`
	// Coordinates of the stop
	Coordinates []float64 `json:"coord"`
	// Level of the stop
	Level int `json:"niveau"`
	// Parent stop information
	Parent *JourneyStop `json:"parent"`
	// Product classes available at the stop
	ProductClasses []int `json:"productClasses"`
	// Departure time according to base timetable
	DepartureTimeBaseTimetable time.Time `json:"departureTimeBaseTimetable"`
	// Planned departure time
	DepartureTimePlanned time.Time `json:"departureTimePlanned"`
	// Estimated departure time
	DepartureTimeEstimated time.Time `json:"departureTimeEstimated"`
	// Arrival time according to base timetable
	ArrivalTimeBaseTimetable time.Time `json:"arrivalTimeBaseTimetable"`
	// Planned arrival time
	ArrivalTimePlanned time.Time `json:"arrivalTimePlanned"`
	// Estimated arrival time
	ArrivalTimeEstimated time.Time `json:"arrivalTimeEstimated"`
	// Properties of the stop
	Properties map[string]any `json:"properties"`
}

type JourneyTransportation struct {
	// ID of the line
	ID string `json:"id"`
	// Name of the line
	Name string `json:"name"`
	// Number of the line
	Number string `json:"number"`
	// Product information
	Product JourneyProduct `json:"product"`
	// Operator information
	Operator JourneyOperator `json:"operator"`
	// Destination information
	Destination JourneyDestination `json:"destination"`
	// Properties of the transportation
	Properties map[string]any `json:"properties"`
	// Whether this is a Samtrafik service
	IsSamtrafik bool `json:"isSamtrafik"`
	// Disassembled name of the line
	DisassembledName string `json:"disassembledName"`
}

type JourneyProduct struct {
	// ID of the product
	ID int `json:"id"`
	// Class of the product
	Class int `json:"class"`
	// Name of the product
	Name string `json:"name"`
	// Icon ID of the product
	IconID int `json:"iconId"`
}

type JourneyOperator struct {
	// ID of the operator
	ID string `json:"id"`
	// Name of the operator
	Name string `json:"name"`
}

type JourneyDestination struct {
	// ID of the destination
	ID string `json:"id"`
	// Name of the destination
	Name string `json:"name"`
	// Type of the destination
	Type string `json:"type"`
}

type JourneyInfo struct {
	// ID of the info
	ID string `json:"id"`
	// List of info links
	InfoLinks []JourneyInfoLink `json:"infoLinks"`
	// Priority of the info
	Priority string `json:"priority"`
	// Type of the info
	Type string `json:"type"`
	// Version of the info
	Version int `json:"version"`
}

type JourneyInfoLink struct {
	// Properties of the info link
	Properties map[string]any `json:"properties"`
	// Title of the info link
	Title string `json:"title"`
	// URL of the info link
	URL string `json:"url"`
}

type LegHint struct {
	// Provider code
	ProviderCode string `json:"providerCode"`
	// Content of the hint
	Content string `json:"content"`
}

type JourneyDaysOfService struct {
	// RVB value
	RVB string `json:"rvb"`
}

type FootPathInfo struct {
	// Duration of the foot path
	Duration int `json:"duration"`
	// Elements of the foot path
	FootPathElements []FootPathElement `json:"footPathElem"`
	// Position of the foot path
	Position string `json:"position"`
}

type FootPathElement struct {
	// Attributes of the foot path element
	Attributes map[string]any `json:"attributes"`
	// Description of the foot path element
	Description string `json:"description"`
	// Origin of the foot path element
	Origin FootPathStop `json:"origin"`
	// Destination of the foot path element
	Destination FootPathStop `json:"destination"`
	// Level of the foot path element
	Level string `json:"level"`
	// Level from
	LevelFrom int `json:"levelFrom"`
	// Level to
	LevelTo int `json:"levelTo"`
	// Opening hours
	OpeningHours []int `json:"openingHours"`
}

type FootPathStop struct {
	// Coordinates of the stop
	Coordinates []float64 `json:"coord"`
	// ID of the stop
	ID string `json:"id"`
	// Whether this is a global ID
	IsGlobalID bool `json:"isGlobalId"`
	// Name of the stop
	Name string `json:"name"`
	// Properties of the stop
	Properties map[string]any `json:"properties"`
	// Type of the stop
	Type string `json:"type"`
}

func (c *Client) Trips(ctx context.Context, tr *TripsRequest) (*TripsResponse, error) {
	if err := tr.Valid(); err != nil {
		return nil, err
	}
	// convert old ids to new efa ids
	fromID := tr.NameOrigin
	toID := tr.NameDestination
	if tr.TypeOrigin == stopFinderTypeAny && len(fromID) != 16 {
		if slidentifiers.IsSiteID(fromID) {
			var err error
			fromID, err = slidentifiers.ConvertSiteIDToEFA(fromID, slidentifiers.EFAPrefix)
			if err != nil {
				return nil, fmt.Errorf("failed to convert site id to efa id: %w", err)
			}
		}
	}
	if tr.TypeDestination == stopFinderTypeAny && len(toID) != 16 {
		if slidentifiers.IsSiteID(toID) {
			var err error
			toID, err = slidentifiers.ConvertSiteIDToEFA(toID, slidentifiers.EFAPrefix)
			if err != nil {
				return nil, fmt.Errorf("failed to convert site id to efa id: %w", err)
			}
		}
	}
	tr.NameOrigin = fromID
	tr.NameDestination = toID

	url := c.baseURL + "/trips"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = tr.toParams().Encode()
	c.addDefaultHeaders(req)
	if c.isDebug {
		fmt.Printf("url: %s\n", req.URL.String())
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d, request: %s, response: %s", resp.StatusCode, req.URL.RawQuery, string(body))
	}

	var res TripsResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w, body: %s", err, string(body))
	}

	return &res, nil
}