package orders

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/httpx"
)

const pathPrefix = "/orders/v0"

type API interface {
	// GetOrders returns orders created or updated during the time frame indicated by the specified parameters.
	// You can also apply a range of filtering criteria to narrow the list of orders returned. If NextToken is present,
	// that will be used to retrieve the orders instead of other criteria.
	GetOrders(filter *GetOrdersFilter) (*apis.CallResponse[GetOrdersResponse], error)
	// GetOrder Returns the order that you specify.
	GetOrder(orderID string) (*apis.CallResponse[GetOrderResponse], error)
	// GetOrderBuyerInfo returns buyer information for the order that you specify.
	GetOrderBuyerInfo(orderID string) (*apis.CallResponse[GetOrderBuyerInfoResponse], error)
	// GetOrderAddress returns the shipping address for the order that you specify.
	GetOrderAddress(orderID string) (*apis.CallResponse[GetOrderAddressResponse], error)
	// GetOrderItems returns detailed order item information for the order that you specify.
	// If NextToken is provided, it's used to retrieve the next page of order items.
	GetOrderItems(orderID string, nextToken *string) (*apis.CallResponse[GetOrderItemsResponse], error)
	// GetOrderItemsBuyerInfo returns buyer information for the order items in the order that you specify.
	GetOrderItemsBuyerInfo(orderID string, nextToken *string) (*apis.CallResponse[GetOrderItemsBuyerInfoResponse], error)
	// UpdateShipmentStatus update the shipment status for an order that you specify.
	UpdateShipmentStatus(orderID string, payload *UpdateShipmentStatusRequest) (*apis.CallResponse[UpdateShipmentStatusErrorResponse], error)
	// GetOrderRegulatedInfo returns regulated information for the order that you specify.
	GetOrderRegulatedInfo(orderID string) (*apis.CallResponse[GetOrderRegulatedInfoResponse], error)
	// UpdateVerificationStatus Updates (approves or rejects) the verification status of an order containing regulated products.
	UpdateVerificationStatus(orderID string, payload *UpdateVerificationStatusRequest) (*apis.CallResponse[UpdateVerificationStatusErrorResponse], error)
}

type api struct {
	HttpClient httpx.Client
}

func NewAPI(httpClient httpx.Client) API {
	return &api{
		HttpClient: httpClient,
	}
}

func (a *api) GetOrders(filter *GetOrdersFilter) (*apis.CallResponse[GetOrdersResponse], error) {
	if len(filter.MarketplaceIDs) > 50 {
		return nil, errors.New("marketplaceIDs must not contain more than 50 elements")
	}
	if len(filter.AmazonOrderIDs) > 50 {
		return nil, errors.New("amazonOrderIDs must not contain more than 50 elements")
	}

	return apis.NewCall[GetOrdersResponse](http.MethodGet, pathPrefix+"/orders").
		WithQueryParams(filter.GetQuery()).
		WithRateLimit(0.0167, time.Second).
		Execute(a.HttpClient)
}

func (a *api) GetOrder(orderID string) (*apis.CallResponse[GetOrderResponse], error) {
	return apis.NewCall[GetOrderResponse](http.MethodGet, pathPrefix+"/orders/"+orderID).
		WithRateLimit(0.0167, time.Second).
		Execute(a.HttpClient)
}

func (a *api) GetOrderBuyerInfo(orderID string) (*apis.CallResponse[GetOrderBuyerInfoResponse], error) {
	return apis.NewCall[GetOrderBuyerInfoResponse](http.MethodGet, pathPrefix+"/orders/"+orderID+"/buyerInfo").
		WithRateLimit(0.0167, time.Second).
		Execute(a.HttpClient)
}

func (a *api) GetOrderAddress(orderID string) (*apis.CallResponse[GetOrderAddressResponse], error) {
	return apis.NewCall[GetOrderAddressResponse](http.MethodGet, pathPrefix+"/orders/"+orderID+"/address").
		WithRateLimit(0.0167, time.Second).
		Execute(a.HttpClient)
}

func (a *api) GetOrderItems(orderID string, nextToken *string) (*apis.CallResponse[GetOrderItemsResponse], error) {
	params := url.Values{}
	if nextToken != nil && *nextToken != "" {
		params.Add("NextToken", *nextToken)
	}

	return apis.NewCall[GetOrderItemsResponse](http.MethodGet, pathPrefix+"/orders/"+orderID+"/orderItems").
		WithQueryParams(params).
		WithRateLimit(0.5, time.Second).
		Execute(a.HttpClient)
}

func (a *api) GetOrderItemsBuyerInfo(orderID string, nextToken *string) (*apis.CallResponse[GetOrderItemsBuyerInfoResponse], error) {
	params := url.Values{}
	if nextToken != nil && *nextToken != "" {
		params.Add("NextToken", *nextToken)
	}

	return apis.NewCall[GetOrderItemsBuyerInfoResponse](http.MethodGet, pathPrefix+"/orders/"+orderID+"/orderItems/buyerInfo").
		WithQueryParams(params).
		WithRateLimit(0.5, time.Second).
		Execute(a.HttpClient)
}

func (a *api) UpdateShipmentStatus(orderID string, payload *UpdateShipmentStatusRequest) (*apis.CallResponse[UpdateShipmentStatusErrorResponse], error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return apis.NewCall[UpdateShipmentStatusErrorResponse](http.MethodPost, pathPrefix+"/orders/"+orderID+"/shipment").
		WithBody(body).
		WithRateLimit(5, time.Second).
		Execute(a.HttpClient)
}

func (a *api) GetOrderRegulatedInfo(orderID string) (*apis.CallResponse[GetOrderRegulatedInfoResponse], error) {
	return apis.NewCall[GetOrderRegulatedInfoResponse](http.MethodGet, pathPrefix+"/orders/"+orderID+"/regulatedInfo").
		WithRateLimit(0.5, time.Second).
		Execute(a.HttpClient)
}

func (a *api) UpdateVerificationStatus(orderID string, payload *UpdateVerificationStatusRequest) (*apis.CallResponse[UpdateVerificationStatusErrorResponse], error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return apis.NewCall[UpdateVerificationStatusErrorResponse](http.MethodPatch, pathPrefix+"/orders/"+orderID+"/regulatedInfo").
		WithBody(body).
		WithRateLimit(0.5, time.Second).
		Execute(a.HttpClient)
}
