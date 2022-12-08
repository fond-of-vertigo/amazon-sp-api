package reports

import (
	"fmt"
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/constants"
	"net/url"
	"strings"
	"time"
)

// Type of report
type Type string

const (
	FBAAmazonFulfilledShipmentsInvoicing Type = "GET_AMAZON_FULFILLED_SHIPMENTS_DATA_INVOICING"
)

// ReportModel Detailed information about the report.
type ReportModel struct {
	// A list of marketplace identifiers for the report.
	MarketplaceIDs []constants.MarketplaceID `json:"marketplaceIds,omitempty"`
	// The identifier for the report. This identifier is unique only in combination with a seller ID.
	ReportID string `json:"reportId"`
	// The report type.
	ReportType Type `json:"reportType"`
	// The start of a date and time range used for selecting the data to report.
	DataStartTime *time.Time `json:"dataStartTime,omitempty"`
	// The end of a date and time range used for selecting the data to report.
	DataEndTime *time.Time `json:"dataEndTime,omitempty"`
	// The identifier of the report schedule that created this report (if any). This identifier is unique only in combination with a seller ID.
	ReportScheduleID *string `json:"reportScheduleId,omitempty"`
	// The date and time when the report was created.
	CreatedTime time.Time `json:"createdTime"`
	// The processing status of the report.
	ProcessingStatus constants.ProcessingStatus `json:"processingStatus"`
	// The date and time when the report processing started, in ISO 8601 date time format.
	ProcessingStartTime *time.Time `json:"processingStartTime,omitempty"`
	// The date and time when the report processing completed, in ISO 8601 date time format.
	ProcessingEndTime *time.Time `json:"processingEndTime,omitempty"`
	// The identifier for the report document. Pass this into the getReportDocument operation to get the information you will need to retrieve the report document's contents.
	ReportDocumentID *string `json:"reportDocumentId,omitempty"`
}

// GetDocumentAPIPath returns the URL /reports/xxxx-xx-xx/documents/documentID which can be
// used for RestrictedDataTokens (RDTs) generation
func (r *ReportModel) GetDocumentAPIPath() string {
	if r.ReportDocumentID == nil {
		return ""
	}

	return pathPrefix + "/documents/" + *r.ReportDocumentID
}

type GetReportsFilter struct {
	// reportTypes is a list of report types used to filter reports.
	// When reportTypes is provided, the other filter parameters
	// (processingStatuses, marketplaceIDs, createdSince, createdUntil) and pageSize may also be provided.
	// Either reportTypes or nextToken is required.
	// Min count 1, max count 10
	reportTypes []string
	// processingStatuses is a list of processing statuses used to filter reports.
	processingStatuses []string
	//marketplaceIDs is a list of marketplace identifiers used to filter reports.
	// The reports returned will match at least one of the marketplaces that you specify.
	// min count 1, max count 10
	marketplaceIDs []string
	// pageSize is the maximum number of reports to return in a single call.
	// min 1, max 100
	pageSize int
	// createdSince is the earliest report creation date and time for reports to include in the response, in ISO 8601 date time format.
	// The default is 90 days ago. ReportsAPI are retained for a maximum of 90 days.
	createdSince apis.JsonTimeISO8601
	// createdUntil is the latest report creation date and time for reports to include in the response, in ISO 8601 date time format.
	// The default is now.
	createdUntil apis.JsonTimeISO8601
	// nextToken is a string token returned in the response to your previous request.
	// nextToken is returned when the number of results exceeds the specified pageSize value.
	//To get the next page of results, call the getReports operation and include this token as the only parameter.
	// Specifying nextToken with any other parameters will cause the request to fail.
	nextToken string
}

func (f *GetReportsFilter) GetQuery() url.Values {
	q := url.Values{}
	q.Add("reportTypes", strings.Join(f.reportTypes, ","))
	q.Add("processingStatuses", strings.Join(f.processingStatuses, ","))
	q.Add("marketplaceIds", strings.Join(f.marketplaceIDs, ","))
	q.Add("pageSize", fmt.Sprint(f.pageSize))
	q.Add("createdSince", f.createdSince.String())
	q.Add("createdUntil", f.createdUntil.String())
	q.Add("nextToken", f.nextToken)
	return q
}

// CreateReportSpecification Information required to create the report.
type CreateReportSpecification struct {
	// Additional information passed to reports. This varies by report type.
	ReportOptions *map[string]string `json:"reportOptions,omitempty"`
	// The report type.
	ReportType Type `json:"reportType"`
	// The start of a date and time range, in ISO 8601 date time format, used for selecting the data to report. The default is now. The value must be prior to or equal to the current date and time. Not all report types make use of this.
	DataStartTime apis.JsonTimeISO8601 `json:"dataStartTime,omitempty"`
	// The end of a date and time range, in ISO 8601 date time format, used for selecting the data to report. The default is now. The value must be prior to or equal to the current date and time. Not all report types make use of this.
	DataEndTime apis.JsonTimeISO8601 `json:"dataEndTime,omitempty"`
	// A list of marketplace identifiers. The report document's contents will contain data for all of the specified marketplaces, unless the report type indicates otherwise.
	MarketplaceIDs []constants.MarketplaceID `json:"marketplaceIds"`
}

// GetReportDocumentResponse Response schema.
type GetReportDocumentResponse struct {
	ReportDocument
}

// CreateReportResponse Response schema.
type CreateReportResponse struct {
	// The identifier for the report. This identifier is unique only in combination with a seller ID.
	ReportID string `json:"reportId"`
}

// GetReportsResponse The response for the getReports operation.
type GetReportsResponse struct {
	// A list of reports.
	Reports []ReportModel `json:"reports"`
	// Returned when the number of results exceeds pageSize. To get the next page of results, call getReports with this token as the only parameter.
	NextToken *string `json:"nextToken,omitempty"`
}

// ReportDocument Information required for the report document.
type ReportDocument struct {
	// The identifier for the report document. This identifier is unique only in combination with a seller ID.
	ReportDocumentID string `json:"reportDocumentId"`
	// A presigned URL for the report document. This URL expires after 5 minutes.
	Url string `json:"url"`
	// If present, the report document contents have been compressed with the provided algorithm.
	CompressionAlgorithm *string `json:"compressionAlgorithm,omitempty"`
}

// GetReportResponse The response for the getReports operation.
type GetReportResponse struct {
	ReportModel
}

// ReportSchedule Detailed information about a report schedule.
type ReportSchedule struct {
	// The identifier for the report schedule. This identifier is unique only in combination with a seller ID.
	ReportScheduleID string `json:"reportScheduleId"`
	// The report type.
	ReportType Type `json:"reportType"`
	// A list of marketplace identifiers. The report document's contents will contain data for all of the specified marketplaces, unless the report type indicates otherwise.
	MarketplaceIDs []constants.MarketplaceID `json:"marketplaceIds,omitempty"`
	// Additional information passed to reports. This varies by report type.
	ReportOptions *map[string]string `json:"reportOptions,omitempty"`
	// An ISO 8601 period value that indicates how often a report should be created.
	Period string `json:"period"`
	// The date and time when the schedule will create its next report, in ISO 8601 date time format.
	NextReportCreationTime apis.JsonTimeISO8601 `json:"nextReportCreationTime,omitempty"`
}

// GetReportScheduleResponse Response schema.
type GetReportScheduleResponse struct {
	ReportSchedule
}

// ReportScheduleList A list of report schedules.
type ReportScheduleList struct {
	ReportSchedules []ReportSchedule `json:"reportSchedules"`
}

// CreateReportScheduleResponse Response schema.
type CreateReportScheduleResponse struct {
	// The identifier for the report schedule. This identifier is unique only in combination with a seller ID.
	ReportScheduleID string `json:"reportScheduleId"`
}

// CreateReportScheduleSpecification struct for CreateReportScheduleSpecification
type CreateReportScheduleSpecification struct {
	// The report type.
	ReportType Type `json:"reportType"`
	// A list of marketplace identifiers for the report schedule.
	MarketplaceIDs []constants.MarketplaceID `json:"marketplaceIds"`
	// Additional information passed to reports. This varies by report type.
	ReportOptions *map[string]string `json:"reportOptions,omitempty"`
	// One of a set of predefined ISO 8601 periods that specifies how often a report should be created.
	Period string `json:"period"`
	// The date and time when the schedule will create its next report, in ISO 8601 date time format.
	NextReportCreationTime apis.JsonTimeISO8601 `json:"nextReportCreationTime,omitempty"`
}
