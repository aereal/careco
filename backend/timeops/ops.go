package timeops

import "time"

func StartOfNextMonth(t time.Time) time.Time {
	if t.Month() == time.December {
		return time.Date(t.Year()+1, time.January, 1, 0, 0, 0, 0, t.Location())
	}
	return time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
}

func StartOfNextYear(t time.Time) time.Time {
	return time.Date(t.Year()+1, time.January, 1, 0, 0, 0, 0, t.Location())
}
