package utils

import (
	"fmt"
	"time"
)

// ParseMMYYYY parses a date string in MM-YYYY format to time.Time
func ParseMMYYYY(dateStr string) (time.Time, error) {
	// Parse the date string in MM-YYYY format
	t, err := time.Parse("01-2006", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, expected MM-YYYY: %w", err)
	}
	return t, nil
}

// FormatMMYYYY formats a time.Time to MM-YYYY string
func FormatMMYYYY(t time.Time) string {
	return t.Format("01-2006")
}

// IsActiveDuringPeriod checks if a subscription was active during a given period
func IsActiveDuringPeriod(start, end, periodStart, periodEnd time.Time) bool {
	// If subscription has no end date, it's active until further notice
	if end.IsZero() {
		return !start.After(periodEnd)
	}
	
	// Check if the periods overlap
	return !start.After(periodEnd) && !end.Before(periodStart)
}