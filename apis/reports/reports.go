package reports

import (
	"encoding/json"
	"fmt"
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"net/http"
	"strings"
)

const pathPrefix = "/reports/2021-06-30"

type Report struct {
	HttpClient apis.HttpRequestDoer
}

// GetReports returns report details for the reports that match the filters that you specify.
//
// filter are optional and can be set to nil
func (r *Report) GetReports(filter GetReportFilter) (response *GetReportsResponse, err error) {
	if filter.pageSize < 1 {
		filter.pageSize = 10
	}
	params := apis.APICall{}
	params.Method = http.MethodGet
	params.APIPath = pathPrefix + "/reports"
	params.QueryParams = filter.GetQuery()
	return apis.CallAPIWithResponseType[GetReportsResponse](params, r.HttpClient)
}

// CreateReport creates a report and returns the reportID.
func (r *Report) CreateReport(specification CreateReportSpecification) (resp *CreateReportResponse, err error) {
	params := apis.APICall{}
	params.Method = http.MethodPost
	params.APIPath = pathPrefix + "/reports"
	params.Body, err = json.Marshal(specification)
	if err != nil {
		return nil, err
	}
	return apis.CallAPIWithResponseType[CreateReportResponse](params, r.HttpClient)
}

// GetReport returns report details (including the reportDocumentID, if available) for the report that you specify.
func (r *Report) GetReport(reportID string) (*ReportModel, error) {
	params := apis.APICall{}
	params.Method = http.MethodGet
	params.APIPath = pathPrefix + "/reports/" + reportID
	return apis.CallAPIWithResponseType[ReportModel](params, r.HttpClient)
}

// CancelReport returns report schedule details that match the filters that you specify.
//
// reportTypes is list of report types used to filter report schedules. This is optional can can be nil.
func (r *Report) CancelReport(reportID string) error {
	params := apis.APICall{}
	params.Method = http.MethodDelete
	params.APIPath = pathPrefix + "/reports/" + reportID
	return apis.CallAPIIgnoreResponse(params, r.HttpClient)
}

// GetReportSchedules returns report schedule details that match the filters that you specify.
//
// reportTypes is list of report types used to filter report schedules. This is optional can can be nil.
func (r *Report) GetReportSchedules(reportTypes []string) (*GetReportsResponse, error) {
	if len(reportTypes) > 10 {
		return nil, fmt.Errorf("reportTypes cannot contain more than 10 reportTypes")
	}
	params := apis.APICall{}
	params.Method = http.MethodGet
	params.APIPath = pathPrefix + "/schedules"
	params.QueryParams.Add("reportTypes", strings.Join(reportTypes, ","))
	return apis.CallAPIWithResponseType[GetReportsResponse](params, r.HttpClient)
}

// CreateReportSchedule creates a report schedule.
// If a report schedule with the same report type and marketplace IDs already exists,
// it will be cancelled and replaced with this one.
func (r *Report) CreateReportSchedule(specification CreateReportScheduleSpecification) (resp *CreateReportScheduleResponse, err error) {
	params := apis.APICall{}
	params.Method = http.MethodPost
	params.APIPath = pathPrefix + "/schedules"
	params.Body, err = json.Marshal(specification)
	if err != nil {
		return nil, err
	}
	return apis.CallAPIWithResponseType[CreateReportScheduleResponse](params, r.HttpClient)
}

// GetReportSchedule returns report schedule details for the report schedule that you specify.
func (r *Report) GetReportSchedule(reportScheduleID string) (*ReportSchedule, error) {
	params := apis.APICall{}
	params.Method = http.MethodGet
	params.APIPath = pathPrefix + "/schedules/" + reportScheduleID
	return apis.CallAPIWithResponseType[ReportSchedule](params, r.HttpClient)
}

// CancelReportSchedule cancels the report schedule that you specify.
func (r *Report) CancelReportSchedule(reportScheduleID string) error {
	params := apis.APICall{}
	params.Method = http.MethodDelete
	params.APIPath = pathPrefix + "/schedules/" + reportScheduleID
	return apis.CallAPIIgnoreResponse(params, r.HttpClient)
}

// GetReportDocument returns the information required for retrieving a report document's contents.
// a restrictedDataToken is optional and may be passed to receive Personally Identifiable Information (PII).
func (r *Report) GetReportDocument(reportDocumentID string, restrictedDataToken *string) (*ReportDocument, error) {
	params := apis.APICall{}
	params.Method = http.MethodGet
	params.APIPath = pathPrefix + "/documents/" + reportDocumentID
	params.RestrictedDataToken = restrictedDataToken
	return apis.CallAPIWithResponseType[ReportDocument](params, r.HttpClient)
}
