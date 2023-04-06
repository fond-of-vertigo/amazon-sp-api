package orders

import (
	"encoding/json"
	"errors"
	"go/types"
	"net/http"
	"net/url"
	"time"

	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/httpx"
)

const pathPrefix = "/orders/v0"

type API struct {
	httpClient *httpx.Client
}

func NewAPI(httpClient *httpx.Client) *API {
	return &API{
		httpClient: httpClient,
	}
}

// GetOrders returns orders created or updated during the time frame indicated by the specified parameters.
// You can also apply a range of filtering criteria to narrow the list of orders returned. If NextToken is present,
// that will be used to retrieve the orders instead of other criteria.
func (a *API) GetOrders(filter *GetOrdersFilter) (*apis.CallResponse[GetOrdersResponse], error) {
	if len(filter.MarketplaceIDs) > 50 {
		return nil, errors.New("marketplaceIDs must not contain more than 50 elements")
	}
	if len(filter.AmazonOrderIDs) > 50 {
		return nil, errors.New("amazonOrderIDs must not contain more than 50 elements")
	}

	return apis.NewCall[GetOrdersResponse](http.MethodGet, pathPrefix+"/orders").
		WithQueryParams(filter.GetQuery()).
		WithRateLimit(0.0167, time.Second).
		Execute(a.httpClient)
}

// GetOrder Returns the order that you specify.
func (a *API) GetOrder(orderID string) (*apis.CallResponse[GetOrderResponse], error) {
	return apis.NewCall[GetOrderResponse](http.MethodGet, pathPrefix+"/orders/"+orderID).
		WithRateLimit(0.0167, time.Second).
		Execute(a.httpClient)
}

// GetOrderBuyerInfo returns buyer information for the order that you specify.
func (a *API) GetOrderBuyerInfo(orderID string) (*apis.CallResponse[GetOrderBuyerInfoResponse], error) {
	return apis.NewCall[GetOrderBuyerInfoResponse](http.MethodGet, pathPrefix+"/orders/"+orderID+"/buyerInfo").
		WithRateLimit(0.0167, time.Second).
		Execute(a.httpClient)
}

// GetOrderAddress returns the shipping address for the order that you specify.
func (a *API) GetOrderAddress(orderID string) (*apis.CallResponse[GetOrderAddressResponse], error) {
	return apis.NewCall[GetOrderAddressResponse](http.MethodGet, pathPrefix+"/orders/"+orderID+"/address").
		WithRateLimit(0.0167, time.Second).
		Execute(a.httpClient)
}

// GetOrderItems returns detailed order item information for the order that you specify.
// If NextToken is provided, it's used to retrieve the next page of order items.
func (a *API) GetOrderItems(orderID string, nextToken *string) (*apis.CallResponse[GetOrderItemsResponse], error) {
	params := url.Values{}
	if nextToken != nil && *nextToken != "" {
		params.Add("NextToken", *nextToken)
	}

	return apis.NewCall[GetOrderItemsResponse](http.MethodGet, pathPrefix+"/orders/"+orderID+"/orderItems").
		WithQueryParams(params).
		WithRateLimit(0.5, time.Second).
		Execute(a.httpClient)
}

// GetOrderItemsBuyerInfo returns buyer information for the order items in the order that you specify.
func (a *API) GetOrderItemsBuyerInfo(orderID string, nextToken *string) (*apis.CallResponse[GetOrderItemsBuyerInfoResponse], error) {
	params := url.Values{}
	if nextToken != nil && *nextToken != "" {
		params.Add("NextToken", *nextToken)
	}

	return apis.NewCall[GetOrderItemsBuyerInfoResponse](http.MethodGet, pathPrefix+"/orders/"+orderID+"/orderItems/buyerInfo").
		WithQueryParams(params).
		WithRateLimit(0.5, time.Second).
		Execute(a.httpClient)
}

// UpdateShipmentStatus update the shipment status for an order that you specify.
func (a *API) UpdateShipmentStatus(orderID string, payload *UpdateShipmentStatusRequest) (*apis.CallResponse[UpdateShipmentStatusErrorResponse], error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return apis.NewCall[UpdateShipmentStatusErrorResponse](http.MethodPost, pathPrefix+"/orders/"+orderID+"/shipment").
		WithBody(body).
		WithRateLimit(5, time.Second).
		Execute(a.httpClient)
}

// GetOrderRegulatedInfo returns regulated information for the order that you specify.
func (a *API) GetOrderRegulatedInfo(orderID string) (*apis.CallResponse[GetOrderRegulatedInfoResponse], error) {
	return apis.NewCall[GetOrderRegulatedInfoResponse](http.MethodGet, pathPrefix+"/orders/"+orderID+"/regulatedInfo").
		WithRateLimit(0.5, time.Second).
		Execute(a.httpClient)
}

// UpdateVerificationStatus Updates (approves or rejects) the verification status of an order containing regulated products.
func (a *API) UpdateVerificationStatus(orderID string, payload *UpdateVerificationStatusRequest) (*apis.CallResponse[UpdateVerificationStatusErrorResponse], error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return apis.NewCall[UpdateVerificationStatusErrorResponse](http.MethodPatch, pathPrefix+"/orders/"+orderID+"/regulatedInfo").
		WithBody(body).
		WithRateLimit(0.5, time.Second).
		Execute(a.httpClient)
}

// GetOrderItemsApprovals returns detailed order items approvals information for the order specified.
// If NextToken is provided, it's used to retrieve the next page of order items approvals.
func (a *API) GetOrderItemsApprovals(orderID string, filter GetOrderItemsApprovalsFilter) (*apis.CallResponse[GetOrderApprovalsResponse], error) {
	if len(filter.ItemApprovalTypes) > 1 {
		return nil, errors.New("itemApprovalTypes must not contain more than 1 element")
	}

	if len(filter.ItemApprovalStatus) > 6 {
		return nil, errors.New("itemApprovalStatus must not contain more than 6 elements")
	}

	return apis.NewCall[GetOrderApprovalsResponse](http.MethodGet, pathPrefix+"/orders/"+orderID+"/orderItems/approvals").
		WithQueryParams(filter.GetQuery()).
		WithRateLimit(0.5, time.Second).
		Execute(a.httpClient)
}

// UpdateOrderItemsApprovals updates the oder items approvals for the specified order.
func (a *API) UpdateOrderItemsApprovals(orderID string, payload *UpdateOrderApprovalsRequest) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = apis.NewCall[types.Nil](http.MethodPost, pathPrefix+"/orders/"+orderID+"/orderItems/approvals").
		WithBody(body).
		WithParseErrorListOnError(true).
		WithRateLimit(5, time.Second).
		Execute(a.httpClient)

	return err
}

// ConfirmShipment updates the shipment status for the specified order.
func (a *API) ConfirmShipment(orderID string, payload *ConfirmShipmentRequest) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = apis.NewCall[types.Nil](http.MethodPost, pathPrefix+"/orders/"+orderID+"/shipmentConfirmation").
		WithBody(body).
		WithParseErrorListOnError(true).
		WithRateLimit(2, time.Second).
		Execute(a.httpClient)
	return err
}
