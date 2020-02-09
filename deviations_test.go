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

func TestDeviationsTimeConvert(t *testing.T) {
	resp, err := client.Deviations(&DeviationsRequest{
		SiteID: "9192",
	})
	if err != nil {
		t.Errorf("caught unexpected error, %v", err)
	}

	for _, d := range resp.Data {
		if _, err := d.FromDate(); err != nil {
			t.Errorf("error caught converting from date, %v", err)
		}
		if _, err := d.ToDate(); err != nil {
			t.Errorf("error caught converting to date, %v", err)
		}
		if _, err := d.UpdatedDate(); err != nil {
			t.Errorf("error caught converting updated date, %v", err)
		}
	}
}
