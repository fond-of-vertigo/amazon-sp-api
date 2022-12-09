package apis

// CallError represents one error occurred while calling the amazon api.
// If the api itself returns a ErrorList as response body calling ErrorList() will provide it.
type CallError interface {
	error
	ErrorList() ErrorList
}

func NewError(err error) CallError {
	return &callError{
		err: err,
	}
}

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

type callError struct {
	err       error
	errorList ErrorList
}

func (a *callError) Error() string {
	return a.err.Error()
}

func (a *callError) ErrorList() ErrorList {
	return a.errorList
}
