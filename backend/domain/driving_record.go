package domain

import (
	"time"

	"github.com/aereal/optional"
)

type DrivingRecord struct {
	Date               time.Time
	CumulativeDistance int64
	Memo               optional.Option[string]
}
