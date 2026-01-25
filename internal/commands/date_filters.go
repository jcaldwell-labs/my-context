package commands

import (
	"fmt"
	"time"
)

// TimePeriod represents the time range for filtering
type TimePeriod struct {
	Start time.Time
	End   time.Time
	Label string
}

// calculatePeriod determines the time period based on flags
func calculatePeriod(today, week, month bool, since, until string) (TimePeriod, error) {
	now := time.Now()
	var period TimePeriod

	switch {
	case today:
		// Start of today in local time
		y, m, d := now.Date()
		period.Start = time.Date(y, m, d, 0, 0, 0, 0, now.Location())
		period.End = now
		period.Label = "Today"

	case week:
		// Start of week (Monday)
		y, m, d := now.Date()
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday = 7
		}
		daysBack := weekday - 1 // Days since Monday
		period.Start = time.Date(y, m, d-daysBack, 0, 0, 0, 0, now.Location())
		period.End = now
		period.Label = "This Week"

	case month:
		// Start of month
		y, m, _ := now.Date()
		period.Start = time.Date(y, m, 1, 0, 0, 0, 0, now.Location())
		period.End = now
		period.Label = "This Month"

	case since != "" || until != "":
		// Custom date range
		if since != "" {
			t, err := time.Parse("2006-01-02", since)
			if err != nil {
				return period, fmt.Errorf("invalid --since date format (use YYYY-MM-DD): %w", err)
			}
			period.Start = t
		} else {
			// Default to beginning of time (1970)
			period.Start = time.Unix(0, 0)
		}

		if until != "" {
			t, err := time.Parse("2006-01-02", until)
			if err != nil {
				return period, fmt.Errorf("invalid --until date format (use YYYY-MM-DD): %w", err)
			}
			// End of the specified day
			period.End = t.Add(24*time.Hour - time.Second)
		} else {
			period.End = now
		}

		switch {
		case since != "" && until != "":
			period.Label = fmt.Sprintf("%s to %s", since, until)
		case since != "":
			period.Label = fmt.Sprintf("Since %s", since)
		default:
			period.Label = fmt.Sprintf("Until %s", until)
		}

	default:
		// All time
		period.Start = time.Unix(0, 0)
		period.End = now
		period.Label = "All Time"
	}

	return period, nil
}
