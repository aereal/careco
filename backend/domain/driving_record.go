package domain

import (
	"time"

	"github.com/aereal/optional"
)

type DrivingRecord struct {
	Date time.Time
	Memo optional.Option[string]
	// OdometerValue はこの日時時点での総走行距離。
	//
	// 日時は [Date] で表される。
	OdometerValue int64
}
