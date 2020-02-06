package trafiklab

import (
	"net/http"
	"os"
	"testing"
)

var (
	client       *Client
	clientNoKeys *Client
)

func TestMain(m *testing.M) {
  client = NewClient(
		WithHTTPClient(http.DefaultClient),
		WithDeparturesAPIKey(os.Getenv("TRAFIKLAB_DEPARTURES_KEY")),
		WithDeviationsAPIKey(os.Getenv("TRAFIKLAB_DEVIATIONS_KEY")),
		WithStopsNearbyAPIKey(os.Getenv("TRAFIKLAB_NEARBY_STOPS_KEY")),
		WithStopsQueryAPIKey(os.Getenv("TRAFIKLAB_TYPEAHEAD_KEY")),
		WithTrafficStatusAPIKey(os.Getenv("TRAFIKLAB_TRAFFIC_SITUATION_KEY")),
		WithTravelplannerAPIKey(os.Getenv("TRAFIKLAB_TRAVEL_PLANNER_KEY")),
	)

	clientNoKeys = NewClient(
		WithHTTPClient(http.DefaultClient),
		WithDeparturesAPIKey(""),
		WithDeviationsAPIKey(""),
		WithStopsNearbyAPIKey(""),
		WithStopsQueryAPIKey(""),
		WithTrafficStatusAPIKey(""),
		WithTravelplannerAPIKey(""),
	)
	code := m.Run()
	os.Exit(code)
}
