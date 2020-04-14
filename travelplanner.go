package trafiklab

import (
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"
)

var (
	travelplannerTripsAPI = &apiConfig{
		host:   "http://api.sl.se",
		path:   "/api2/travelplannerv3_1/trip.xml",
		method: "GET",
	}

	travelplannerReconstructionAPI = &apiConfig{
		host:   "http://api.sl.se",
		path:   "/api2/travelplannerv3_1/reconstruction.xml",
		method: "GET",
	}

	travelplannerJourneyDetailAPI = &apiConfig{
		host:   "http://api.sl.se",
		path:   "/api2/travelplannerv3_1/journeydetail.xml",
		method: "GET",
	}
)

type Travelplanner struct {
	common *Client
}

type ProductRef int32

const (
	ProductRefTrain   ProductRef = 1
	ProductRefMetro   ProductRef = 2
	ProductRefTram    ProductRef = 4
	ProductRefBus     ProductRef = 8
	ProductRefBoat    ProductRef = 96
	ProductRefCommute ProductRef = 128
)

func (c *Travelplanner) JourneyDetail(journeyDetailReq *JourneyDetailRequest) (*Leg, error) {
	journeyDetailReq.key = c.common.apiKeys[keyTravelplanner]
	legResp := &Leg{}
	if _, err := c.common.doXML(travelplannerJourneyDetailAPI, journeyDetailReq, legResp); err != nil {
		return nil, err
	}

	return legResp, nil
}

type JourneyDetailRequest struct {
	key  string
	ID   string
	Poly bool
}

func (r JourneyDetailRequest) body() (io.Reader, error) { return nil, nil }
func (r JourneyDetailRequest) params() url.Values {
	params := url.Values{}
	if r.key != "" {
		params.Set("key", r.key)
	}
	if r.ID != "" {
		params.Set("id", r.ID)
	}
	if r.Poly {
		params.Set("poly", "1")
	}

	return params
}

func (c *Travelplanner) Reconstruction(ctx string) (*TripResp, error) {
	tripResp := &TripResp{}
	if _, err := c.common.doXML(travelplannerReconstructionAPI, ReconstructionRequest{
		key:     c.common.apiKeys[keyTravelplanner],
		Context: ctx,
	}, tripResp); err != nil {
		return nil, err
	}

	return tripResp, nil
}

type ReconstructionRequest struct {
	key     string
	Context string
}

func (r ReconstructionRequest) body() (io.Reader, error) { return nil, nil }
func (r ReconstructionRequest) params() url.Values {
	params := url.Values{}
	if r.key != "" {
		params.Set("key", r.key)
	}
	if r.Context != "" {
		params.Set("ctx", r.Context)
	}

	return params
}

func (c *Travelplanner) Trips(tripsReq *TripsRequest) (*TripsResp, error) {
	tripsReq.key = c.common.apiKeys[keyTravelplanner]
	tripsResp := &TripsResp{}
	if _, err := c.common.doXML(travelplannerTripsAPI, tripsReq, tripsResp); err != nil {
		return nil, err
	}

	return tripsResp, nil
}

type LegContextualFunc func(leg, prevLeg, prevTransportLeg, nextLeg, nextTransportLeg *Leg, i int) error

type Via struct {
	ViaID    string       `json:"via_id"`
	WaitTime int          `json:"wait_time,string"`
	Status   string       `json:"status"`
	Products []ProductRef `json:"products"`
}

type Avoid struct {
	AvoidID     string `json:"avoid_id"`
	AvoidStatus string `json:"avoid_status"`
}

type Walk struct {
	Allow  bool `json:"allow,string"`
	Min    int  `json:"min,string"`
	Max    int  `json:"max,string"`
	Speed  int  `json:"speed,string"`
	Linear bool `json:"linear"`
}

type TripsRequest struct {
	key string

	Lang              string   `json:"lang"`
	OriginID          string   `json:"origin_id"`
	OriginExtID       string   `json:"origin_ext_id"`
	OriginCoordLat    string   `json:"origin_coord_lat"`
	OriginCoordLong   string   `json:"origin_coord_long"`
	DestID            string   `json:"dest_id"`
	DestExtID         string   `json:"dest_ext_id"`
	DestCoordLat      string   `json:"dest_coord_lat"`
	DestCoordLong     string   `json:"dest_coord_long"`
	Via               []string `json:"via"`
	ViaID             string   `json:"via_id"`
	ViaWaitTime       string   `json:"via_wait_time"`
	Avoid             []string `json:"avoid"`
	AvoidID           string   `json:"avoid_id"`
	ChangeTimePercent string   `json:"change_time_percent"`
	MinChangeTime     string   `json:"min_change_time"`
	MaxChangeTime     string   `json:"max_change_time"`
	AddChangeTime     string   `json:"add_change_time"`
	MaxChange         string   `json:"max_change"`
	Time              time.Time
	SearchForArrival  bool   `json:"search_for_arrival"`
	NumF              string `json:"num_f"`
	NumB              string `json:"num_b"`
	Products          []ProductRef
	AvoidProducts     []ProductRef
	Lines             []string `json:"lines"`
	Context           string   `json:"context"`
	Poly              bool     `json:"poly"`
	Passlist          bool     `json:"passlist"`
	OriginWalk        Walk     `json:"origin_walk"`
	DestWalk          Walk     `json:"dest_walk"`
}

func (r TripsRequest) body() (io.Reader, error) { return nil, nil }
func (r TripsRequest) params() url.Values {
	params := url.Values{}
	if r.key != "" {
		params.Set("key", r.key)
	}
	if r.Lang == "" {
		params.Set("lang", "en")
	} else {
		params.Set("lang", r.Lang)
	}
	if r.OriginID != "" {
		params.Set("originId", r.OriginID)
	}
	if r.OriginExtID != "" {
		params.Set("originExtId", r.OriginExtID)
	}
	if r.OriginCoordLat != "" {
		params.Set("originCoordLat", r.OriginCoordLat)
	}
	if r.OriginCoordLong != "" {
		params.Set("originCoordLong", r.OriginCoordLong)
	}
	if r.DestID != "" {
		params.Set("destId", r.DestID)
	}
	if r.DestExtID != "" {
		params.Set("destExtId", r.DestExtID)
	}
	if r.DestCoordLat != "" {
		params.Set("destCoordLat", r.DestCoordLat)
	}
	if r.DestCoordLong != "" {
		params.Set("destCoordLong", r.DestCoordLong)
	}
	if r.Via != nil && len(r.Via) > 0 {
		// TODO
	}
	if r.ViaID != "" {
		params.Set("viaId", r.ViaID)
	}
	if r.ViaWaitTime != "" {
		params.Set("viaWaitTime", r.ViaWaitTime)
	}
	if r.Avoid != nil && len(r.Avoid) > 0 {
		// TODO
	}
	if r.AvoidID != "" {
		params.Set("avoidId", r.AvoidID)
	}
	if r.ChangeTimePercent != "" {
		params.Set("changeTimePercent", r.ChangeTimePercent)
	}
	if r.MinChangeTime != "" {
		params.Set("minChangeTime", r.MinChangeTime)
	}
	if r.MaxChangeTime != "" {
		params.Set("maxChangeTime", r.MaxChangeTime)
	}
	if r.AddChangeTime != "" {
		params.Set("addChangeTime", r.AddChangeTime)
	}
	if r.MaxChange != "" {
		params.Set("maxChange", r.MaxChange)
	}
	if r.Time != (time.Time{}) {
		params.Set("date", r.Time.In(LocationEuropeStockholm).Format("2006-01-02"))
		params.Set("time", r.Time.In(LocationEuropeStockholm).Format("15:04"))
	}
	if r.SearchForArrival {
		params.Set("searchForArrival", "1")
	} else {
		params.Set("searchForArrival", "0")
	}
	if r.NumF != "" {
		params.Set("numF", r.NumF)
	}
	if r.NumB != "" {
		params.Set("numB", r.NumB)
	}
	if r.AvoidProducts != nil && len(r.AvoidProducts) > 0 {
		// TODO
	}
	if r.Products != nil && len(r.Products) > 0 {
		p := 0

		for _, product := range r.Products {
			p += int(product)
		}

		params.Set("products", strconv.Itoa(p))
	}
	if r.Lines != nil && len(r.Lines) > 0 {
		lines := ""

		for i, l := range r.Lines {
			if l == "" {
				continue
			}

			if i != 0 && lines != "" {
				lines += ","
			}

			lines += l
		}

		params.Set("lines", lines)
	}
	if r.Context != "" {
		params.Set("context", r.Context)
	}
	if r.Poly {
		params.Set("poly", "1")
	} else {
		params.Set("poly", "0")
	}
	if r.Passlist {
		params.Set("passlist", "1")
	} else {
		params.Set("passlist", "0")
	}
	if r.OriginWalk != (Walk{}) {
		allow := "0"
		if r.OriginWalk.Allow {
			allow = "1"
		}
		linear := "0"
		if r.OriginWalk.Linear {
			linear = "1"
		}
		params.Set("originWalk", fmt.Sprintf("%v,%v,%v,%v", allow, strconv.Itoa(r.OriginWalk.Min), strconv.Itoa(r.OriginWalk.Max), linear))
	}
	if r.DestWalk != (Walk{}) {
		allow := "0"
		if r.DestWalk.Allow {
			allow = "1"
		}
		linear := "0"
		if r.DestWalk.Linear {
			linear = "1"
		}
		params.Set("destWalk", fmt.Sprintf("%v,%v,%v,%v", allow, strconv.Itoa(r.DestWalk.Min), strconv.Itoa(r.DestWalk.Max), linear))
	}
	return params
}

type TripResp struct {
	StatusCode int32  `json:"status_code"`
	Message    string `json:"message"`
	ScrB       string `json:"scr_b" xml:"scrB,attr"`
	ScrF       string `json:"scr_f" xml:"scrF,attr"`
	Trip       *Trip  `json:"trips" xml:"Trip"`
}

type TripsResp struct {
	StatusCode int32  `json:"status_code"`
	Message    string `json:"message"`
	ScrB       string `json:"scr_b" xml:"scrB,attr"`
	ScrF       string `json:"scr_f" xml:"scrF,attr"`
	Trips      []Trip `json:"trips" xml:"Trip"`
}

type LegResp struct {
	Leg
	StatusCode int32  `json:"status_code"`
	Message    string `json:"message"`
	ScrB       string `json:"scr_b" xml:"scrB,attr"`
	ScrF       string `json:"scr_f" xml:"scrF,attr"`
}

func (d *TripsResp) CombineWalks() {
	for ti, _ := range d.Trips {
		d.Trips[ti].CombineWalks()
	}
}

type Trip struct {
	Idx         string        `json:"idx" xml:"idx,attr"`
	CtxRecon    string        `json:"ctx_recon" xml:"ctxRecon,attr"`
	Checksum    string        `json:"checksum" xml:"checksum,attr"`
	TripID      string        `json:"trip_id" xml:"tripId,attr"`
	Valid       string        `json:"valid,string" xml:"valid,attr"`
	Duration    string        `json:"duration" xml:"duration,attr"`
	ServiceDays []ServiceDay  `json:"service_days"`
	Legs        []Leg         `json:"legs" xml:"LegList>Leg"`
	Tarriff     []FareSetItem `json:"tarriff,omitempty" xml:"TarriffResult>fareSetItem"`
}

func (trip *Trip) CombineWalks() {
	// combine adjecent walks
	combinedLegs := []Leg{}
	for _, leg := range trip.Legs {
		if leg.Type != "WALK" {
			combinedLegs = append(combinedLegs, leg)

			continue
		}

		var prevWalk *Leg
		if len(combinedLegs) > 0 && combinedLegs[len(combinedLegs)-1].Type == "WALK" {
			prevWalk = &combinedLegs[len(combinedLegs)-1]
		}

		if prevWalk != nil {
			prevWalk.Distance += leg.Distance
			prevWalk.Destination = leg.Destination
		} else {
			combinedLegs = append(combinedLegs, leg)
		}
	}

	// remove short walks
	legs := []Leg{}
	intermediate := false
	for i, leg := range combinedLegs {
		if leg.Type != "WALK" {
			intermediate = true
		}
		if len(combinedLegs)-1 == i {
			intermediate = false
		}
		if leg.Type != "WALK" ||
			((!intermediate && leg.Distance > 40) ||
				(intermediate && leg.Distance > 150)) {
			legs = append(legs, leg)
		}
	}

	trip.Legs = legs
}

func (trip *Trip) EachLegContextual(fn LegContextualFunc) error {
	if len(trip.Legs) == 0 {
		return nil
	}

	prevLeg := &Leg{}
	prevTransportLeg := &Leg{}
	nextLeg := &Leg{}
	nextTransportLeg := &Leg{}
	legCount := len(trip.Legs) - 1

	for i, _ := range trip.Legs {
		leg := &trip.Legs[i]

		if i < legCount {
			nextLeg = &trip.Legs[i+1]
			for _, leg := range trip.Legs[i+1:] {
				if leg.Type != "WALK" {
					nextTransportLeg = &leg
					break
				}
			}
		}

		if err := fn(leg, prevLeg, prevTransportLeg, nextLeg, nextTransportLeg, i); err != nil {
			return err
		}

		nextLeg = &Leg{}
		nextTransportLeg = &Leg{}
		prevLeg = &trip.Legs[i]

		if prevLeg.Type != "WALK" {
			prevTransportLeg = prevLeg
		}
	}

	return nil
}

type ServiceDay struct {
	SDaysR              string `json:"s_days_r" xml:"sDaysR,attr"`
	SDaysI              string `json:"s_days_i" xml:"sDaysI,attr"`
	SDaysB              string `json:"s_days_b" xml:"sDaysB,attr"`
	PlanningPeriodBegin string `json:"planning_period_being" xml:"planningPeriodBeing,attr"`
	PlanningPeriodEnd   string `json:"planning_period_end" xml:"planningPeriodEnd,attr"`
}

type Leg struct {
	Distance      int           `json:"distance" xml:"dist,attr"`
	Type          string        `json:"type" xml:"type,attr"`
	Idx           int           `json:"idx,string" xml:"idx,attr"`
	Cancelled     bool          `json:"cancelled,string" xml:"cancelled,attr"`
	Name          string        `json:"name" xml:"name,attr"`
	Number        int           `json:"number,string" xml:"number,attr"`
	Category      string        `json:"category" xml:"category,attr"`
	Reachable     bool          `json:"reachable,string" xml:"reachable,attr"`
	Direction     string        `json:"direction" xml:"direction,attr"`
	Origin        Location      `json:"origin"`
	Destination   Location      `json:"destination"`
	JourneyDetail JourneyDetail `json:"journey_detail" xml:"JourneyDetailRef"`
	Messages      []Message     `json:"messages,omitempty" xml:"Messages>Message"`
	Notes         []Note        `json:"notes,omitempty" xml:"Notes>Note"`
	JourneyStatus string        `json:"journey_status"`
	Product       *Product      `json:"product,omitempty"`
	Polyline      *Polyline     `json:"polyline,omitempty"`
	Stops         []Stop        `json:"stops,omitempty" xml:"Stops>Stop"`
}

type Location struct {
	ID            string  `json:"id" xml:"id,attr"`
	ExtID         string  `json:"ext_id" xml:"extId,attr"`
	Name          string  `json:"name" xml:"name,attr"`
	Type          string  `json:"type" xml:"type,attr"`
	Lon           float64 `json:"lon" xml:"lon,attr"`
	Lat           float64 `json:"lat" xml:"lat,attr"`
	HasMainMast   bool    `json:"has_main_mast,string" xml:"hasMainMast,attr"`
	MainMastID    string  `json:"main_mast_id" xml:"mainMastId,attr"`
	MainMastExtID string  `json:"main_mast_ext_id" xml:"mainMastExtId,attr"`
	Date          string  `json:"date" xml:"date,attr"`
	RtDate        string  `json:"rt_date" xml:"rtDate,attr"`
	Time          string  `json:"time" xml:"time,attr"`
	RtTime        string  `json:"rt_time" xml:"rtTime,attr"`
	Track         string  `json:"track" xml:"track,attr"`
	PrognosisType string  `json:"prognosis_type" xml:"prognosisType,attr"`
}

func (l Location) ParseTime() (st time.Time, rt time.Time, err error) {
	if l.Date != "" && l.Time != "" {
		st, err = time.ParseInLocation("2006-01-02 15:04:05", l.Date+" "+l.Time, LocationEuropeStockholm)
		if err != nil {
			return
		}
	}

	if l.RtDate != "" && l.RtTime != "" {
		rt, err = time.ParseInLocation("2006-01-02 15:04:05", l.RtDate+" "+l.RtTime, LocationEuropeStockholm)
		if err != nil {
			return
		}
	}

	if rt == (time.Time{}) {
		rt = st
	}

	return
}

type JourneyDetail struct {
	Ref string `json:"ref" xml:"ref,attr"`
}

type Message struct {
	ID        string `json:"id" xml:"id,attr"`
	Act       bool   `json:"act,string" xml:"act,attr"`
	Head      string `json:"head" xml:"head,attr"`
	Text      string `json:"text" xml:"text,attr"`
	Priority  int    `json:"priority,string" xml:"priority,attr"`
	Category  string `json:"category" xml:"category,attr"`
	Products  int    `json:"products,string" xml:"products,attr"`
	StartTime string `json:"start_time" xml:"sTime,attr"`
	StartDate string `json:"start_date" xml:"sDate,attr"`
	EndTime   string `json:"end_time" xml:"eTime,attr"`
	EndDate   string `json:"end_date" xml:"eDate,attr"`
}

type Note struct {
	Priority int    `json:"priority,string" xml:"priority,attr"`
	Text     string `json:"text" xml:",chardata"`
}

type Product struct {
	CategoryCode      int    `json:"category_code,string" xml:"catCode,attr"`
	CategoryIn        string `json:"category_in" xml:"catIn,attr"`
	CategoryOut       string `json:"category_out" xml:"catOut,attr"`
	CateogryOutLocale string `json:"category_out_locale" xml:"catOutL,attr"`
	CategoryOutShort  string `json:"category_out_short" xml:"catOutS,attr"`
	Line              string `json:"line" xml:"line,attr"`
	Name              string `json:"name" xml:"name,attr"`
	Num               int    `json:"num,string" xml:"num,attr"`
	Operator          string `json:"operator" xml:"operator,attr"`
	OperatorCode      string `json:"operator_code" xml:"operatorCode,attr"`
	Admin             string `json:"admin" xml:"admin,attr"`
}

type Polyline struct {
	Type                       string    `json:"type" xml:"type,attr"`
	Dim                        string    `json:"dim" xml:"dim,attr"`
	CoordinatesEncryptedString string    `json:"coordinates_encrypted_string" xml:"crdEncS,attr"`
	Delta                      bool      `json:"delta,string" xml:"delta,attr"`
	Coordinates                []float64 `json:"coordinates,string" xml:"crd"`
}

type Stop struct {
	DepartureDate   string  `json:"departure_date" xml:"depDate,attr"`
	RtDepartureDate string  `json:"rt_departure_date" xml:"rtDepDate,attr"`
	DepartureTime   string  `json:"departure_time" xml:"depTime,attr"`
	RtDepartureTime string  `json:"rt_departure_time" xml:"rtDepTime,attr"`
	ArrivalDate     string  `json:"arrival_date" xml:"arrDate,attr"`
	RtArrivalDate   string  `json:"rt_arrival_date" xml:"rtArrDate,attr"`
	ArrivalTime     string  `json:"arrival_time" xml:"arrTime,attr"`
	RtArrivalTime   string  `json:"rt_arrival_time" xml:"rtArrTime,attr"`
	RouteIdx        int     `json:"route_idx,string" xml:"routeIdx,attr"`
	Name            string  `json:"name" xml:"name,attr"`
	ID              string  `json:"id" xml:"id,attr"`
	ExtId           string  `json:"ext_id" xml:"extId,attr"`
	Lon             float64 `json:"lon" xml:"lon,attr"`
	Lat             float64 `json:"lat" xml:"lat,attr"`
	HasMainMast     bool    `json:"has_main_mast,string" xml:"hasMainMast,attr"`
	MainMastID      string  `json:"main_mast_id" xml:"mainMastId,attr"`
	MainMastExtID   string  `json:"main_mast_ext_id" xml:"mainMastExtId,attr"`
	DepartureTrack  string  `json:"departure_track" xml:"depTrack,attr"`
	ArrivalTrack    string  `json:"arrival_track" xml:"arrTrack,attr"`
}

func (s Stop) ParseArrival() (st time.Time, rt time.Time, err error) {
	if s.ArrivalDate != "" && s.ArrivalTime != "" {
		st, err = time.ParseInLocation("2006-01-02 15:04:05", s.ArrivalDate+" "+s.ArrivalTime, LocationEuropeStockholm)
		if err != nil {
			return
		}
	}

	if s.RtArrivalDate != "" && s.RtArrivalTime != "" {
		rt, err = time.ParseInLocation("2006-01-02 15:04:05", s.RtArrivalDate+" "+s.RtArrivalTime, LocationEuropeStockholm)
		if err != nil {
			return
		}
	}

	if rt == (time.Time{}) {
		rt = st
	}

	return
}

func (s Stop) ParseDeparture() (st time.Time, rt time.Time, err error) {
	if s.DepartureDate != "" && s.DepartureTime != "" {
		st, err = time.ParseInLocation("2006-01-02 15:04:05", s.DepartureDate+" "+s.DepartureTime, LocationEuropeStockholm)
		if err != nil {
			return
		}
	}

	if s.RtDepartureDate != "" && s.RtDepartureTime != "" {
		rt, err = time.ParseInLocation("2006-01-02 15:04:05", s.RtDepartureDate+" "+s.RtDepartureTime, LocationEuropeStockholm)
		if err != nil {
			return
		}
	}

	if rt == (time.Time{}) {
		rt = st
	}

	return
}

type FareSetItem struct {
	Name        string     `json:"name" xml:"name,attr"`
	Description string     `json:"desc" xml:"desc,attr"`
	Fares       []FareItem `json:"fares" xml:"fareItem"`
}

type FareItem struct {
	Name        string `json:"name" xml:"name,attr"`
	Description string `json:"desc" xml:"desc,attr"`
	Currency    string `json:"currency" xml:"cur,attr"`
	Price       int    `json:"price,string" xml:"price,attr"`
}
