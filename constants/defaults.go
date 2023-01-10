package constants

import "time"

// ExpiryDelta describes the puffer-time for a token update before it will expire.
const ExpiryDelta = 1 * time.Minute

// MaxRetryCount describes the maximum number of retries for a request.
const MaxRetryCount = 10

// StartingRetryDelay describes the first delay between retries.
// retries grow exponentially 2^attempts * StartingRetryDelay.
// Where attempts starts internally at 0.
const StartingRetryDelay = 1 * time.Second
