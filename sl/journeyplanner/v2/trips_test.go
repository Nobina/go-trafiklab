package journeyplanner

import (
	"testing"
	"time"

	"github.com/nobina/go-trafiklab/timeutils"
)

func TestTripsRequest_toParams_TimestampMapping(t *testing.T) {
	// Test cases for timestamp and dep/arr parameter mapping
	tests := []struct {
		name                   string
		at                     time.Time
		tripDateTimeDepArr     string
		expectedDateParam      string
		expectedTimeParam      string
		expectedDepArrParam    string
		expectTimezoneConversion bool
	}{
		{
			name:               "Departure search with UTC time",
			at:                 time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			tripDateTimeDepArr: "dep",
			expectedDateParam:  "20240115", // Should be in YYYYMMDD format
			expectedTimeParam:  "1130",     // Should be converted to Stockholm time (UTC+1) and in HHMM format
			expectedDepArrParam: "dep",
			expectTimezoneConversion: true,
		},
		{
			name:               "Arrival search with UTC time",
			at:                 time.Date(2024, 1, 15, 14, 45, 0, 0, time.UTC),
			tripDateTimeDepArr: "arr",
			expectedDateParam:  "20240115", // Should be in YYYYMMDD format
			expectedTimeParam:  "1545",     // Should be converted to Stockholm time (UTC+1) and in HHMM format
			expectedDepArrParam: "arr",
			expectTimezoneConversion: true,
		},
		{
			name:               "Default departure (empty dep/arr)",
			at:                 time.Date(2024, 6, 15, 8, 15, 0, 0, time.UTC),
			tripDateTimeDepArr: "",
			expectedDateParam:  "20240615", // Should be in YYYYMMDD format
			expectedTimeParam:  "1015",     // Should be converted to Stockholm time (UTC+2 in summer) and in HHMM format
			expectedDepArrParam: "",        // Should not be set when empty
			expectTimezoneConversion: true,
		},
		{
			name:               "Stockholm time input",
			at:                 time.Date(2024, 1, 15, 11, 30, 0, 0, timeutils.EuropeStockholm()),
			tripDateTimeDepArr: "dep",
			expectedDateParam:  "20240115",
			expectedTimeParam:  "1130", // Should remain the same since already in Stockholm time
			expectedDepArrParam: "dep",
			expectTimezoneConversion: false, // No conversion needed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &TripsRequest{
				At:                 tt.at,
				NumTrips:           1,
				TypeOrigin:         "any",
				NameOrigin:         "Stockholm Central",
				TypeDestination:    "any",
				NameDestination:    "Arlanda",
				TripDateTimeDepArr: tt.tripDateTimeDepArr,
			}

			params := req.toParams()

			// Check correct parameter names are used (API expects itd_date and itd_time)
			if params.Get("itd_date") == "" {
				t.Errorf("Expected 'itd_date' parameter to be set, but it's missing. Found 'date': %s", params.Get("date"))
			}
			if params.Get("itd_time") == "" {
				t.Errorf("Expected 'itd_time' parameter to be set, but it's missing. Found 'time': %s", params.Get("time"))
			}

			// Check date format (should be YYYYMMDD, not YYYY-MM-DD)
			actualDate := params.Get("itd_date")
			if actualDate != tt.expectedDateParam {
				t.Errorf("Expected date parameter '%s', got '%s'", tt.expectedDateParam, actualDate)
			}

			// Check time format (should be HHMM, not HH:MM)
			actualTime := params.Get("itd_time")
			if actualTime != tt.expectedTimeParam {
				t.Errorf("Expected time parameter '%s', got '%s'", tt.expectedTimeParam, actualTime)
			}

			// Check dep/arr parameter
			actualDepArr := params.Get("itd_trip_date_time_dep_arr")
			if actualDepArr != tt.expectedDepArrParam {
				t.Errorf("Expected dep/arr parameter '%s', got '%s'", tt.expectedDepArrParam, actualDepArr)
			}

			// Test timezone conversion
			if tt.expectTimezoneConversion {
				// Convert the input time to Stockholm timezone and check if it matches
				stockholmTime := tt.at.In(timeutils.EuropeStockholm())
				expectedDate := stockholmTime.Format("20060102")
				expectedTime := stockholmTime.Format("1504")

				if actualDate != expectedDate {
					t.Errorf("Timezone conversion failed for date. Expected '%s', got '%s'", expectedDate, actualDate)
				}
				if actualTime != expectedTime {
					t.Errorf("Timezone conversion failed for time. Expected '%s', got '%s'", expectedTime, actualTime)
				}
			}
		})
	}
}

func TestTripsRequest_toParams_ParameterNames(t *testing.T) {
	// Test that we use the correct parameter names according to the API spec
	req := &TripsRequest{
		At:                 time.Now(),
		NumTrips:           1,
		TypeOrigin:         "any",
		NameOrigin:         "Test Origin",
		TypeDestination:    "any",
		NameDestination:    "Test Destination",
		TripDateTimeDepArr: "dep",
	}

	params := req.toParams()

	// Check that we're using the correct parameter names according to API spec
	correctParams := map[string]bool{
		"itd_date":                    true,
		"itd_time":                    true,
		"itd_trip_date_time_dep_arr":  true,
	}

	incorrectParams := map[string]bool{
		"date":                false, // Wrong parameter name
		"time":                false, // Wrong parameter name
		"searchForArrival":    false, // V1 parameter, not used in V2
	}

	for param := range correctParams {
		if params.Get(param) == "" {
			t.Errorf("Missing required parameter: %s", param)
		}
	}

	for param := range incorrectParams {
		if params.Get(param) != "" {
			t.Errorf("Found incorrect parameter '%s' with value '%s', this should not be used in V2 API", param, params.Get(param))
		}
	}
}

func TestTripsRequest_DepArrLogic(t *testing.T) {
	// Test the dep/arr logic specifically
	testTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name           string
		depArr         string
		expectedParam  string
		shouldBeSet    bool
	}{
		{
			name:          "Departure search",
			depArr:        "dep",
			expectedParam: "dep",
			shouldBeSet:   true,
		},
		{
			name:          "Arrival search",
			depArr:        "arr",
			expectedParam: "arr",
			shouldBeSet:   true,
		},
		{
			name:          "Empty (default to departure)",
			depArr:        "",
			expectedParam: "",
			shouldBeSet:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &TripsRequest{
				At:                 testTime,
				NumTrips:           1,
				TypeOrigin:         "any",
				NameOrigin:         "Origin",
				TypeDestination:    "any",
				NameDestination:    "Destination",
				TripDateTimeDepArr: tt.depArr,
			}

			params := req.toParams()
			actualParam := params.Get("itd_trip_date_time_dep_arr")

			if tt.shouldBeSet {
				if actualParam != tt.expectedParam {
					t.Errorf("Expected parameter 'itd_trip_date_time_dep_arr' to be '%s', got '%s'", tt.expectedParam, actualParam)
				}
			} else {
				if actualParam != "" {
					t.Errorf("Expected parameter 'itd_trip_date_time_dep_arr' to be empty, got '%s'", actualParam)
				}
			}
		})
	}
}

func TestTripsRequest_DepartAfterArriveBeforeMapping(t *testing.T) {
	// Test the mapping from depart_after/arrive_before concepts to Trafiklab API parameters
	tests := []struct {
		name                    string
		scenario                string
		at                      time.Time
		tripDateTimeDepArr      string
		expectedBehavior        string
		expectedAPICall         string
	}{
		{
			name:                "depart_after scenario",
			scenario:            "User wants trips that depart AFTER 17:30",
			at:                  time.Date(2024, 1, 15, 17, 30, 0, 0, timeutils.EuropeStockholm()),
			tripDateTimeDepArr:  "dep",
			expectedBehavior:    "Should find trips departing at or after 17:30",
			expectedAPICall:     "itd_trip_date_time_dep_arr=dep with time=17:30",
		},
		{
			name:                "arrive_before scenario",
			scenario:            "User wants trips that arrive BEFORE 09:00",
			at:                  time.Date(2024, 1, 15, 9, 0, 0, 0, timeutils.EuropeStockholm()),
			tripDateTimeDepArr:  "arr",
			expectedBehavior:    "Should find trips arriving at or before 09:00",
			expectedAPICall:     "itd_trip_date_time_dep_arr=arr with time=09:00",
		},
		{
			name:                "pagination scenario",
			scenario:            "User loaded 3 trips, last one departed at 18:15, now wants next 3",
			at:                  time.Date(2024, 1, 15, 18, 15, 0, 0, timeutils.EuropeStockholm()),
			tripDateTimeDepArr:  "dep",
			expectedBehavior:    "Should find trips departing AFTER 18:15 (not including 18:15)",
			expectedAPICall:     "itd_trip_date_time_dep_arr=dep with time=18:15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &TripsRequest{
				At:                 tt.at,
				NumTrips:           3,
				TypeOrigin:         "any",
				NameOrigin:         "9091001000001002", // Stockholm Central
				TypeDestination:    "any",
				NameDestination:    "9091001001002800", // Test destination
				TripDateTimeDepArr: tt.tripDateTimeDepArr,
			}

			params := req.toParams()

			// Check that the parameters are correctly set
			expectedDate := tt.at.In(timeutils.EuropeStockholm()).Format("20060102")
			expectedTime := tt.at.In(timeutils.EuropeStockholm()).Format("1504")

			actualDate := params.Get("itd_date")
			actualTime := params.Get("itd_time")
			actualDepArr := params.Get("itd_trip_date_time_dep_arr")

			if actualDate != expectedDate {
				t.Errorf("Expected date '%s', got '%s'", expectedDate, actualDate)
			}

			if actualTime != expectedTime {
				t.Errorf("Expected time '%s', got '%s'", expectedTime, actualTime)
			}

			if actualDepArr != tt.tripDateTimeDepArr {
				t.Errorf("Expected dep/arr parameter '%s', got '%s'", tt.tripDateTimeDepArr, actualDepArr)
			}

			t.Logf("Scenario: %s", tt.scenario)
			t.Logf("Expected behavior: %s", tt.expectedBehavior)
			t.Logf("API call: %s", tt.expectedAPICall)
			t.Logf("Generated params: itd_date=%s, itd_time=%s, itd_trip_date_time_dep_arr=%s",
				actualDate, actualTime, actualDepArr)
		})
	}
}

func TestTripsRequest_UnixTimestampConversion(t *testing.T) {
	// Test conversion from Unix timestamp (as used in depart_after/arrive_before) to Trafiklab API
	tests := []struct {
		name               string
		unixTimestamp      int64
		depArr             string
		expectedDate       string
		expectedTime       string
	}{
		{
			name:          "depart_after Unix timestamp",
			unixTimestamp: 1750089240, // This is the example from the user's curl request
			depArr:        "dep",
			expectedDate:  "20250426", // Expected date in YYYYMMDD format
			expectedTime:  "1234",     // Expected time in HHMM format (converted to Stockholm time)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert Unix timestamp to time.Time
			inputTime := time.Unix(tt.unixTimestamp, 0)

			req := &TripsRequest{
				At:                 inputTime,
				NumTrips:           3,
				TypeOrigin:         "any",
				NameOrigin:         "9091001000001002",
				TypeDestination:    "any",
				NameDestination:    "9091001001002800",
				TripDateTimeDepArr: tt.depArr,
			}

			params := req.toParams()

			actualDate := params.Get("itd_date")
			actualTime := params.Get("itd_time")
			actualDepArr := params.Get("itd_trip_date_time_dep_arr")

			// Convert to Stockholm time for comparison
			stockholmTime := inputTime.In(timeutils.EuropeStockholm())
			expectedDate := stockholmTime.Format("20060102")
			expectedTime := stockholmTime.Format("1504")

			if actualDate != expectedDate {
				t.Errorf("Expected date '%s', got '%s'", expectedDate, actualDate)
			}

			if actualTime != expectedTime {
				t.Errorf("Expected time '%s', got '%s'", expectedTime, actualTime)
			}

			if actualDepArr != tt.depArr {
				t.Errorf("Expected dep/arr parameter '%s', got '%s'", tt.depArr, actualDepArr)
			}

			t.Logf("Unix timestamp %d converted to Stockholm time: %s",
				tt.unixTimestamp, stockholmTime.Format("2006-01-02 15:04:05"))
			t.Logf("API params: itd_date=%s, itd_time=%s, itd_trip_date_time_dep_arr=%s",
				actualDate, actualTime, actualDepArr)
		})
	}
}

func TestTripsRequest_PaginationScenario(t *testing.T) {
	// Test a realistic pagination scenario
	// First request: Get initial trips
	// Second request: Get next trips after the last departure time

	// Simulate the timestamp of the last trip from the first request
	lastDepartureTime := time.Unix(1750089240, 0) // User's example timestamp

	// For pagination, we want trips that depart AFTER this time
	nextPageReq := &TripsRequest{
		At:                 lastDepartureTime,
		NumTrips:           3,
		TypeOrigin:         "any",
		NameOrigin:         "9091001000001002",
		TypeDestination:    "any",
		NameDestination:    "9091001001002800",
		TripDateTimeDepArr: "dep", // We want departures after this time
	}

	params := nextPageReq.toParams()

	// The issue mentioned is that the API returns trips that depart BEFORE the specified time
	// This test verifies the parameters are set correctly for the Trafiklab API

	stockholmTime := lastDepartureTime.In(timeutils.EuropeStockholm())
	expectedDate := stockholmTime.Format("20060102")
	expectedTime := stockholmTime.Format("1504")

	actualDate := params.Get("itd_date")
	actualTime := params.Get("itd_time")
	actualDepArr := params.Get("itd_trip_date_time_dep_arr")

	if actualDate != expectedDate {
		t.Errorf("Expected date '%s', got '%s'", expectedDate, actualDate)
	}

	if actualTime != expectedTime {
		t.Errorf("Expected time '%s', got '%s'", expectedTime, actualTime)
	}

	if actualDepArr != "dep" {
		t.Errorf("Expected dep/arr parameter 'dep', got '%s'", actualDepArr)
	}

	t.Logf("Pagination scenario:")
	t.Logf("Last departure time: %s UTC", lastDepartureTime.Format("2006-01-02 15:04:05"))
	t.Logf("Stockholm time: %s", stockholmTime.Format("2006-01-02 15:04:05"))
	t.Logf("Trafiklab API params: itd_date=%s, itd_time=%s, itd_trip_date_time_dep_arr=%s",
		actualDate, actualTime, actualDepArr)
	t.Logf("Expected behavior: Should return trips departing AFTER %s", stockholmTime.Format("15:04"))
}

func TestTripsRequest_ParameterCorrectness(t *testing.T) {
	// Verify that the SDK generates the correct parameter names and formats
	testTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	req := &TripsRequest{
		At:                 testTime,
		NumTrips:           3,
		TypeOrigin:         "any",
		NameOrigin:         "Origin",
		TypeDestination:    "any",
		NameDestination:    "Destination",
		TripDateTimeDepArr: "dep",
	}

	params := req.toParams()

	// According to the API documentation, these should be the parameter names
	requiredParams := map[string]string{
		"itd_date":                   "20240115", // YYYYMMDD format
		"itd_time":                   "1130",     // Stockholm time (UTC+1) in HHMM format
		"itd_trip_date_time_dep_arr": "dep",
	}

	for paramName, expectedValue := range requiredParams {
		actualValue := params.Get(paramName)
		if actualValue != expectedValue {
			t.Errorf("Parameter '%s': expected '%s', got '%s'", paramName, expectedValue, actualValue)
		}
	}

	// Verify we're not using the old parameter names
	deprecatedParams := []string{"date", "time", "searchForArrival"}
	for _, paramName := range deprecatedParams {
		if params.Get(paramName) != "" {
			t.Errorf("Found deprecated parameter '%s' with value '%s'", paramName, params.Get(paramName))
		}
	}
}

func TestTripsRequest_ParameterConversion(t *testing.T) {
	// Test that TripsRequest.toParams() generates correct Trafiklab API parameters
	tests := []struct {
		name                string
		request             *TripsRequest
		expectedParams      map[string]string
		description         string
	}{
		{
			name: "Basic departure search",
			request: &TripsRequest{
				At:                 time.Date(2024, 1, 15, 17, 30, 0, 0, timeutils.EuropeStockholm()),
				NumTrips:           3,
				TypeOrigin:         "any",
				NameOrigin:         "9091001000001002",
				TypeDestination:    "any",
				NameDestination:    "9091001001002800",
				TripDateTimeDepArr: "dep",
			},
			expectedParams: map[string]string{
				"itd_date":                   "20240115",
				"itd_time":                   "1730",
				"itd_trip_date_time_dep_arr": "dep",
				"calc_number_of_trips":       "3",
			},
			description: "Should search for departures at/around 17:30",
		},
		{
			name: "Arrival search",
			request: &TripsRequest{
				At:                 time.Date(2024, 1, 15, 9, 0, 0, 0, timeutils.EuropeStockholm()),
				NumTrips:           3,
				TypeOrigin:         "any",
				NameOrigin:         "9091001000001002",
				TypeDestination:    "any",
				NameDestination:    "9091001001002800",
				TripDateTimeDepArr: "arr",
			},
			expectedParams: map[string]string{
				"itd_date":                   "20240115",
				"itd_time":                   "0900",
				"itd_trip_date_time_dep_arr": "arr",
				"calc_number_of_trips":       "3",
			},
			description: "Should search for arrivals at/before 09:00",
		},
		{
			name: "Pagination scenario - depart_after",
			request: &TripsRequest{
				At:                 time.Unix(1750089240, 0), // Unix timestamp from user's example
				NumTrips:           3,
				TypeOrigin:         "any",
				NameOrigin:         "9091001000001002",
				TypeDestination:    "any",
				NameDestination:    "9091001001002800",
				TripDateTimeDepArr: "dep", // This is the key - does "dep" mean "after" or "at"?
			},
			expectedParams: map[string]string{
				"itd_date":                   "20250616", // Stockholm time conversion
				"itd_time":                   "1754",
				"itd_trip_date_time_dep_arr": "dep",
				"calc_number_of_trips":       "3",
			},
			description: "CRITICAL: When user wants trips departing AFTER 17:54, does dep=17:54 give trips after or at 17:54?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := tt.request.toParams()

			// Verify expected parameters
			for paramName, expectedValue := range tt.expectedParams {
				actualValue := params.Get(paramName)
				if actualValue != expectedValue {
					t.Errorf("Parameter '%s': expected '%s', got '%s'",
						paramName, expectedValue, actualValue)
				}
			}

			// Log the full query for manual inspection
			t.Logf("Description: %s", tt.description)
			t.Logf("Full query string: %s", params.Encode())
		})
	}
}

func TestTrips_PaginationBehavior(t *testing.T) {
	// This test simulates the actual pagination scenario described by the user

	// Scenario: User gets initial 3 trips, last one departs at timestamp 1750089240
	// Then user wants NEXT 3 trips (departing AFTER that time)
	// But the API returns trips that depart BEFORE that time (20-30 minutes earlier)

	lastDepartureTime := time.Unix(1750089240, 0)

	// This is what your API does when user requests next page with depart_after
	nextPageRequest := &TripsRequest{
		At:                 lastDepartureTime, // Set to the last departure time
		NumTrips:           3,
		TypeOrigin:         "any",
		NameOrigin:         "9091001000001002", // Stockholm Central
		TypeDestination:    "any",
		NameDestination:    "9091001001002800",
		TripDateTimeDepArr: "dep", // Request departures
	}

	params := nextPageRequest.toParams()

	stockholmTime := lastDepartureTime.In(timeutils.EuropeStockholm())

	t.Logf("PAGINATION SCENARIO:")
	t.Logf("Last departure from previous page: %s UTC", lastDepartureTime.Format("2006-01-02 15:04:05"))
	t.Logf("Stockholm time: %s", stockholmTime.Format("2006-01-02 15:04:05"))
	t.Logf("User wants: Trips departing AFTER %s", stockholmTime.Format("15:04"))
	t.Logf("")
	t.Logf("SDK generates these Trafiklab API parameters:")
	t.Logf("  itd_date=%s", params.Get("itd_date"))
	t.Logf("  itd_time=%s", params.Get("itd_time"))
	t.Logf("  itd_trip_date_time_dep_arr=%s", params.Get("itd_trip_date_time_dep_arr"))
	t.Logf("")
	t.Logf("ISSUE: If Trafiklab API interprets 'dep' + time as 'departures AT this time'")
	t.Logf("instead of 'departures AFTER this time', then it will return the same results")
	t.Logf("or results that start before %s, breaking pagination.", stockholmTime.Format("15:04"))
	t.Logf("")
	t.Logf("Full query: %s", params.Encode())

	// The test passes regardless - this is for documentation/debugging
	// The real test is whether the Trafiklab API behaves correctly with these parameters
}

func TestTrips_UnixTimestampConversion(t *testing.T) {
	// Test that Unix timestamps from your API are correctly converted

	// Using the exact timestamp from the user's curl example
	unixTimestamp := int64(1750089240)
	inputTime := time.Unix(unixTimestamp, 0)

	request := &TripsRequest{
		At:                 inputTime,
		NumTrips:           3,
		TypeOrigin:         "any",
		NameOrigin:         "9091001000001002",
		TypeDestination:    "any",
		NameDestination:    "9091001001002800",
		TripDateTimeDepArr: "dep",
	}

	params := request.toParams()

	// Convert to Stockholm time
	stockholmTime := inputTime.In(timeutils.EuropeStockholm())
	expectedDate := stockholmTime.Format("20060102")
	expectedTime := stockholmTime.Format("1504")

	if params.Get("itd_date") != expectedDate {
		t.Errorf("Date conversion failed: expected %s, got %s",
			expectedDate, params.Get("itd_date"))
	}

	if params.Get("itd_time") != expectedTime {
		t.Errorf("Time conversion failed: expected %s, got %s",
			expectedTime, params.Get("itd_time"))
	}

	t.Logf("Unix timestamp %d converted to:", unixTimestamp)
	t.Logf("  UTC: %s", inputTime.Format("2006-01-02 15:04:05"))
	t.Logf("  Stockholm: %s", stockholmTime.Format("2006-01-02 15:04:05"))
	t.Logf("  API params: itd_date=%s, itd_time=%s",
		params.Get("itd_date"), params.Get("itd_time"))
}

func TestTrips_SemanticsQuestion(t *testing.T) {
	// This test documents the core question about API semantics

	t.Log("SEMANTIC QUESTION: What does 'dep' + specific time mean in Trafiklab API?")
	t.Log("")
	t.Log("Option A: 'dep' means departures AT OR AFTER the specified time")
	t.Log("  - This would work correctly for pagination")
	t.Log("  - depart_after=17:54 → returns trips departing from 17:54 onwards")
	t.Log("")
	t.Log("Option B: 'dep' means departures AROUND the specified time")
	t.Log("  - This would break pagination (current behavior)")
	t.Log("  - depart_after=17:54 → returns trips departing around 17:54 (including before)")
	t.Log("")
	t.Log("Option C: 'dep' means departures AT the specified time")
	t.Log("  - This would also break pagination")
	t.Log("  - depart_after=17:54 → returns trips departing exactly at 17:54")
	t.Log("")
	t.Log("HYPOTHESIS: The issue might be that we need to add some offset to the time")
	t.Log("when using 'dep' to get true 'after' behavior, or there's a different parameter")
	t.Log("or flag that controls this behavior.")
}

func TestTrips_CalcOneDirectionFix(t *testing.T) {
	// This test demonstrates the solution for the pagination issue
	// The calc_one_direction flag prevents getting trips before the specified time

	lastDepartureTime := time.Unix(1750089240, 0)

	t.Log("SOLUTION: Use calc_one_direction flag for pagination")
	t.Log("")
	t.Log("Problem: When requesting trips with depart_after, Trafiklab API returns")
	t.Log("trips that depart BEFORE the specified time (20-30 minutes earlier)")
	t.Log("")
	t.Log("Solution: Set calc_one_direction=true to prevent trips before the time")
	t.Log("")

	// Test current behavior (without calc_one_direction)
	currentRequest := &TripsRequest{
		At:                 lastDepartureTime,
		NumTrips:           3,
		TypeOrigin:         "any",
		NameOrigin:         "9091001000001002",
		TypeDestination:    "any",
		NameDestination:    "9091001001002800",
		TripDateTimeDepArr: "dep",
		Flags:              []string{}, // No calc_one_direction
	}

	currentParams := currentRequest.toParams()

	// Test fixed behavior (with calc_one_direction)
	fixedRequest := &TripsRequest{
		At:                 lastDepartureTime,
		NumTrips:           3,
		TypeOrigin:         "any",
		NameOrigin:         "9091001000001002",
		TypeDestination:    "any",
		NameDestination:    "9091001001002800",
		TripDateTimeDepArr: "dep",
		Flags:              []string{"calc_one_direction"}, // ADD THIS FLAG
	}

	fixedParams := fixedRequest.toParams()

	stockholmTime := lastDepartureTime.In(timeutils.EuropeStockholm())

	t.Logf("Current (broken) request:")
	t.Logf("  Query: %s", currentParams.Encode())
	t.Logf("  Result: Returns trips departing BEFORE %s (causing pagination issues)", stockholmTime.Format("15:04"))
	t.Logf("")
	t.Logf("Fixed request:")
	t.Logf("  Query: %s", fixedParams.Encode())
	t.Logf("  Result: Should only return trips departing AT OR AFTER %s", stockholmTime.Format("15:04"))

	// Verify the flag is set
	if fixedParams.Get("calc_one_direction") != "true" {
		t.Errorf("Expected calc_one_direction=true, got '%s'", fixedParams.Get("calc_one_direction"))
	}

	// Verify the flag is NOT set in current request
	if currentParams.Get("calc_one_direction") != "" {
		t.Errorf("Expected calc_one_direction to be empty in current request, got '%s'", currentParams.Get("calc_one_direction"))
	}
}