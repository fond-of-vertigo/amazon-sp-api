package reports

import (
	"encoding/json"
	"fmt"
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/httpx"
	"go/types"
	"net/http"
	"net/url"
	"strings"
)

const pathPrefix = "/reports/2021-06-30"

type API interface {
	// GetReports returns report details for the reports that match the filters that you specify.
	// filter are optional and can be set to nil
	GetReports(filter *GetReportsFilter) (*apis.CallResponse[GetReportsResponse], error)
	// CreateReport creates a report and returns the reportID.
	CreateReport(specification *CreateReportSpecification) (*apis.CallResponse[CreateReportResponse], error)
	// GetReport returns report details (including the reportDocumentID, if available) for the report that you specify.
	GetReport(reportID string) (*apis.CallResponse[GetReportResponse], error)
	// CancelReport returns report schedule details that match the filters that you specify.
	// reportTypes is list of report types used to filter report schedules. This is optional can can be nil.
	CancelReport(reportID string) error
	// GetReportSchedules returns report schedule details that match the filters that you specify.
	// reportTypes is list of report types used to filter report schedules. This is optional can can be nil.
	GetReportSchedules(reportTypes []string) (*apis.CallResponse[GetReportsResponse], error)
	// CreateReportSchedule creates a report schedule.
	// If a report schedule with the same report type and marketplace IDs already exists,
	// it will be cancelled and replaced with this one.
	CreateReportSchedule(specification *CreateReportScheduleSpecification) (*apis.CallResponse[CreateReportScheduleResponse], error)
	// GetReportSchedule returns report schedule details for the report schedule that you specify.
	GetReportSchedule(reportScheduleID string) (*apis.CallResponse[GetReportScheduleResponse], error)
	// CancelReportSchedule cancels the report schedule that you specify.
	CancelReportSchedule(reportScheduleID string) error
	// GetReportDocument returns the information required for retrieving a report document's contents.
	// a restrictedDataToken is optional and may be passed to receive Personally Identifiable Information (PII).
	GetReportDocument(reportDocumentID string, restrictedDataToken *string) (*apis.CallResponse[GetReportDocumentResponse], error)
}
type api struct {
	HttpClient httpx.Client
}

func NewAPI(httpClient httpx.Client) API {
	return &api{
		HttpClient: httpClient,
	}
}

func (r *api) GetReports(filter *GetReportsFilter) (*apis.CallResponse[GetReportsResponse], error) {
	if filter.pageSize < 1 {
		filter.pageSize = 10
	}
	return apis.NewCall[GetReportsResponse](http.MethodGet, pathPrefix+"/reports").
		WithQueryParams(filter.GetQuery()).
		WithParseErrorListOnError(true).
		Execute(r.HttpClient)
}

func (r *api) CreateReport(specification *CreateReportSpecification) (*apis.CallResponse[CreateReportResponse], error) {
	body, err := json.Marshal(specification)
	if err != nil {
		return nil, err
	}
	return apis.NewCall[CreateReportResponse](http.MethodPost, pathPrefix+"/reports").
		WithBody(body).
		WithParseErrorListOnError(true).
		Execute(r.HttpClient)
}

func (r *api) GetReport(reportID string) (*apis.CallResponse[GetReportResponse], error) {
	return apis.NewCall[GetReportResponse](http.MethodGet, pathPrefix+"/reports/"+reportID).
		WithParseErrorListOnError(true).
		Execute(r.HttpClient)
}

func (r *api) CancelReport(reportID string) error {
	_, err := apis.NewCall[types.Nil](http.MethodDelete, pathPrefix+"/reports/"+reportID).
		Execute(r.HttpClient)
	return err
}

func (r *api) GetReportSchedules(reportTypes []string) (*apis.CallResponse[GetReportsResponse], error) {
	if len(reportTypes) > 10 {
		return nil, fmt.Errorf("reportTypes cannot contain more than 10 reportTypes")
	}
	params := url.Values{}
	params.Add("reportTypes", strings.Join(reportTypes, ","))
	return apis.NewCall[GetReportsResponse](http.MethodGet, pathPrefix+"/schedules").
		WithQueryParams(params).
		WithParseErrorListOnError(true).
		Execute(r.HttpClient)
}

func (r *api) CreateReportSchedule(specification *CreateReportScheduleSpecification) (*apis.CallResponse[CreateReportScheduleResponse], error) {
	body, err := json.Marshal(specification)
	if err != nil {
		return nil, err
	}
	return apis.NewCall[CreateReportScheduleResponse](http.MethodPost, pathPrefix+"/schedules").
		WithBody(body).
		WithParseErrorListOnError(true).
		Execute(r.HttpClient)
}

func (r *api) GetReportSchedule(reportScheduleID string) (*apis.CallResponse[GetReportScheduleResponse], error) {
	return apis.NewCall[GetReportScheduleResponse](http.MethodGet, pathPrefix+"/schedules/"+reportScheduleID).
		WithParseErrorListOnError(true).
		Execute(r.HttpClient)
}

func (r *api) CancelReportSchedule(reportScheduleID string) error {
	_, err := apis.NewCall[types.Nil](http.MethodDelete, pathPrefix+"/schedules/"+reportScheduleID).
		Execute(r.HttpClient)
	return err
}

func (r *api) GetReportDocument(reportDocumentID string, restrictedDataToken *string) (*apis.CallResponse[GetReportDocumentResponse], error) {
	return apis.NewCall[GetReportDocumentResponse](http.MethodGet, pathPrefix+"/documents/"+reportDocumentID).
		WithRestrictedDataToken(restrictedDataToken).
		WithParseErrorListOnError(true).
		Execute(r.HttpClient)
}
