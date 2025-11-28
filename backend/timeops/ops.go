package timeops

import "time"

func StartOfNextMonth(t time.Time) time.Time {
	if t.Month() == time.December {
		return startOfMonth(t.Year()+1, time.January, t.Location())
	}
	return startOfMonth(t.Year(), t.Month()+1, t.Location())
}

func StartOfNextYear(t time.Time) time.Time {
	return startOfMonth(t.Year()+1, time.January, t.Location())
}

func startOfMonth(year int, month time.Month, loc *time.Location) time.Time {
	return time.Date(year, month, 1, 0, 0, 0, 0, loc)
}
