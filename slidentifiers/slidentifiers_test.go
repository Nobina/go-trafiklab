package slidentifiers_test

import (
	"testing"

	"github.com/nobina/go-trafiklab/slidentifiers"
)

func TestConvertHafasToEFA(t *testing.T) {
	prefix := "909100100"
	tests := []struct {
		hafasID string
		want    string
		wantErr bool
	}{
		{hafasID: "300104400", want: "9091001000004400", wantErr: false},
		{hafasID: "300109192", want: "9091001000009192", wantErr: false},
		{hafasID: "300109669", want: "9091001000009669", wantErr: false},
		{hafasID: "123", want: "", wantErr: true},                    // too short
		{hafasID: "400109192", want: "", wantErr: true},              // doesn't start with '3'
	}

	for _, tt := range tests {
		got, err := slidentifiers.ConvertHafasToEFA(tt.hafasID, prefix)
		if (err != nil) != tt.wantErr {
			t.Errorf("ConvertHafasToEFA(%q, %q) error = %v, wantErr %v", tt.hafasID, prefix, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("ConvertHafasToEFA(%q, %q) = %q, want %q", tt.hafasID, prefix, got, tt.want)
		}
	}
}

func TestConvertIDToHafas(t *testing.T) {
	tests := []struct {
		sid     string
		want    string
		wantErr bool
	}{
		{sid: "4400", want: "300104400", wantErr: false},
		{sid: "9192", want: "300109192", wantErr: false},
		{sid: "9669", want: "300109669", wantErr: false},
		{sid: "notnum", want: "", wantErr: true},
		{sid: "12345678", want: "12345678", wantErr: false}, // length > 7 returns unchanged
	}

	for _, tt := range tests {
		got, err := slidentifiers.ConvertIDToHafas(tt.sid)
		if (err != nil) != tt.wantErr {
			t.Errorf("convertIDToHafas(%q) error = %v, wantErr %v", tt.sid, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("convertIDToHafas(%q) = %q, want %q", tt.sid, got, tt.want)
		}
	}
}