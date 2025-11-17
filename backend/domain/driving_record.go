package domain

import (
	"time"

	"github.com/aereal/optional"
)

type DrivingRecord struct {
	Date               time.Time
	DistanceKilometers int64
	Memo               optional.Option[string]
}
