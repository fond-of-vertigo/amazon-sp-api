package apis

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/fond-of-vertigo/amazon-sp-api/constants"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
	GetEndpoint() constants.Endpoint
	Close()
}
type CallResponse[responseBodyType any] struct {
	Status       int
	ResponseBody *responseBodyType
	ErrorList    *ErrorList
}
type Call[responseType any] struct {
	Method                  string
	URL                     string
	QueryParams             url.Values
	Body                    []byte
	RestrictedDataToken     *string
	ParseErrorListOnError   bool
	WaitDurationOnRateLimit time.Duration
}

func NewCall[responseType any](method string, url string) *Call[responseType] {
	return &Call[responseType]{
		Method:                  method,
		URL:                     url,
		WaitDurationOnRateLimit: constants.DefaultWaitDurationOnTooManyRequestsError,
	}
}

// sleeper func as type for mocking
type sleeper func(d time.Duration)

var sleepFunc sleeper = time.Sleep

func (a *Call[responseType]) WithQueryParams(queryParams url.Values) *Call[responseType] {
	a.QueryParams = queryParams
	return a
}

func (a *Call[responseType]) WithBody(body []byte) *Call[responseType] {
	a.Body = body
	return a
}

// WithRestrictedDataToken is optional and can be passed to replace the existing accessToken
func (a *Call[responseType]) WithRestrictedDataToken(token *string) *Call[responseType] {
	a.RestrictedDataToken = token
	return a
}

func (a *Call[responseType]) WithParseErrorListOnError() *Call[responseType] {
	a.ParseErrorListOnError = true
	return a
}

func (a *Call[responseType]) WithRateLimit(callsPer float32, duration time.Duration) *Call[responseType] {
	a.WaitDurationOnRateLimit = calcWaitTimeByRateLimit(callsPer, duration)
	return a
}

// Execute will return response object on success
func (a *Call[responseType]) Execute(httpClient HTTPClient) (*CallResponse[responseType], error) {
	resp, err := a.execute(httpClient)
	if err != nil {
		return nil, err
	}

	callResp := &CallResponse[responseType]{
		Status: resp.StatusCode,
	}

	if callResp.IsSuccess() { // Response is OK
		if resp.ContentLength > 0 {
			if err = unmarshalBody(resp, &callResp.ResponseBody); err != nil {
				return nil, err
			}
		}
		return callResp, nil
	}

	// Response is not OK
	err = fmt.Errorf("request with URL=%v returned with non-OK statuscode=%d", a.URL, callResp.Status)
	if a.ParseErrorListOnError { // ErrorList is activated, try to parse it
		if parseErr := unmarshalBody(resp, &callResp.ErrorList); parseErr != nil {
			return nil, errors.Join(err, parseErr)
		}

		return callResp, errors.Join(err, mapErrorListToError(callResp.ErrorList))
	}

	return callResp, err
}

func (a *Call[responseType]) execute(httpClient HTTPClient) (*http.Response, error) {
	for attempts := 0; attempts < constants.MaxRetryCountOnTooManyRequestsError; attempts++ {
		req, err := a.createNewRequest(httpClient.GetEndpoint())
		if err != nil {
			return nil, err
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			sleepFunc(a.WaitDurationOnRateLimit)
			continue
		}

		return resp, nil
	}

	return nil, ErrMaxRetryCountReached
}

func (a *Call[responseType]) createNewRequest(endpoint constants.Endpoint) (*http.Request, error) {
	callURL, err := url.Parse(string(endpoint) + a.URL)
	if err != nil {
		return nil, err
	}
	callURL.RawQuery = a.QueryParams.Encode()

	req, err := http.NewRequest(a.Method, callURL.String(), bytes.NewBuffer(a.Body))
	if err == nil {
		if a.RestrictedDataToken != nil && *a.RestrictedDataToken != "" {
			req.Header.Add(constants.AccessTokenHeader, *a.RestrictedDataToken)
		}
	}
	return req, err
}

// IsSuccess checks if the status is in range 2xx
func (r *CallResponse[any]) IsSuccess() bool {
	return r.Status >= http.StatusOK && r.Status < http.StatusMultipleChoices
}

// IsError checks if the status is in range 4xx - 5xx
func (r *CallResponse[any]) IsError() bool {
	return r.Status >= http.StatusBadRequest && r.Status < 600
}

// ErrorsAsString returns all errors as string or an empty string.
func (r *CallResponse[any]) ErrorsAsString() string {
	if r == nil || !r.IsError() {
		return ""
	}

	msg := fmt.Sprintf("received HTTP status code %d", r.Status)
	if r.ErrorList != nil {
		msg = fmt.Sprintf("%s\n%v", msg, r.ErrorList)
	}
	return msg
}

func calcWaitTimeByRateLimit(callsPer float32, duration time.Duration) time.Duration {
	return time.Duration(float32(duration) / callsPer)
}
