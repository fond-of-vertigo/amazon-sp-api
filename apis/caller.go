package apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HttpRequestDoer interface {
	Do(*http.Request) (*http.Response, error)
	GetEndpoint() string
}

type APICall struct {
	Method      string
	APIPath     string
	QueryParams url.Values
	Body        []byte
	// RestrictedDataToken is optional and can be passed to replace the existing accessToken
	RestrictedDataToken *string
}

func CallAPIWithResponseType[responseType any](callParams APICall, httpClient HttpRequestDoer) (*responseType, error) {
	_, bodyBytes, err := CallAPI(callParams, httpClient)
	if err != nil {
		return nil, err
	}

	var reportResp responseType
	err = json.Unmarshal(bodyBytes, &reportResp)
	return &reportResp, err
}

func CallAPIIgnoreResponse(callParams APICall, httpClient HttpRequestDoer) error {
	_, _, err := CallAPI(callParams, httpClient)
	return err
}

func CallAPI(callParams APICall, httpClient HttpRequestDoer) (*http.Response, []byte, error) {
	callParams.APIPath = httpClient.GetEndpoint() + callParams.APIPath

	req, err := createNewRequest(callParams)
	if err != nil {
		return nil, nil, err
	}

	if callParams.RestrictedDataToken != nil && *callParams.RestrictedDataToken != "" {
		req.Header.Add("X-Amz-Access-Token", *callParams.RestrictedDataToken)
	}

	return executeRequest(httpClient, req)
}

func createNewRequest(callParams APICall) (*http.Request, error) {
	apiPath, err := url.Parse(callParams.APIPath)
	if err != nil {
		return nil, err
	}
	apiPath.RawQuery = callParams.QueryParams.Encode()

	return http.NewRequest(callParams.Method, apiPath.String(), bytes.NewBuffer(callParams.Body))
}

func executeRequest(httpClient HttpRequestDoer, req *http.Request) (*http.Response, []byte, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	bodyBytes, err := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		var errorList ErrorList
		if err = json.Unmarshal(bodyBytes, &errorList); err != nil {
			return nil, nil, fmt.Errorf("could not unmarshal ErrorList %w", err)
		}

		return nil, nil, fmt.Errorf("%v", errorList)
	}

	return resp, bodyBytes, err
}
