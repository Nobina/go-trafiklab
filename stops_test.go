package trafiklab

import (
	"testing"
)

func TestStopsNearbyNoKeys(t *testing.T) {
	_, err := clientNoKeys.Stops.Nearby(&StopsNearbyRequest{
		OriginCoordLat:  "59.348572",
		OriginCoordLong: "17.997520",
	})
	if err != ErrMissingAPIKey {
		t.Errorf("expected ErrMissingAPIKey")
	}
}

func TestStopsNearby(t *testing.T) {
	_, err := client.Stops.Nearby(&StopsNearbyRequest{
		OriginCoordLat:  "59.348572",
		OriginCoordLong: "17.997520",
	})
	if err != nil {
		t.Errorf("caught unexpected error, %v", err)
	}
}

func TestStopsQueryNoKeys(t *testing.T) {
	_, err := clientNoKeys.Stops.Query(&StopsQueryRequest{
		SearchString: "slussen",
	})
	if err != ErrMissingAPIKey {
		t.Errorf("expected ErrMissingAPIKey")
	}
}

func TestStopsQuery(t *testing.T) {
	_, err := client.Stops.Query(&StopsQueryRequest{
		SearchString: "slussen",
	})
	if err != nil {
		t.Errorf("caught unexpected error, %v", err)
	}
}
