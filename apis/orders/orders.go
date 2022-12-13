package orders

import (
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/httpx"
	"golang.org/x/time/rate"
	"net/http"
	"net/url"
	"time"
)

const pathPrefix = "/orders/v0"

type API interface {
	// GetOrderItems returns detailed order item information for the order that you specify.
	// If NextToken is provided, it's used to retrieve the next page of order items.
	GetOrderItems(orderID string, nextToken *string) (*apis.CallResponse[GetOrderItemsResponse], error)
}

var (
	RateLimitGetOrderItems *rate.Limiter
)

type api struct {
	HttpClient httpx.Client
}

func NewAPI(httpClient httpx.Client) API {
	RateLimitGetOrderItems = rate.NewLimiter(rate.Every(time.Second/2), 30)
	return &api{
		HttpClient: httpClient,
	}
}

func (a *api) GetOrderItems(orderID string, nextToken *string) (*apis.CallResponse[GetOrderItemsResponse], error) {
	params := url.Values{}
	if nextToken != nil {
		params.Add("NextToken", *nextToken)
	}
	return apis.NewCall[GetOrderItemsResponse](http.MethodGet, pathPrefix+"/orders/"+orderID+"/orderItems").
		WithQueryParams(params).
		WithRateLimiter(RateLimitGetOrderItems).
		Execute(a.HttpClient)
}
