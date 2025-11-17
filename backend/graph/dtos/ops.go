package dtos

import "time"

func (mr *MonthlyReport) StartOfMonth() time.Time {
	return time.Date(mr.Year, mr.Month, 1, 0, 0, 0, 0, time.Local)
}

func (yr *YearlyReport) StartOfYear() time.Time {
	return time.Date(yr.Year, time.January, 1, 0, 0, 0, 0, time.Local)
}
