package orders

import (
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"net/http"
	"net/url"
)

const pathPrefix = "/orders/v0"

type API interface {
	// GetOrderItems returns detailed order item information for the order that you specify.
	// If NextToken is provided, it's used to retrieve the next page of order items.
	GetOrderItems(orderID string, nextToken *string) (*GetOrderItemsResponse, apis.CallError)
}

type api struct {
	HttpClient apis.HttpRequestDoer
}

func NewAPI(httpClient apis.HttpRequestDoer) API {
	return &api{
		HttpClient: httpClient,
	}
}

func (a *api) GetOrderItems(orderID string, nextToken *string) (*GetOrderItemsResponse, apis.CallError) {
	params := url.Values{}
	if nextToken != nil {
		params.Add("NextToken", *nextToken)
	}
	return apis.NewCall[GetOrderItemsResponse](http.MethodGet, pathPrefix+"/orders/"+orderID+"/orderItems").
		WithQueryParams(params).
		Execute(a.HttpClient)
}
