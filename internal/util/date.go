package util

import "time"

// Today returns the current date in YYYY-MM-DD format.
func Today() string {
	return time.Now().Format("2006-01-02")
}

// IsValidDate returns true if the date is valid.
func IsValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

// DiffDays returns the difference in days from the given date to today.
func DiffDays(date string) int {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0
	}
	return int(time.Since(t).Hours() / 24)
}
