package reports

import (
	"encoding/json"
	"fmt"
	"go/types"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/httpx"
)

const pathPrefix = "/reports/2021-06-30"

type API struct {
	httpClient *httpx.Client
}

func NewAPI(httpClient *httpx.Client) *API {
	return &API{
		httpClient: httpClient,
	}
}

// GetReports returns report details for the reports that match the filters that you specify.
// filter are optional and can be set to nil
func (r *API) GetReports(filter *GetReportsFilter) (*apis.CallResponse[GetReportsResponse], error) {
	if filter.PageSize < 1 {
		filter.PageSize = 10
	}
	return apis.NewCall[GetReportsResponse](http.MethodGet, pathPrefix+"/reports").
		WithQueryParams(filter.GetQuery()).
		WithParseErrorListOnError(true).
		WithRateLimit(0.0222, time.Second).
		Execute(r.httpClient)
}

// CreateReport creates a report and returns the reportID.
func (r *API) CreateReport(specification *CreateReportSpecification) (*apis.CallResponse[CreateReportResponse], error) {
	body, err := json.Marshal(specification)
	if err != nil {
		return nil, err
	}
	return apis.NewCall[CreateReportResponse](http.MethodPost, pathPrefix+"/reports").
		WithBody(body).
		WithParseErrorListOnError(true).
		WithRateLimit(0.0167, time.Second).
		Execute(r.httpClient)
}

// GetReport returns report details (including the reportDocumentID, if available) for the report that you specify.
func (r *API) GetReport(reportID string) (*apis.CallResponse[GetReportResponse], error) {
	return apis.NewCall[GetReportResponse](http.MethodGet, pathPrefix+"/reports/"+reportID).
		WithParseErrorListOnError(true).
		WithRateLimit(2.0, time.Second).
		Execute(r.httpClient)
}

// CancelReport returns report schedule details that match the filters that you specify.
// reportTypes is list of report types used to filter report schedules. This is optional can can be nil.
func (r *API) CancelReport(reportID string) error {
	_, err := apis.NewCall[types.Nil](http.MethodDelete, pathPrefix+"/reports/"+reportID).
		WithRateLimit(0.0222, time.Second).
		Execute(r.httpClient)
	return err
}

// GetReportSchedules returns report schedule details that match the filters that you specify.
// reportTypes is list of report types used to filter report schedules. This is optional can can be nil.
func (r *API) GetReportSchedules(reportTypes []string) (*apis.CallResponse[GetReportsResponse], error) {
	if len(reportTypes) > 10 {
		return nil, fmt.Errorf("reportTypes cannot contain more than 10 reportTypes")
	}
	params := url.Values{}
	params.Add("reportTypes", strings.Join(reportTypes, ","))
	return apis.NewCall[GetReportsResponse](http.MethodGet, pathPrefix+"/schedules").
		WithQueryParams(params).
		WithParseErrorListOnError(true).
		WithRateLimit(0.0222, time.Second).
		Execute(r.httpClient)
}

// CreateReportSchedule creates a report schedule.
// If a report schedule with the same report type and marketplace IDs already exists,
// it will be cancelled and replaced with this one.
func (r *API) CreateReportSchedule(specification *CreateReportScheduleSpecification) (*apis.CallResponse[CreateReportScheduleResponse], error) {
	body, err := json.Marshal(specification)
	if err != nil {
		return nil, err
	}
	return apis.NewCall[CreateReportScheduleResponse](http.MethodPost, pathPrefix+"/schedules").
		WithBody(body).
		WithParseErrorListOnError(true).
		WithRateLimit(0.0222, time.Second).
		Execute(r.httpClient)
}

// GetReportSchedule returns report schedule details for the report schedule that you specify.
func (r *API) GetReportSchedule(reportScheduleID string) (*apis.CallResponse[GetReportScheduleResponse], error) {
	return apis.NewCall[GetReportScheduleResponse](http.MethodGet, pathPrefix+"/schedules/"+reportScheduleID).
		WithParseErrorListOnError(true).
		WithRateLimit(0.0222, time.Second).
		Execute(r.httpClient)
}

// CancelReportSchedule cancels the report schedule that you specify.
func (r *API) CancelReportSchedule(reportScheduleID string) error {
	_, err := apis.NewCall[types.Nil](http.MethodDelete, pathPrefix+"/schedules/"+reportScheduleID).
		WithRateLimit(0.0222, time.Second).
		Execute(r.httpClient)
	return err
}

// GetReportDocument returns the information required for retrieving a report document's contents.
// a restrictedDataToken is optional and may be passed to receive Personally Identifiable Information (PII).
func (r *API) GetReportDocument(reportDocumentID string, restrictedDataToken *string) (*apis.CallResponse[GetReportDocumentResponse], error) {
	return apis.NewCall[GetReportDocumentResponse](http.MethodGet, pathPrefix+"/documents/"+reportDocumentID).
		WithRestrictedDataToken(restrictedDataToken).
		WithParseErrorListOnError(true).
		WithRateLimit(0.0167, time.Second).
		Execute(r.httpClient)
}
