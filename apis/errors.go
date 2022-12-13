package apis

// Error response returned when the request is unsuccessful.
type Error struct {
	// An error code that identifies the type of error that occurred.
	Code string `json:"code"`
	// A message that describes the error condition in a human-readable form.
	Message string `json:"message"`
	// Additional details that can help the caller understand or fix the issue.
	Details *string `json:"details,omitempty"`
}

// ErrorList A list of error responses returned when a request is unsuccessful.
type ErrorList struct {
	Errors []Error `json:"errors"`
}
