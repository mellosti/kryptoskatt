package util

import (
	"fmt"
	"strconv"
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
