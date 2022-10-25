package apis

import (
	"fmt"
	"time"
)

// JsonTimeISO8601 marshals time.Time to ISO 8601 format
type JsonTimeISO8601 time.Time

func (t JsonTimeISO8601) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t)), nil
}

func (t JsonTimeISO8601) String() string {
	return time.Time(t).Format(time.RFC3339)
}
