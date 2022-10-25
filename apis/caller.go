package apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fond-of-vertigo/amazon-sp-api/apis/reports"
	"io"
	"net/http"
	"net/url"
)

type HttpRequestDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type APICall struct {
	Method      string
	APIPath     string
	QueryParams url.Values
	Body        []byte
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
	req, err := createNewRequest(callParams)
	if err != nil {
		return nil, nil, err
	}

	return executeRequest(err, httpClient, req)
}

func createNewRequest(callParams APICall) (*http.Request, error) {
	apiPath, err := url.Parse(callParams.APIPath)
	if err != nil {
		return nil, err
	}
	apiPath.RawQuery = callParams.QueryParams.Encode()

	return http.NewRequest(callParams.Method, apiPath.RawPath, bytes.NewReader(callParams.Body))
}

func executeRequest(err error, httpClient HttpRequestDoer, req *http.Request) (*http.Response, []byte, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	bodyBytes, err := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		var errorList reports.ErrorList
		if err = json.Unmarshal(bodyBytes, &errorList); err != nil {
			return nil, nil, fmt.Errorf("could not unmarshal ErrorList %w", err)
		}

		return nil, nil, fmt.Errorf("%v", errorList)
	}

	return resp, bodyBytes, err
}
