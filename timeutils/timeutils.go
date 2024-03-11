package timeutils

import (
	"log"
	"time"
)

const sweden = "Europe/Stockholm"

var _ = EuropeStockholm() // crash on init if location not available

// GetDefaultLocation gets the location with the correct timezone that we should use
// Panics if locale is not found, the only reason this should happen is if we're
// on an alpine docker image and the timezone data is not installed
func EuropeStockholm() *time.Location {
	loc, err := time.LoadLocation(sweden)
	if err != nil {
		log.Fatalf("Could not load location, something is very broken: %s", err.Error())
	}
	return loc
}
