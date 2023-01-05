package apis

import (
	"time"
)

// JsonTimeISO8601 marshals time.Time to ISO 8601 format
type JsonTimeISO8601 struct {
	time.Time
}

func (t JsonTimeISO8601) MarshalJSON() ([]byte, error) {
	value := "\"" + t.Format(time.RFC3339Nano) + "\""
	return []byte(value), nil
}

func (t JsonTimeISO8601) String() string {
	return t.Format(time.RFC3339Nano)
}
