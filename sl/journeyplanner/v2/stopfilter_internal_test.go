package journeyplanner

import (
	"testing"
)
func TestAllValuesAccountedFor(t *testing.T) {
	for _, v := range stopFilterNameToValue {
		if _, ok := stopFilterValueToName[v]; !ok {
			t.Errorf("value %v not accounted for", v)
		}
	}

	for _, v := range stopFilterValueToName {
		if _, ok := stopFilterNameToValue[v]; !ok {
			t.Errorf("value %v not accounted for", v)
		}
	}
}
