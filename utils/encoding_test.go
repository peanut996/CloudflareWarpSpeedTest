package utils

import (
	"reflect"
	"testing"
)

func TestParseReservedString(t *testing.T) {
	tests := []struct {
		name           string
		reservedString string
		wantReserved   [3]byte
		wantErr        bool
	}{
		{
			name:           "empty string",
			reservedString: "",
			wantReserved:   [3]byte{},
			wantErr:        false,
		},
		{
			name:           "valid array",
			reservedString: "[1, 2, 3]",
			wantReserved:   [3]byte{1, 2, 3},
			wantErr:        false,
		},
		{
			name:           "invalid json",
			reservedString: "invalid",
			wantReserved:   [3]byte{},
			wantErr:        true,
		},
		{
			name:           "wrong length",
			reservedString: "[1, 2, 3, 4]",
			wantReserved:   [3]byte{},
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReserved, err := ParseReservedString(tt.reservedString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReservedString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(gotReserved, tt.wantReserved) {
				t.Errorf("ParseReservedString() = %v, want %v", gotReserved, tt.wantReserved)
			}
		})
	}
}
