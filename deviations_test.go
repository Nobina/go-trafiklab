package trafiklab

import (
	"testing"
)

func TestDeviationsNoKeys(t *testing.T) {
	_, err := clientNoKeys.Deviations(&DeviationsRequest{
		SiteID: "9192",
	})
	if err != ErrMissingAPIKey {
		t.Errorf("expected ErrMissingAPIKey")
	}
}

func TestDeviations(t *testing.T) {
	_, err := client.Deviations(&DeviationsRequest{
		SiteID: "9192",
	})
	if err != nil {
		t.Errorf("caught unexpected error, %v", err)
	}
}
