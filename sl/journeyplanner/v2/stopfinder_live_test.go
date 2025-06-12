package journeyplanner_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/nobina/go-trafiklab/sl/journeyplanner/v2"
)

func TestStopFinder(t *testing.T) {
	_ = godotenv.Load()
	apiKey := os.Getenv("SL_API_KEY")
	//TODO(pb): API key not needed anymore?
	client := journeyplanner.NewClient(&journeyplanner.JourneyPlannerConfig{
		APIKey: apiKey,
	}, http.DefaultClient)

	ctx := context.Background()

	t.Run("basic stop search", func(t *testing.T) {
		req := journeyplanner.NewStopFinderSearchRequest("T-Centralen", []string{"stop"})
		resp, err := client.StopFinder(ctx, req)
		if err != nil {
			t.Fatalf("StopFinder failed: %v", err)
		}

		if len(resp.StopLocations) == 0 {
			t.Fatal("expected at least one stop location")
		}
	})

	t.Run("search with multiple types", func(t *testing.T) {
		req := journeyplanner.NewStopFinderSearchRequest("Stockholm", []string{"stop", "poi"})
		resp, err := client.StopFinder(ctx, req)
		if err != nil {
			t.Fatalf("StopFinder failed: %v", err)
		}

		if len(resp.StopLocations) == 0 {
			t.Fatal("expected at least one stop location")
		}

		// Verify we got different types of results
		types := make(map[string]bool)
		for _, loc := range resp.StopLocations {
			types[loc.Type] = true
		}
		if len(types) < 2 {
			t.Error("expected results of different types")
		}
	})

	t.Run("invalid request - missing name", func(t *testing.T) {
		req := journeyplanner.NewStopFinderSearchRequest("", []string{"stop"})

		_, err := client.StopFinder(ctx, req)
		if err == nil {
			t.Error("expected error for missing name")
		}
	})

	t.Run("invalid request - missing type", func(t *testing.T) {
		req := journeyplanner.NewStopFinderSearchRequest("T-Centralen", []string{})

		_, err := client.StopFinder(ctx, req)
		if err == nil {
			t.Error("expected error for missing type")
		}
	})

	t.Run("response validation", func(t *testing.T) {
		req := journeyplanner.NewStopFinderSearchRequest("T-Centralen", []string{"stop"})

		resp, err := client.StopFinder(ctx, req)
		if err != nil {
			t.Fatalf("StopFinder failed: %v", err)
		}

		for _, loc := range resp.StopLocations {
			// Verify required fields are present
			if loc.ID == "" {
				t.Error("expected non-empty ID")
			}
			if loc.Name == "" {
				t.Error("expected non-empty Name")
			}
			if loc.Type == "" {
				t.Error("expected non-empty Type")
			}

			// Verify coordinates are valid
			if loc.Coordinates.Latitude == 0 && loc.Coordinates.Longitude == 0 {
				t.Error("expected non-zero coordinates")
			}
		}
	})
}


func TestStopFinder_CoordSearch(t *testing.T) {
	//TODO(pb): API key not needed anymore?
	client := journeyplanner.NewClient(&journeyplanner.JourneyPlannerConfig{
	}, http.DefaultClient)

	ctx := context.Background()

	t.Run("basic coord search", func(t *testing.T) {
		req := journeyplanner.NewStopFinderPosRequest(journeyplanner.LatLng{
			Latitude: 59.332581,
			Longitude: 18.064924,
		}, []string{"stop"})
		resp, err := client.StopFinder(ctx, req)
		if err != nil {
			t.Fatalf("StopFinder failed: %v", err)
		}

		if len(resp.StopLocations) == 0 {
			t.Fatal("expected at least one stop location")
		}
	})
}
