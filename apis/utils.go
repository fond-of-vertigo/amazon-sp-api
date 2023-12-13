package apis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func FirstNElementsOfSlice[Element any](slice []Element, n int) []Element {
	if len(slice) < n {
		return slice
	}
	return slice[:n]
}

func AddToQueryIfSet(q url.Values, key string, value string) {
	if value != "" {
		q.Set(key, value)
	}
}

func MapToCommaString[t any](slice []t) string {
	var result []string
	for _, v := range slice {
		result = append(result, fmt.Sprintf("%v", v))
	}
	return strings.Join(result, ",")
}

func mapErrorListToError(errorList *ErrorList) (errs error) {
	if errorList == nil || len(errorList.Errors) == 0 {
		return nil
	}

	for _, err := range errorList.Errors {
		errs = errors.Join(errs, fmt.Errorf("code=%s, message=%s, details=%v", err.Code, err.Message, err.Details))
	}

	return errs
}

func unmarshalBody(resp *http.Response, into any) error {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(bodyBytes) == 0 {
		return nil
	}
	return json.Unmarshal(bodyBytes, into)
}
