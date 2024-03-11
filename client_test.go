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
		WithStopsQueryAPIKey(os.Getenv("TRAFIKLAB_TYPEAHEAD_KEY")),
		WithTrafficStatusAPIKey(os.Getenv("TRAFIKLAB_TRAFFIC_SITUATION_KEY")),
	)

	clientNoKeys = NewClient(
		WithHTTPClient(http.DefaultClient),
		WithDeparturesAPIKey(""),
		WithDeviationsAPIKey(""),
		WithStopsQueryAPIKey(""),
		WithTrafficStatusAPIKey(""),
	)
	code := m.Run()
	os.Exit(code)
}
