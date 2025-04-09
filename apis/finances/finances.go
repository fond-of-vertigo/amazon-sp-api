package finances

import (
	"errors"
	"net/http"
	"time"

	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/httpx"
)

const pathPrefix = "/finances/v0"

type API struct {
	httpClient *httpx.Client
}

func NewAPI(httpClient *httpx.Client) *API {
	return &API{
		httpClient: httpClient,
	}
}

// ListFinancialEventGroups returns financial event groups for a given date range.
func (a *API) ListFinancialEventGroups(filter *ListFinancialEventGroupsFilter) (*apis.CallResponse[ListFinancialEventGroupsResponse], error) {
	if filter.MaxResultsPerPage != nil && (*filter.MaxResultsPerPage < 1 || *filter.MaxResultsPerPage > 100) {
		return nil, errors.New("maxResultsPerPage must be between 1 and 100")
	}

	return apis.NewCall[ListFinancialEventGroupsResponse](http.MethodGet, pathPrefix+"/financialEventGroups").
		WithQueryParams(filter.GetQuery()).
		WithRateLimit(0.5, time.Second).
		Execute(a.httpClient)
}

// ListFinancialEventsByGroupID returns all financial events for the specified financial event group.
func (a *API) ListFinancialEventsByGroupID(eventGroupID string, filter *ListFinancialEventsByIDFilter) (*apis.CallResponse[ListFinancialEventsResponse], error) {
	if filter.MaxResultsPerPage != nil && (*filter.MaxResultsPerPage < 1 || *filter.MaxResultsPerPage > 100) {
		return nil, errors.New("maxResultsPerPage must be between 1 and 100")
	}

	return apis.NewCall[ListFinancialEventsResponse](http.MethodGet, pathPrefix+"/financialEventGroups/"+eventGroupID+"/financialEvents").
		WithQueryParams(filter.GetQuery()).
		WithRateLimit(0.5, time.Second).
		WithParseErrorListOnError().
		Execute(a.httpClient)
}

// ListFinancialEventsByOrderID returns all financial events for the specified order.
func (a *API) ListFinancialEventsByOrderID(orderID string, filter *ListFinancialEventsByIDFilter) (*apis.CallResponse[ListFinancialEventsResponse], error) {
	if filter.MaxResultsPerPage != nil && (*filter.MaxResultsPerPage < 1 || *filter.MaxResultsPerPage > 100) {
		return nil, errors.New("maxResultsPerPage must be between 1 and 100")
	}

	return apis.NewCall[ListFinancialEventsResponse](http.MethodGet, pathPrefix+"/orders/"+orderID+"/financialEvents").
		WithQueryParams(filter.GetQuery()).
		WithRateLimit(0.5, time.Second).
		WithParseErrorListOnError().
		Execute(a.httpClient)
}

// ListFinancialEvents returns financial events for the specified data range.
func (a *API) ListFinancialEvents(filter *ListFinancialEventsFilter) (*apis.CallResponse[ListFinancialEventsResponse], error) {
	if filter.MaxResultsPerPage != nil && (*filter.MaxResultsPerPage < 1 || *filter.MaxResultsPerPage > 100) {
		return nil, errors.New("maxResultsPerPage must be between 1 and 100")
	}

	return apis.NewCall[ListFinancialEventsResponse](http.MethodGet, pathPrefix+"/financialEvents").
		WithQueryParams(filter.GetQuery()).
		WithRateLimit(0.5, time.Second).
		WithParseErrorListOnError().
		Execute(a.httpClient)
}
