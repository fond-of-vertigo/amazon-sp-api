package feeds

import (
	"encoding/json"
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/httpx"
	"go/types"
	"net/http"
)

const pathPrefix = "/feeds/2021-06-30"

type API interface {
	// GetFeeds returns feed details for the feeds that match the filters that you specify.
	GetFeeds(filter *GetFeedsRequestFilter) (*apis.CallResponse[GetFeedsResponse], error)
	// CreateFeed creates a feed. Upload the contents of the feed document before calling this operation.
	CreateFeed(specification *CreateFeedSpecification) (*apis.CallResponse[CreateFeedResponse], error)
	// GetFeed returns feed details (including the resultDocumentId, if available) for the feed that you specify.
	GetFeed(feedID string) (*apis.CallResponse[Feed], error)
	// CancelFeed cancels the feed that you specify. Only feeds with processingStatus=IN_QUEUE can be cancelled.
	// Cancelled feeds are returned in subsequent calls to the getFeed and getFeeds operations.
	CancelFeed(feedID string) error
	// CreateFeedDocument creates a feed document for the feed type that you specify.
	// This operation returns a presigned URL for uploading the feed document contents.
	// It also returns a feedDocumentId value that you can pass in with a subsequent call to the createFeed operation.
	CreateFeedDocument(specification *CreateFeedDocumentSpecification) (*apis.CallResponse[CreateFeedDocumentResponse], error)
	// GetFeedDocument the information required for retrieving a feed document's contents.
	GetFeedDocument(feedDocumentID string) (*apis.CallResponse[FeedDocument], error)
}

type api struct {
	HttpClient httpx.Client
}

func NewAPI(httpClient httpx.Client) API {
	return &api{
		HttpClient: httpClient,
	}
}

func (a *api) GetFeeds(filter *GetFeedsRequestFilter) (*apis.CallResponse[GetFeedsResponse], error) {
	return apis.NewCall[GetFeedsResponse](http.MethodGet, pathPrefix+"/feeds").
		WithQueryParams(filter.GetQuery()).
		WithParseErrorListOnError(true).
		Execute(a.HttpClient)
}

func (a *api) CreateFeed(specification *CreateFeedSpecification) (*apis.CallResponse[CreateFeedResponse], error) {
	body, err := json.Marshal(specification)
	if err != nil {
		return nil, err
	}

	return apis.NewCall[CreateFeedResponse](http.MethodPost, pathPrefix+"/feeds").
		WithBody(body).
		WithParseErrorListOnError(true).
		Execute(a.HttpClient)
}

func (a *api) GetFeed(feedID string) (*apis.CallResponse[Feed], error) {
	return apis.NewCall[Feed](http.MethodGet, pathPrefix+"/feeds/"+feedID).
		WithParseErrorListOnError(true).
		Execute(a.HttpClient)
}

func (a *api) CancelFeed(feedID string) error {
	_, err := apis.NewCall[types.Nil](http.MethodDelete, pathPrefix+"/feeds/"+feedID).
		WithParseErrorListOnError(true).
		Execute(a.HttpClient)
	return err
}

func (a *api) CreateFeedDocument(specification *CreateFeedDocumentSpecification) (*apis.CallResponse[CreateFeedDocumentResponse], error) {
	body, err := json.Marshal(specification)
	if err != nil {
		return nil, err
	}

	return apis.NewCall[CreateFeedDocumentResponse](http.MethodPost, pathPrefix+"/documents").
		WithBody(body).
		WithParseErrorListOnError(true).
		Execute(a.HttpClient)
}

func (a *api) GetFeedDocument(feedDocumentID string) (*apis.CallResponse[FeedDocument], error) {
	return apis.NewCall[FeedDocument](http.MethodGet, pathPrefix+"/documents/"+feedDocumentID).
		WithParseErrorListOnError(true).
		Execute(a.HttpClient)
}
