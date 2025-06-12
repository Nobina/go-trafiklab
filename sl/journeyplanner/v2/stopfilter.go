package journeyplanner

import (
	"errors"
	"fmt"
)

// StopFilter is a filter for the stop finder.
// It uses bitwise operations to combine multiple filters.
type StopFilter int
var StopFilterNone StopFilter = 0
const (
	StopFilterSuburb StopFilter = 1 << iota
	StopFilterStop // 2
	StopFilterStreet // 4
	StopFilterAddress // 8
	StopFilterUnknown // 16
	StopFilterPointOfInterest // 32
)

var stopFilterNameToValue = map[string]StopFilter{
	"none": StopFilterNone,
	"suburb": StopFilterSuburb,
	"stop": StopFilterStop,
	"street": StopFilterStreet,
	"singlehouse": StopFilterAddress,
	"unknown2": StopFilterUnknown,
	"poi": StopFilterPointOfInterest,
}

var stopFilterValueToName = map[StopFilter]string{
	StopFilterNone: "none",
	StopFilterSuburb: "suburb",
	StopFilterStop: "stop",
	StopFilterStreet: "street",
	StopFilterAddress: "singlehouse",
	StopFilterUnknown: "unknown2",
	StopFilterPointOfInterest: "poi",
}

func StopFilterFromString(s ...string) (StopFilter, error) {
	if len(s) == 0 {
		return -1, errors.New("no stop filter provided")
	}

	sf := StopFilterNone
	for _, s := range s {
		if v, ok := stopFilterNameToValue[s]; ok {
			sf = sf.Add(v)
		} else {
			return -1, fmt.Errorf("invalid stop filter: %s", s)
		}
	}
	return sf, nil
}


func (sf StopFilter) Has(filter StopFilter) bool { return sf&filter != 0 }
func (sf StopFilter) Add(filter StopFilter) StopFilter { return sf | filter }
func (sf StopFilter) Remove(filter StopFilter) StopFilter { return sf &^ filter }
