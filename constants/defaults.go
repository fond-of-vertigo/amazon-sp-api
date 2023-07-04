package constants

import "time"

// ExpiryDelta describes the puffer-time for a token update before it will expire
const ExpiryDelta time.Duration = 1 * time.Minute

// MaxRetryCountOnTooManyRequestsError is the maximum retry if we get HTTP 429 error
const MaxRetryCountOnTooManyRequestsError int = 20

// DefaultWaitDurationOnTooManyRequestsError is the default wait time between two requests
// on HTTP 429 error
const DefaultWaitDurationOnTooManyRequestsError time.Duration = 1 * time.Second
