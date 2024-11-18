package utils

import (
	"encoding/json"
	"fmt"
)

func ParseReservedString(reservedString string) (reserved [3]byte, err error) {
	if reservedString == "" {
		return
	}
	
	// First unmarshal into a slice to validate length
	var tempSlice []byte
	if err = json.Unmarshal([]byte(reservedString), &tempSlice); err != nil {
		return
	}
	
	// Validate length
	if len(tempSlice) != 3 {
		err = fmt.Errorf("reserved array must have exactly 3 elements, got %d", len(tempSlice))
		return
	}
	
	// Copy to fixed-size array
	reserved = [3]byte{tempSlice[0], tempSlice[1], tempSlice[2]}
	return
}
