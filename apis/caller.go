package apis

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/fond-of-vertigo/amazon-sp-api/constants"
	"github.com/fond-of-vertigo/amazon-sp-api/httpx"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"net/url"
)

type CallResponse[responseBodyType any] struct {
	Status       int
	ResponseBody *responseBodyType
	ErrorList    *ErrorList
}
type Call[responseType any] interface {
	WithQueryParams(url.Values) Call[responseType]
	WithBody([]byte) Call[responseType]
	// WithRestrictedDataToken is optional and can be passed to replace the existing accessToken
	WithRestrictedDataToken(*string) Call[responseType]
	WithParseErrorListOnError(bool) Call[responseType]
	WithRateLimiter(rateLimiter *rate.Limiter) Call[responseType]
	// Execute will return response object on success
	Execute(httpClient httpx.Client) (*CallResponse[responseType], error)
}

func NewCall[responseType any](method string, url string) Call[responseType] {
	return &call[responseType]{
		Method: method,
		URL:    url,
	}
}

type call[responseType any] struct {
	Method                string
	URL                   string
	QueryParams           url.Values
	Body                  []byte
	RestrictedDataToken   *string
	ParseErrorListOnError bool
	RateLimiter           *rate.Limiter
}

func (a *call[responseType]) WithQueryParams(queryParams url.Values) Call[responseType] {
	a.QueryParams = queryParams
	return a
}
func (a *call[responseType]) WithBody(body []byte) Call[responseType] {
	a.Body = body
	return a
}

func (a *call[responseType]) WithRestrictedDataToken(token *string) Call[responseType] {
	a.RestrictedDataToken = token
	return a
}

func (a *call[responseType]) WithParseErrorListOnError(parseErrList bool) Call[responseType] {
	a.ParseErrorListOnError = parseErrList
	return a
}
func (a *call[responseType]) WithRateLimiter(rateLimiter *rate.Limiter) Call[responseType] {
	a.RateLimiter = rateLimiter
	return a
}
func (a *call[responseType]) Execute(httpClient httpx.Client) (*CallResponse[responseType], error) {
	resp, err := a.execute(httpClient)

	if err != nil {
		return nil, err
	}

	callResp := &CallResponse[responseType]{
		Status: resp.StatusCode,
	}

	if a.ParseErrorListOnError && !callResp.IsSuccess() {
		if err = unmarshalBody(resp, &callResp.ErrorList); err != nil {
			return nil, err
		}
		return callResp, nil
	}
	if resp.ContentLength > 0 {
		if err = unmarshalBody(resp, &callResp.ResponseBody); err != nil {
			return nil, err
		}
	}

	return callResp, nil
}

func (a *call[responseType]) execute(httpClient httpx.Client) (*http.Response, error) {

	req, err := a.createNewRequest(httpClient.GetEndpoint())
	if err != nil {
		return nil, err
	}
	if a.RateLimiter != nil {
		err = a.RateLimiter.Wait(context.Background())
		if err != nil {
			return nil, err
		}
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (a *call[responseType]) createNewRequest(endpoint constants.Endpoint) (*http.Request, error) {
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
func unmarshalBody(resp *http.Response, into any) error {
	if resp.ContentLength == 0 {
		return nil
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(bodyBytes, into)
}
