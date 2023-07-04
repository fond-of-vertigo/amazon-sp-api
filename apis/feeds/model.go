package feeds

import (
	"github.com/fond-of-vertigo/amazon-sp-api/internal/utils"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/constants"
)

type ProcessingStatus string

const (
	// ProcessingStatusCanceled The feed was cancelled before it started processing.
	ProcessingStatusCanceled ProcessingStatus = "CANCELLED"
	// ProcessingStatusDone The feed has completed processing. Examine the contents of the result document to determine if there were any errors during processing.
	ProcessingStatusDone ProcessingStatus = "DONE"
	// ProcessingStatusFatal The feed was aborted due to a fatal error. Some, none, or all of the operations within the feed may have completed successfully.
	ProcessingStatusFatal ProcessingStatus = "FATAL"
	// ProcessingStatusInProgress The feed is being processed.
	ProcessingStatusInProgress ProcessingStatus = "IN_PROGRESS"
	// ProcessingStatusInQueue The feed has not yet started processing. It may be waiting for another IN_PROGRESS feed.
	ProcessingStatusInQueue ProcessingStatus = "IN_QUEUE"
)

// Feed contains detailed information about the feed.
type Feed struct {
	// The identifier for the feed. This identifier is unique only in combination with a seller ID.
	FeedId string `json:"feedId"`
	// The feed type.
	FeedType string `json:"feedType"`
	// A list of identifiers for the marketplaces that the feed is applied to.
	MarketplaceIDs []constants.MarketplaceID `json:"marketplaceIds,omitempty"`
	// The date and time when the feed was created, in ISO 8601 date time format.
	CreatedTime time.Time `json:"createdTime"`
	// The processing status of the feed.
	ProcessingStatus ProcessingStatus `json:"processingStatus"`
	// The date and time when feed processing started, in ISO 8601 date time format.
	ProcessingStartTime *time.Time `json:"processingStartTime,omitempty"`
	// The date and time when feed processing completed, in ISO 8601 date time format.
	ProcessingEndTime *time.Time `json:"processingEndTime,omitempty"`
	// The identifier for the feed document. This identifier is unique only in combination with a seller ID.
	ResultFeedDocumentId *string `json:"resultFeedDocumentId,omitempty"`
}

// CreateFeedDocumentResponse is the response schema for the createFeedDocument operation.
type CreateFeedDocumentResponse struct {
	// The identifier of the feed document.
	FeedDocumentId string `json:"feedDocumentId"`
	// The presigned URL for uploading the feed contents. This URL expires after 5 minutes.
	Url string `json:"url"`
}

// CreateFeedSpecification information required to create the feed."
type CreateFeedSpecification struct {
	// The feed type.
	FeedType string `json:"feedType"`
	// A list of identifiers for marketplaces that you want the feed to be applied to.
	MarketplaceIDs []constants.MarketplaceID `json:"marketplaceIds"`
	// The document identifier returned by the createFeedDocument operation. Upload the feed document contents before
	// calling the createFeed operation.
	InputFeedDocumentId string `json:"inputFeedDocumentId"`
	// Additional options to control the feed. These vary by feed type.
	FeedOptions *map[string]string `json:"feedOptions,omitempty"`
}

// CreateFeedDocumentSpecification specifies the content type for the createFeedDocument operation.
type CreateFeedDocumentSpecification struct {
	// The content type of the feed.
	ContentType string `json:"contentType"`
}

// GetFeedsRequestFilter specifies optional filters for the getFeeds operation.
type GetFeedsRequestFilter struct {
	// A list of feed types used to filter feeds. When feedTypes is provided, the other filter parameters
	// (processingStatuses, marketplaceIds, createdSince, createdUntil) and pageSize may also be provided.
	// Either feedTypes or nextToken is required. Maximum 10 feed types. If longer the first 10 will be used.
	FeedTypes []string `json:"feedTypes,omitempty"`
	// A list of marketplace identifiers used to filter feeds.
	// The feeds returned will match at least one of the marketplaces that you specify.
	// Maximum 10 marketplace identifiers. If longer the first 10 will be used.
	MarketplaceIDs []constants.MarketplaceID `json:"marketplaceIds,omitempty"`
	// The maximum number of feeds to return in a single call.
	// Minimum 1. Maximum 100.
	PageSize int `json:"pageSize,omitempty"`
	// A list of processing statuses used to filter feeds.
	ProcessingStatuses []string `json:"processingStatuses,omitempty"`
	// The earliest feed creation date and time for feeds included in the response, in ISO 8601 format.
	//The default is 90 days ago. Feeds are retained for a maximum of 90 days.
	CreatedSince apis.JsonTimeISO8601 `json:"createdSince,omitempty"`
	// The latest feed creation date and time for feeds included in the response, in ISO 8601 format.
	// The default is now.
	CreatedUntil apis.JsonTimeISO8601 `json:"createdUntil,omitempty"`
	// The token returned by a previous call to this operation.
	NextToken string `json:"nextToken,omitempty"`
}

func (f *GetFeedsRequestFilter) GetQuery() url.Values {
	q := url.Values{}

	feedTypes := strings.Join(utils.FirstNElementsOfSlice(f.FeedTypes, 10), ",")
	if feedTypes != "" {
		q.Set("feedTypes", feedTypes)
	}

	topTenMarketplaceIDs := utils.FirstNElementsOfSlice(f.MarketplaceIDs, 10)
	marketplaceIds := utils.MapToCommaString(topTenMarketplaceIDs)
	if marketplaceIds != "" {
		q.Set("marketplaceIds", marketplaceIds)
	}

	if f.PageSize > 0 && f.PageSize <= 100 {
		q.Set("pageSize", strconv.Itoa(f.PageSize))
	}

	processingStatuses := strings.Join(f.ProcessingStatuses, ",")
	if processingStatuses != "" {
		q.Set("processingStatuses", processingStatuses)
	}

	if !f.CreatedSince.IsZero() {
		q.Set("createdSince", f.CreatedSince.String())
	}

	if !f.CreatedUntil.IsZero() {
		q.Set("createdUntil", f.CreatedUntil.String())
	}

	if f.NextToken != "" {
		q.Set("nextToken", f.NextToken)
	}

	return q
}

// GetFeedsResponse is the response schema for the getFeeds operation.
type GetFeedsResponse struct {
	// A list of feeds.
	Feeds []Feed `json:"feeds"`
	// Returned when the number of results exceeds pageSize. To get the next page of results, call the getFeeds
	// operation with this token as the only parameter.
	NextToken *string `json:"nextToken,omitempty"`
}

// FeedDocument contains information about the feed document.
type FeedDocument struct {
	// The identifier for the feed document. This identifier is unique only in combination with a seller ID.
	FeedDocumentId string `json:"feedDocumentId"`
	// A presigned URL for the feed document. If `compressionAlgorithm` is not returned, you can download the feed
	// directly from this URL. This URL expires after 5 minutes.
	Url string `json:"url"`
	// If the feed document contents have been compressed, the compression algorithm used is returned in this property
	// and you must decompress the feed when you download. Otherwise, you can download the feed directly.
	CompressionAlgorithm *string `json:"compressionAlgorithm,omitempty"`
}

// CreateFeedResponse is the response schema for the createFeed operation.
type CreateFeedResponse struct {
	// The identifier for the feed. This identifier is unique only in combination with a seller ID.
	FeedId string `json:"feedId"`
}
