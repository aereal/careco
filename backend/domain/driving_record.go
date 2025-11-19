package domain

import (
	"time"

	"github.com/aereal/optional"
)

type DrivingRecordToRecord struct {
	Date               time.Time
	DistanceKilometers int64
	Memo               optional.Option[string]
}

type DrivingRecord struct {
	Date                    time.Time
	DistanceKilometers      int64
	Memo                    optional.Option[string]
	TotalDistanceKilometers int64
}
