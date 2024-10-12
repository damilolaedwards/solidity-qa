package utils

import (
	"fmt"
	"strings"
)

// MapToDictString converts a map to a dictionary-like string
func MapToDictString(inputMap map[string]any) string {
	var sb strings.Builder
	sb.WriteString("{")

	// Counter to help manage commas between key-value pairs
	counter := 0
	totalItems := len(inputMap)

	// Iterate through the map
	for key, value := range inputMap {
		sb.WriteString(fmt.Sprintf("\"%s\": ", key))

		// Handle different types of values
		switch v := value.(type) {
		case string:
			sb.WriteString(fmt.Sprintf("\"%s\"", v))
		default:
			sb.WriteString(fmt.Sprintf("%v", v))
		}

		// Add a comma between items, except for the last one
		if counter < totalItems-1 {
			sb.WriteString(", ")
		}
		counter++
	}

	sb.WriteString("}")
	return sb.String()
}

// SliceContains returns whether
func SliceContains(arr []string, target string) bool {
	m := make(map[string]bool)
	for _, s := range arr {
		m[s] = true
	}
	return m[target]
}
