package attr

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
)

func ISO8601Value(t time.Time) attribute.Value {
	return attribute.StringValue(t.Format(time.RFC3339Nano))
}

func DurationSecondsValue(d time.Duration) attribute.Value {
	return attribute.Float64Value(d.Seconds())
}
