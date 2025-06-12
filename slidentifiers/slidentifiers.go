// Package slidentifiers provides utilities for converting between different SL (Storstockholms Lokaltrafik)
// identification systems used in their journey planning APIs.
//
// Background:
// SL updated their journey planning system from HAFAS to EFA, which broke their ID system.
// This package provides conversion functions to help migrate saved favorites and other stored IDs
// between the old HAFAS format and the new EFA format.
//
// HAFAS ID Format (9 digits): XFGYEDCBA
//   - X: Always 3 for Sites (position 0)
//   - F,G,E,D,C,B,A: The 7-digit site number from positions 1,2,4,5,6,7,8
//   - Y: TransportAuthority.Number modulus 10 (position 3) - always 1 for SL
//
// EFA ID Format (16 digits): [9-digit prefix][7-digit site number]
//   - Prefix: "909100100" for SL Stockholm region
//   - Site number: 7-digit number extracted from HAFAS ID
//
// About the EFA prefix:
// The prefix "909100100" is defined in SL's official Pubtrans GID (Global ID) specification.
// According to SL's documentation, the prefix corresponds to the "Place (site)" entity type
// in their standardized ID scheme for the Stockholm region.
//
// IMPORTANT: This prefix may not apply to all stops. Different entity types in the Pubtrans
// specification (Stop Area, Stop Point, Station Entrance Point, etc.) or different transport
// authorities/regions may use different prefixes. The "909100100" prefix has been observed
// for standard SL Stockholm region stops, but edge cases may exist for:
// - Different transport authorities (e.g., Waxholmsbolaget)
// - Different entity types (Stop Area vs Stop Point vs Place)
// - Different regions or municipalities
// - Regional transport operators
//
// For robust conversion, consider making the prefix configurable or detecting it dynamically
// from actual EFA API responses when possible.
//
// Example conversion:
//
//	HAFAS: "300104400" (Stavsnäs vinterhamn)
//	EFA:   "9091001000004400"
package slidentifiers

import (
	"fmt"
	"strconv"
	"unicode"
)

// ConvertHafasToEFA converts a HAFAS site ID to an EFA Global ID (GID).
//
// This function extracts the 7-digit site number from a 9-digit HAFAS ID and combines it
// with a 9-digit EFA prefix to create a 16-digit EFA GID.
//
// HAFAS ID structure (XFGYEDCBA):
//   Position: 0 1 2 3 4 5 6 7 8
//   Meaning:  X F G Y E D C B A
//   Where:
//   - X = 3 (always, indicates Site)
//   - Y = 1 (SL transport authority number mod 10)
//   - GFEDCBA = 7-digit site number (extracted from positions 2,1,4,5,6,7,8)
//
// Parameters:
//   - hafasID: 9-digit HAFAS ID (must start with '3' and be all digits)
//   - prefix: 9-digit EFA prefix (typically "909100100" for SL Stockholm)
//
// Returns:
//   - 16-digit EFA GID string
//   - error if validation fails
//
// Example:
//   hafasID := "300104400"  // Stavsnäs vinterhamn
//   prefix := "909100100"   // SL Stockholm prefix
//   efaID, err := ConvertHafasToEFA(hafasID, prefix)
//   // Result: "9091001000004400"
func ConvertHafasToEFA(hafasID string, prefix string) (string, error) {
    // 1) Validate HAFAS ID format
    if len(hafasID) != 9 {
        return "", fmt.Errorf("invalid HAFAS ID %q: must be exactly 9 digits", hafasID)
    }
    if hafasID[0] != '3' {
        return "", fmt.Errorf("invalid HAFAS ID %q: must start with '3'", hafasID)
    }
    for i, r := range hafasID {
        if !unicode.IsDigit(r) {
            return "", fmt.Errorf("invalid HAFAS ID %q: character %d (%q) is not a digit", hafasID, i, r)
        }
    }

    // 2) Extract the 7-digit site number from HAFAS ID
    //    The HAFAS format XFGYEDCBA stores the site number as GFEDCBA
    //    We need to reorder: positions 2,1,4,5,6,7,8 → GFEDCBA
    //
    //    Example: "300104400"
    //             X=3, F=0, G=0, Y=1, E=0, D=4, C=4, B=0, A=0
    //             Site number: G(0) + F(0) + E(0) + D(4) + C(4) + B(0) + A(0) = "0004400"
    extractedSiteNumber := []byte{
        hafasID[2], // G (position 2)
        hafasID[1], // F (position 1)
        hafasID[4], // E (position 4)
        hafasID[5], // D (position 5)
        hafasID[6], // C (position 6)
        hafasID[7], // B (position 7)
        hafasID[8], // A (position 8)
    }

    // Verify the extracted site number is numeric
    if _, err := strconv.Atoi(string(extractedSiteNumber)); err != nil {
        return "", fmt.Errorf("extracted site number %q is not numeric", string(extractedSiteNumber))
    }

    // 3) Validate EFA prefix format
    if len(prefix) != 9 {
        return "", fmt.Errorf("invalid EFA prefix %q: must be exactly 9 digits", prefix)
    }
    for i, r := range prefix {
        if !unicode.IsDigit(r) {
            return "", fmt.Errorf("invalid EFA prefix %q: character %d (%q) is not a digit", prefix, i, r)
        }
    }

    // 4) Combine prefix and site number to create 16-digit EFA GID
    return prefix + string(extractedSiteNumber), nil
}

// ConvertIDToHafas converts a legacy SL site ID to HAFAS format.
//
// This is a temporary conversion function needed after SL's domain update broke their ID system.
// It handles the conversion from old short site IDs to the 9-digit HAFAS format.
//
// Conversion logic for IDs ≤ 7 digits:
//   1. Split the ID into first 2 digits and last 5 digits
//   2. Format as: "3" + first2digits + "1" (SL authority) + last5digits
//   3. The middle "1" represents SL's transport authority number
//
// Parameters:
//   - sid: Site ID string (if > 7 digits, returned unchanged)
//
// Returns:
//   - HAFAS ID string (9 digits if converted, original if > 7 digits)
//   - error if the input cannot be converted to integer
//
// Examples:
//   "4400" → "300104400"     (00 + 04400 with SL authority "1")
//   "9192" → "300109192"     (00 + 09192 with SL authority "1")
//   "12345678" → "12345678"  (unchanged, already > 7 digits)
//
// Note: This function will be deprecated once SL's other major breaking changes are resolved.
func ConvertIDToHafas(sid string) (string, error) {
	// If ID is already long enough (> 7 digits), assume it's already in correct format
	if len(sid) > 7 {
		return sid, nil
	}

	// Convert string to integer for manipulation
	id, err := strconv.Atoi(sid)
	if err != nil {
		return "", fmt.Errorf("failed to convert id to hafas: %w", err)
	}

	// Split the ID: first 2 digits and last 5 digits
	// For IDs shorter than 7 digits, this will pad with zeros appropriately
	firstTwoDigits := id / 100000  // Integer division gets first 2 digits (or 0 if < 100000)
	lastFiveDigits := id % 100000  // Modulo gets last 5 digits

	// Build HAFAS ID: "3" + first2digits + "1" (SL authority) + last5digits
	// Format ensures proper zero-padding: %02d for 2 digits, %05d for 5 digits
	return fmt.Sprintf("3%02d1%05d", firstTwoDigits, lastFiveDigits), nil
}