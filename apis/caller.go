package apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fond-of-vertigo/amazon-sp-api/constants"
	"io"
	"net/http"
	"net/url"
)

type HttpRequestDoer interface {
	Do(*http.Request) (*http.Response, error)
	GetEndpoint() constants.Endpoint
}

type Call[responseType any] interface {
	WithQueryParams(url.Values) Call[responseType]
	WithBody([]byte) Call[responseType]
	// WithRestrictedDataToken is optional and can be passed to replace the existing accessToken
	WithRestrictedDataToken(*string) Call[responseType]
	// Execute will return response object on success
	Execute(httpClient HttpRequestDoer) (*responseType, CallError)
}

func NewCall[responseType any](method string, url string) Call[responseType] {
	return &call[responseType]{
		Method: method,
		URL:    url,
	}
}

type call[responseType any] struct {
	Method              string
	URL                 string
	QueryParams         url.Values
	Body                []byte
	RestrictedDataToken *string
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

func (a *call[responseType]) Execute(httpClient HttpRequestDoer) (*responseType, CallError) {
	resp, err := a.execute(httpClient)

	if err != nil {
		return nil, &callError{err: err}
	}

	if !IsSuccess(resp.StatusCode) {
		callErr := &callError{
			err: fmt.Errorf("request %s %s failed with status code %d", resp.Request.Method, resp.Request.URL, resp.StatusCode),
		}
		if err = unmarshalBody(resp, &callErr.errorList); err != nil {
			return nil, &callError{err: err}
		}
		return nil, callErr
	}
	if resp.ContentLength == 0 {
		return nil, nil
	}
	var reportResp responseType
	if err = unmarshalBody(resp, &reportResp); err != nil {
		return nil, &callError{err: err}
	}
	return &reportResp, nil
}

func (a *call[responseType]) execute(httpClient HttpRequestDoer) (*http.Response, error) {

	req, err := a.createNewRequest(httpClient.GetEndpoint())
	if err != nil {
		return nil, err
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
func IsSuccess(status int) bool {
	return status >= http.StatusOK && status < http.StatusMultipleChoices
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
