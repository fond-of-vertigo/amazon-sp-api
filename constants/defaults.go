package constants

import "time"

// ExpiryDelta describes the puffer-time for a token update before it will expire.
const ExpiryDelta = 1 * time.Minute

// MaxRetryCountOnTooManyRequestsError is the maximum retry if we get HTTP 429 error
const MaxRetryCountOnTooManyRequestsError = 20

// DefaultWaitDurationOnTooManyRequestsError is the default wait time between two requests
// on HTTP 429 error
const DefaultWaitDurationOnTooManyRequestsError = 1 * time.Second
