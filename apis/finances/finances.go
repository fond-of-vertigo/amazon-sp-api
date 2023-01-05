package finances

import (
	"errors"
	"net/http"
	"time"

	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/httpx"
	"golang.org/x/time/rate"
)

const pathPrefix = "/finances/v0"

type API interface {
	// ListFinancialEventGroups returns financial event groups for a given date range.
	ListFinancialEventGroups(filter *ListFinancialEventGroupsFilter) (*apis.CallResponse[ListFinancialEventGroupsResponse], error)
	// ListFinancialEventsByGroupID returns all financial events for the specified financial event group.
	ListFinancialEventsByGroupID(eventGroupID string, filter *ListFinancialEventsByIDFilter) (*apis.CallResponse[ListFinancialEventsResponse], error)
	// ListFinancialEventsByOrderID returns all financial events for the specified order.
	ListFinancialEventsByOrderID(orderID string, filter *ListFinancialEventsByIDFilter) (*apis.CallResponse[ListFinancialEventsResponse], error)
	// ListFinancialEvents returns financial events for the specified data range.
	ListFinancialEvents(filter *ListFinancialEventsFilter) (*apis.CallResponse[ListFinancialEventsResponse], error)
}

type api struct {
	HttpClient                            httpx.Client
	RateLimitListFinancialEventGroups     *rate.Limiter
	RateLimitListFinancialEventsByGroupID *rate.Limiter
	RateLimitListFinancialEventsByOrderID *rate.Limiter
	RateLimitListFinancialEvents          *rate.Limiter
}

func NewAPI(httpClient httpx.Client) API {
	return &api{
		HttpClient:                            httpClient,
		RateLimitListFinancialEventGroups:     rate.NewLimiter(rate.Every(time.Minute/30), 30),
		RateLimitListFinancialEventsByGroupID: rate.NewLimiter(rate.Every(time.Minute/30), 30),
		RateLimitListFinancialEventsByOrderID: rate.NewLimiter(rate.Every(time.Minute/30), 30),
		RateLimitListFinancialEvents:          rate.NewLimiter(rate.Every(time.Minute/30), 30),
	}
}

func (a *api) ListFinancialEventGroups(filter *ListFinancialEventGroupsFilter) (*apis.CallResponse[ListFinancialEventGroupsResponse], error) {
	if filter.MaxResultsPerPage != nil && (*filter.MaxResultsPerPage < 1 || *filter.MaxResultsPerPage > 100) {
		return nil, errors.New("maxResultsPerPage must be between 1 and 100")
	}

	return apis.NewCall[ListFinancialEventGroupsResponse](http.MethodGet, pathPrefix+"/financialEventGroups").
		WithQueryParams(filter.GetQuery()).
		WithRateLimiter(a.RateLimitListFinancialEventGroups).
		Execute(a.HttpClient)
}

func (a *api) ListFinancialEventsByGroupID(eventGroupID string, filter *ListFinancialEventsByIDFilter) (*apis.CallResponse[ListFinancialEventsResponse], error) {
	if filter.MaxResultsPerPage != nil && (*filter.MaxResultsPerPage < 1 || *filter.MaxResultsPerPage > 100) {
		return nil, errors.New("maxResultsPerPage must be between 1 and 100")
	}

	return apis.NewCall[ListFinancialEventsResponse](http.MethodGet, pathPrefix+"/financialEventGroups/"+eventGroupID+"/financialEvents").
		WithQueryParams(filter.GetQuery()).
		WithRateLimiter(a.RateLimitListFinancialEventsByGroupID).
		Execute(a.HttpClient)
}

func (a *api) ListFinancialEventsByOrderID(orderID string, filter *ListFinancialEventsByIDFilter) (*apis.CallResponse[ListFinancialEventsResponse], error) {
	if filter.MaxResultsPerPage != nil && (*filter.MaxResultsPerPage < 1 || *filter.MaxResultsPerPage > 100) {
		return nil, errors.New("maxResultsPerPage must be between 1 and 100")
	}

	return apis.NewCall[ListFinancialEventsResponse](http.MethodGet, pathPrefix+"/orders/"+orderID+"/financialEvents").
		WithQueryParams(filter.GetQuery()).
		WithRateLimiter(a.RateLimitListFinancialEventsByOrderID).
		Execute(a.HttpClient)
}

func (a *api) ListFinancialEvents(filter *ListFinancialEventsFilter) (*apis.CallResponse[ListFinancialEventsResponse], error) {
	if filter.MaxResultsPerPage != nil && (*filter.MaxResultsPerPage < 1 || *filter.MaxResultsPerPage > 100) {
		return nil, errors.New("maxResultsPerPage must be between 1 and 100")
	}

	return apis.NewCall[ListFinancialEventsResponse](http.MethodGet, pathPrefix+"/financialEvents").
		WithQueryParams(filter.GetQuery()).
		WithRateLimiter(a.RateLimitListFinancialEvents).
		Execute(a.HttpClient)
}
