package constants

import "time"

const (
	// ExpiryDelta describes the puffer-time for a token update before it will expire
	ExpiryDelta time.Duration = 1 * time.Minute

	// MaxRetryCountOnTooManyRequestsError is the maximum retry if we get HTTP 429 error
	MaxRetryCountOnTooManyRequestsError int = 20
	// DefaultWaitDurationOnTooManyRequestsError is the default wait time between two requests
	// on HTTP 429 error
	DefaultWaitDurationOnTooManyRequestsError time.Duration = 1 * time.Second

	//DefaultTokenUpdaterBackoffTime is the default backoff time for the token updater when a request fails
	DefaultTokenUpdaterBackoffTime time.Duration = 15 * time.Second
)
