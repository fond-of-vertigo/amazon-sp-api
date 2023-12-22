package apis

import (
	"time"
)

// JsonTimeISO8601 marshals time.Time to ISO 8601 format
type JsonTimeISO8601 struct {
	time.Time
}

func (t JsonTimeISO8601) MarshalJSON() ([]byte, error) {
	value := "\"" + t.UTC().Format(time.RFC3339) + "\""
	return []byte(value), nil
}

func (t JsonTimeISO8601) String() string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}
