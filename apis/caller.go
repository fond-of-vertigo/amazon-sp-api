package apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/fond-of-vertigo/amazon-sp-api/constants"
	"github.com/fond-of-vertigo/amazon-sp-api/httpx"
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
	attempts := 1
	for {
		if attempts >= constants.MaxRetryCount {
			return nil, fmt.Errorf("max retry count of %d reached", constants.MaxRetryCount)
		}

		req, err := a.createNewRequest(httpClient.GetEndpoint())
		if err != nil {
			return nil, err
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			if err := waitForRetry(resp); err != nil {
				return resp, err
			}
			attempts++
			continue
		}
		return resp, err
	}
}

func waitForRetry(resp *http.Response) error {
	backOffTime, err := getBackoffDelay(resp)
	if err != nil {
		return err
	}

	time.Sleep(backOffTime)
	return nil
}

func getBackoffDelay(resp *http.Response) (time.Duration, error) {
	backOffTimeHeader := resp.Header.Get(constants.RateLimitHeader)
	if backOffTimeHeader == "" {
		return constants.RetryDelay, nil
	}

	reqPerSecond, err := strconv.ParseFloat(backOffTimeHeader, 64)
	if err != nil {
		return constants.RetryDelay, fmt.Errorf("error parsing backoff delay: %w", err)
	}

	backOffTime := time.Duration(math.Abs(1/reqPerSecond)) * time.Second
	return backOffTime, nil
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
