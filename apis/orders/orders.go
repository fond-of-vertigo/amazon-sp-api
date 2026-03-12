package orders

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/httpx"
)

const pathPrefix = "/orders/2026-01-01"

type API struct {
	httpClient *httpx.Client
}

func NewAPI(httpClient *httpx.Client) *API {
	return &API{
		httpClient: httpClient,
	}
}

// GetOrder returns the order that you specify.
// includedData is optional and specifies which datasets to include in the response.
// A restrictedDataToken is optional and may be passed to receive Personally Identifiable Information (PII).
func (a *API) GetOrder(orderID string, includedData []IncludedData, restrictedDataToken *string) (*apis.CallResponse[GetOrderResponse], error) {
	call := apis.NewCall[GetOrderResponse](http.MethodGet, pathPrefix+"/orders/"+orderID).
		WithRateLimit(0.0167, time.Second).
		WithRestrictedDataToken(restrictedDataToken)

	if len(includedData) > 0 {
		vals := make([]string, len(includedData))
		for i, d := range includedData {
			vals[i] = string(d)
		}
		call = call.WithQueryParams(url.Values{
			"includedData": {strings.Join(vals, ",")},
		})
	}

	return call.Execute(a.httpClient)
}
