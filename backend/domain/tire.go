package domain

import (
	"time"

	"github.com/aereal/optional"
)

type TireID string

type Tire struct {
	ID          TireID
	Name        string
	PurchasedOn optional.Option[time.Time]
	Retired     bool
}
