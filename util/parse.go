package util

import (
	"fmt"
	"strconv"
	"time"
)

// ParseFloat64 converts a string to float64, returning 0 if parsing fails
func ParseFloat64(s string) float64 {
	if s == "" {
		return 0
	}
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Printf("Error parsing '%s' as float: %v\n", s, err)
		return 0
	}
	return value
}

// ParseFloat32 converts a string to float32, returning 0 if parsing fails
func ParseFloat32(s string) float32 {
	if s == "" {
		return 0
	}
	value, err := strconv.ParseFloat(s, 32)
	if err != nil {
		fmt.Printf("Error parsing '%s' as float: %v\n", s, err)
		return 0
	}
	return float32(value)
}

// FormatTimestampToISO converts a unix timestamp string to ISO 8601 format
// Returns empty string if parsing fails
func FormatTimestampToISO(timestamp string) string {
	if timestamp == "" {
		return ""
	}

	// Try to parse as int64 (seconds)
	sec, err := strconv.ParseInt(timestamp, 10, 64)
	if err == nil {
		// Check if it's in milliseconds (13 digits) instead of seconds (10 digits)
		if len(timestamp) > 10 {
			return time.UnixMilli(sec).Format(time.RFC3339)
		}
		return time.Unix(sec, 0).Format(time.RFC3339)
	}

	// Try to parse as float64 for timestamps with decimals
	secFloat, err := strconv.ParseFloat(timestamp, 64)
	if err != nil {
		fmt.Printf("Error parsing timestamp '%s': %v\n", timestamp, err)
		return ""
	}

	sec = int64(secFloat)
	nsec := int64((secFloat - float64(sec)) * 1e9)
	return time.Unix(sec, nsec).Format(time.RFC3339)
}
