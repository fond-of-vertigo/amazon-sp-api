package apis

import (
	"fmt"
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
