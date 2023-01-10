package constants

import "time"

// ExpiryDelta describes the puffer-time for a token update before it will expire.
const ExpiryDelta = 1 * time.Minute
const MaxRetryCount = 10
const RetryDelay = 1 * time.Second
