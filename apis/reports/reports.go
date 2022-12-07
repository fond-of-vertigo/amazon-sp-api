package reports

import (
	"encoding/json"
	"fmt"
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"go/types"
	"net/http"
	"net/url"
	"strings"
)

const pathPrefix = "/reports/2021-06-30"

type API interface {
	// GetReports returns report details for the reports that match the filters that you specify.
	// filter are optional and can be set to nil
	GetReports(filter *GetReportsFilter) (*GetReportsResponse, apis.CallError)
	// CreateReport creates a report and returns the reportID.
	CreateReport(specification *CreateReportSpecification) (*CreateReportResponse, apis.CallError)
	// GetReport returns report details (including the reportDocumentID, if available) for the report that you specify.
	GetReport(reportID string) (*GetReportResponse, apis.CallError)
	// CancelReport returns report schedule details that match the filters that you specify.
	// reportTypes is list of report types used to filter report schedules. This is optional can can be nil.
	CancelReport(reportID string) apis.CallError
	// GetReportSchedules returns report schedule details that match the filters that you specify.
	// reportTypes is list of report types used to filter report schedules. This is optional can can be nil.
	GetReportSchedules(reportTypes []string) (*GetReportsResponse, apis.CallError)
	// CreateReportSchedule creates a report schedule.
	// If a report schedule with the same report type and marketplace IDs already exists,
	// it will be cancelled and replaced with this one.
	CreateReportSchedule(specification *CreateReportScheduleSpecification) (*CreateReportScheduleResponse, apis.CallError)
	// GetReportSchedule returns report schedule details for the report schedule that you specify.
	GetReportSchedule(reportScheduleID string) (*GetReportScheduleResponse, apis.CallError)
	// CancelReportSchedule cancels the report schedule that you specify.
	CancelReportSchedule(reportScheduleID string) apis.CallError
	// GetReportDocument returns the information required for retrieving a report document's contents.
	// a restrictedDataToken is optional and may be passed to receive Personally Identifiable Information (PII).
	GetReportDocument(reportDocumentID string, restrictedDataToken *string) (*GetReportDocumentResponse, apis.CallError)
}
type api struct {
	HttpClient apis.HttpRequestDoer
}

func NewAPI(httpClient apis.HttpRequestDoer) API {
	return &api{
		HttpClient: httpClient,
	}
}

func (r *api) GetReports(filter *GetReportsFilter) (*GetReportsResponse, apis.CallError) {
	if filter.pageSize < 1 {
		filter.pageSize = 10
	}
	return apis.NewCall[GetReportsResponse](http.MethodGet, pathPrefix+"/reports").
		WithQueryParams(filter.GetQuery()).
		Execute(r.HttpClient)
}

func (r *api) CreateReport(specification *CreateReportSpecification) (*CreateReportResponse, apis.CallError) {
	body, err := json.Marshal(specification)
	if err != nil {
		return nil, apis.NewError(err)
	}
	return apis.NewCall[CreateReportResponse](http.MethodPost, pathPrefix+"/reports").
		WithBody(body).
		Execute(r.HttpClient)
}

func (r *api) GetReport(reportID string) (*GetReportResponse, apis.CallError) {
	return apis.NewCall[GetReportResponse](http.MethodGet, pathPrefix+"/reports/"+reportID).
		Execute(r.HttpClient)
}

func (r *api) CancelReport(reportID string) apis.CallError {
	_, err := apis.NewCall[types.Nil](http.MethodDelete, pathPrefix+"/reports/"+reportID).
		Execute(r.HttpClient)
	return err
}

func (r *api) GetReportSchedules(reportTypes []string) (*GetReportsResponse, apis.CallError) {
	if len(reportTypes) > 10 {
		return nil, apis.NewError(fmt.Errorf("reportTypes cannot contain more than 10 reportTypes"))
	}
	params := url.Values{}
	params.Add("reportTypes", strings.Join(reportTypes, ","))
	return apis.NewCall[GetReportsResponse](http.MethodGet, pathPrefix+"/schedules").
		WithQueryParams(params).
		Execute(r.HttpClient)
}

func (r *api) CreateReportSchedule(specification *CreateReportScheduleSpecification) (*CreateReportScheduleResponse, apis.CallError) {
	body, err := json.Marshal(specification)
	if err != nil {
		return nil, apis.NewError(err)
	}
	return apis.NewCall[CreateReportScheduleResponse](http.MethodPost, pathPrefix+"/schedules").
		WithBody(body).
		Execute(r.HttpClient)
}

func (r *api) GetReportSchedule(reportScheduleID string) (*GetReportScheduleResponse, apis.CallError) {
	return apis.NewCall[GetReportScheduleResponse](http.MethodGet, pathPrefix+"/schedules/"+reportScheduleID).
		Execute(r.HttpClient)
}

func (r *api) CancelReportSchedule(reportScheduleID string) apis.CallError {
	_, err := apis.NewCall[types.Nil](http.MethodDelete, pathPrefix+"/schedules/"+reportScheduleID).
		Execute(r.HttpClient)
	return err
}

func (r *api) GetReportDocument(reportDocumentID string, restrictedDataToken *string) (*GetReportDocumentResponse, apis.CallError) {
	return apis.NewCall[GetReportDocumentResponse](http.MethodGet, pathPrefix+"/documents/"+reportDocumentID).
		WithRestrictedDataToken(restrictedDataToken).
		Execute(r.HttpClient)
}
