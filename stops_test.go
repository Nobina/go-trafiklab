package trafiklab

import (
	"testing"
)

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
