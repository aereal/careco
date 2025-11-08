package resolver

import (
	"iter"
	"time"

	"careco/backend/graph/dtos"
)

func generateDummyData(base time.Time) iter.Seq[*dtos.DailyReport] {
	return func(yield func(*dtos.DailyReport) bool) {
		var i int
		for {
			r := &dtos.DailyReport{
				RecordedAt:         base.Add(time.Hour * -24 * time.Duration(i)),
				DistanceKilometers: 12 * i,
			}
			if i%2 == 1 {
				memo := "blah blah"
				r.Memo = &memo
			}
			if !yield(r) {
				return
			}
			i++
		}
	}
}
