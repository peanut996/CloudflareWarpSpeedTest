package utils

import "encoding/json"

func ParseReservedString(reservedString string) (reserved [3]byte, err error) {
	if reservedString == "" {
		return
	}
	reserved = [3]byte{}
	err = json.Unmarshal([]byte(reservedString), &reserved)
	return
}
