package trafiklab

import (
	"testing"
)

func TestTravelplannerTripsNoKeys(t *testing.T) {
	_, err := clientNoKeys.Travelplanner.Trips(&TripsRequest{
		OriginID: "9192",
		DestID:   "9306",
	})
	if err != ErrMissingAPIKey {
		t.Errorf("expected ErrMissingAPIKey")
	}
}

func TestTravelplannerTrips(t *testing.T) {
	resp, err := client.Travelplanner.Trips(&TripsRequest{
		OriginID: "9192",
		DestID:   "9306",
		Poly:     true,
	})
	if err != nil {
		t.Errorf("caught unexpected error, %v", err)
	}

	for _, trip := range resp.Trips {
		for _, leg := range trip.Legs {
			if leg.Polyline != nil {
				leg.Polyline.LatLng()
			}
		}
	}

}

func TestTravelplannerReconstructionNoKeys(t *testing.T) {
	_, err := clientNoKeys.Travelplanner.Reconstruction("")
	if err != ErrMissingAPIKey {
		t.Errorf("expected ErrMissingAPIKey")
	}
}

func TestTravelplannerReconstruction(t *testing.T) {
	_, err := client.Travelplanner.Reconstruction("T$A=1@O=Slussen@L=400102011@a=128@$A=1@O=T-Centralen@L=400101051@a=128@$202002060907$202002060911$        $§W$A=1@O=T-Centralen@L=400101051@a=128@$A=1@O=T-Centralen@L=400103051@a=128@$202002060911$202002060915$$§T$A=1@O=T-Centralen@L=400103051@a=128@$A=1@O=Västra skogen@L=400103201@a=128@$202002060915$202002060923$        $ Checksum:CAE2C55C_4")
	if err != nil {
		t.Errorf("caught unexpected error, %v", err)
	}
}

func TestTravelplannerJourneyDetailNoKeys(t *testing.T) {
	_, err := clientNoKeys.Travelplanner.JourneyDetail(&JourneyDetailRequest{
		ID: "1|33724|0|74|6022020",
	})
	if err != ErrMissingAPIKey {
		t.Errorf("expected ErrMissingAPIKey")
	}
}

func TestTravelplannerJourneyDetail(t *testing.T) {
	_, err := client.Travelplanner.JourneyDetail(&JourneyDetailRequest{
		ID: "1|33724|0|74|6022020",
	})
	if err != nil {
		t.Errorf("caught unexpected error, %v", err)
	}
}
