package feeds

import (
	"encoding/json"
	"go/types"
	"net/http"
	"time"

	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/httpx"
)

const pathPrefix = "/feeds/2021-06-30"

type API struct {
	httpClient *httpx.Client
}

func NewAPI(httpClient *httpx.Client) *API {
	return &API{
		httpClient: httpClient,
	}
}

// GetFeeds returns feed details for the feeds that match the filters that you specify.
func (a *API) GetFeeds(filter *GetFeedsRequestFilter) (*apis.CallResponse[GetFeedsResponse], error) {
	return apis.NewCall[GetFeedsResponse](http.MethodGet, pathPrefix+"/feeds").
		WithQueryParams(filter.GetQuery()).
		WithParseErrorListOnError().
		WithRateLimit(0.0222, time.Second).
		Execute(a.httpClient)
}

// CreateFeed creates a feed. Upload the contents of the feed document before calling this operation.
func (a *API) CreateFeed(specification *CreateFeedSpecification) (*apis.CallResponse[CreateFeedResponse], error) {
	body, err := json.Marshal(specification)
	if err != nil {
		return nil, err
	}

	return apis.NewCall[CreateFeedResponse](http.MethodPost, pathPrefix+"/feeds").
		WithBody(body).
		WithParseErrorListOnError().
		WithRateLimit(0.0083, time.Second).
		Execute(a.httpClient)
}

// GetFeed returns feed details (including the resultDocumentId, if available) for the feed that you specify.
func (a *API) GetFeed(feedID string) (*apis.CallResponse[Feed], error) {
	return apis.NewCall[Feed](http.MethodGet, pathPrefix+"/feeds/"+feedID).
		WithParseErrorListOnError().
		WithRateLimit(2, time.Second).
		Execute(a.httpClient)
}

// CancelFeed cancels the feed that you specify. Only feeds with processingStatus=IN_QUEUE can be cancelled.
// Cancelled feeds are returned in subsequent calls to the getFeed and getFeeds operations.
func (a *API) CancelFeed(feedID string) error {
	_, err := apis.NewCall[types.Nil](http.MethodDelete, pathPrefix+"/feeds/"+feedID).
		WithParseErrorListOnError().
		WithRateLimit(0.0222, time.Second).
		Execute(a.httpClient)
	return err
}

// CreateFeedDocument creates a feed document for the feed type that you specify.
// This operation returns a presigned URL for uploading the feed document contents.
// It also returns a feedDocumentId value that you can pass in with a subsequent call to the createFeed operation.
func (a *API) CreateFeedDocument(specification *CreateFeedDocumentSpecification) (*apis.CallResponse[CreateFeedDocumentResponse], error) {
	body, err := json.Marshal(specification)
	if err != nil {
		return nil, err
	}

	return apis.NewCall[CreateFeedDocumentResponse](http.MethodPost, pathPrefix+"/documents").
		WithBody(body).
		WithParseErrorListOnError().
		WithRateLimit(0.0083, time.Second).
		Execute(a.httpClient)
}

// GetFeedDocument the information required for retrieving a feed document's contents.
func (a *API) GetFeedDocument(feedDocumentID string) (*apis.CallResponse[FeedDocument], error) {
	return apis.NewCall[FeedDocument](http.MethodGet, pathPrefix+"/documents/"+feedDocumentID).
		WithParseErrorListOnError().
		WithRateLimit(1.0, time.Minute). // documented value (2/sec) seems way too much (many http 429 errors)
		Execute(a.httpClient)
}
