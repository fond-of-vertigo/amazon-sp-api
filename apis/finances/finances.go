package finances

import (
	"errors"
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/httpx"
	"golang.org/x/time/rate"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const pathPrefix = "/finances/v0"

type API interface {
	// ListFinancialEventGroups returns financial event groups for a given date range.
	ListFinancialEventGroups(maxResultsPerPage *int, financialEventGroupStartedBefore *apis.JsonTimeISO8601, financialEventGroupStartedAfter *apis.JsonTimeISO8601, nextToken *string) (*apis.CallResponse[ListFinancialEventGroupsResponse], error)
	// ListFinancialEventsByGroupID returns all financial events for the specified financial event group.
	ListFinancialEventsByGroupID(eventGroupID string, maxResultsPerPage *int, nextToken *string) (*apis.CallResponse[ListFinancialEventsResponse], error)
	// ListFinancialEventsByOrderID returns all financial events for the specified order.
	ListFinancialEventsByOrderID(orderID string, maxResultsPerPage *int, nextToken *string) (*apis.CallResponse[ListFinancialEventsResponse], error)
	// ListFinancialEvents returns financial events for the specified data range.
	ListFinancialEvents(maxResultsPerPage *int, postedAfter *apis.JsonTimeISO8601, postedBefore *apis.JsonTimeISO8601, nextToken *string) (*apis.CallResponse[ListFinancialEventsResponse], error)
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
		RateLimitListFinancialEventGroups:     rate.NewLimiter(rate.Every(time.Second/2), 30),
		RateLimitListFinancialEventsByGroupID: rate.NewLimiter(rate.Every(time.Second/2), 30),
		RateLimitListFinancialEventsByOrderID: rate.NewLimiter(rate.Every(time.Second/2), 30),
		RateLimitListFinancialEvents:          rate.NewLimiter(rate.Every(time.Second/2), 30),
	}
}

func (a *api) ListFinancialEventGroups(maxResultsPerPage *int, financialEventGroupStartedBefore *apis.JsonTimeISO8601, financialEventGroupStartedAfter *apis.JsonTimeISO8601, nextToken *string) (*apis.CallResponse[ListFinancialEventGroupsResponse], error) {
	params := url.Values{}
	if maxResultsPerPage != nil {
		if *maxResultsPerPage < 1 || *maxResultsPerPage > 100 {
			return nil, errors.New("maxResultsPerPage must be between 1 and 100")
		}

		params.Add("MaxResultsPerPage", strconv.Itoa(*maxResultsPerPage))
	}
	if financialEventGroupStartedBefore != nil {
		params.Add("FinancialEventGroupStartedBefore", financialEventGroupStartedBefore.String())
	}
	if financialEventGroupStartedAfter != nil {
		params.Add("FinancialEventGroupStartedAfter", financialEventGroupStartedAfter.String())
	}
	if nextToken != nil {
		params.Add("NextToken", *nextToken)
	}

	return apis.NewCall[ListFinancialEventGroupsResponse](http.MethodGet, pathPrefix+"/financialEventGroups").
		WithQueryParams(params).
		WithRateLimiter(a.RateLimitListFinancialEventGroups).
		Execute(a.HttpClient)
}

func (a *api) ListFinancialEventsByGroupID(eventGroupID string, maxResultsPerPage *int, nextToken *string) (*apis.CallResponse[ListFinancialEventsResponse], error) {
	params := url.Values{}
	if maxResultsPerPage != nil {
		if *maxResultsPerPage < 1 || *maxResultsPerPage > 100 {
			return nil, errors.New("maxResultsPerPage must be between 1 and 100")
		}

		params.Add("MaxResultsPerPage", strconv.Itoa(*maxResultsPerPage))
	}
	if nextToken != nil {
		params.Add("NextToken", *nextToken)
	}

	return apis.NewCall[ListFinancialEventsResponse](http.MethodGet, pathPrefix+"/financialEventGroups/"+eventGroupID+"/financialEvents").
		WithQueryParams(params).
		WithRateLimiter(a.RateLimitListFinancialEventsByGroupID).
		Execute(a.HttpClient)
}

func (a *api) ListFinancialEventsByOrderID(orderID string, maxResultsPerPage *int, nextToken *string) (*apis.CallResponse[ListFinancialEventsResponse], error) {
	params := url.Values{}
	if maxResultsPerPage != nil {
		if *maxResultsPerPage < 1 || *maxResultsPerPage > 100 {
			return nil, errors.New("maxResultsPerPage must be between 1 and 100")
		}

		params.Add("MaxResultsPerPage", strconv.Itoa(*maxResultsPerPage))
	}
	if nextToken != nil {
		params.Add("NextToken", *nextToken)
	}

	return apis.NewCall[ListFinancialEventsResponse](http.MethodGet, pathPrefix+"/orders/"+orderID+"/financialEvents").
		WithQueryParams(params).
		WithRateLimiter(a.RateLimitListFinancialEventsByOrderID).
		Execute(a.HttpClient)
}

func (a *api) ListFinancialEvents(maxResultsPerPage *int, postedAfter *apis.JsonTimeISO8601, postedBefore *apis.JsonTimeISO8601, nextToken *string) (*apis.CallResponse[ListFinancialEventsResponse], error) {
	params := url.Values{}
	if maxResultsPerPage != nil {
		if *maxResultsPerPage < 1 || *maxResultsPerPage > 100 {
			return nil, errors.New("maxResultsPerPage must be between 1 and 100")
		}

		params.Add("MaxResultsPerPage", strconv.Itoa(*maxResultsPerPage))
	}
	if postedAfter != nil {
		params.Add("PostedAfter", postedAfter.String())
	}
	if postedBefore != nil {
		params.Add("PostedBefore", postedBefore.String())
	}
	if nextToken != nil {
		params.Add("NextToken", *nextToken)
	}

	return apis.NewCall[ListFinancialEventsResponse](http.MethodGet, pathPrefix+"/financialEvents").
		WithQueryParams(params).
		WithRateLimiter(a.RateLimitListFinancialEvents).
		Execute(a.HttpClient)
}
