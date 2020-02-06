package trafiklab

import (
	"testing"
)

func TestDeparturesNoKeys(t *testing.T) {
	_, err := clientNoKeys.Departures(&DeparturesRequest{
		SiteID: "9192",
	})
	if err != ErrMissingAPIKey {
		t.Errorf("expected ErrMissingAPIKey")
	}
}

func TestDepartures(t *testing.T) {
	_, err := client.Departures(&DeparturesRequest{
		SiteID: "9192",
	})
	if err != nil {
		t.Errorf("caught unexpected error, %v", err)
	}
}
