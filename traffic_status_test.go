package trafiklab

import (
	"testing"
)

func TestTrafficStatusNoKeys(t *testing.T) {
	_, err := clientNoKeys.TrafficStatus()
	if err != ErrMissingAPIKey {
		t.Errorf("expected ErrMissingAPIKey")
	}
}

func TestTrafficStatus(t *testing.T) {
	_, err := client.TrafficStatus()
	if err != nil {
		t.Errorf("caught unexpected error, %v", err)
	}
}
