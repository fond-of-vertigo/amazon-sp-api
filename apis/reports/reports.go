package reports

import (
	"encoding/json"
	"fmt"
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"strings"
)

type Report struct {
	HttpClient apis.Doer
}

// GetReports returns report details for the reports that match the filters that you specify.
//
// filter are optional and can be set to nil
func (r *Report) GetReports(filter GetReportFilter) (response *GetReportsResponse, err error) {
	if filter.pageSize < 1 {
		filter.pageSize = 10
	}
	params := apis.APICall{}
	params.Method = "GET"
	params.APIPath = config.routePrefix() + "/reports"
	params.QueryParams = filter.GetQuery()
	return apis.CallAPIWithResponseType[GetReportsResponse](params, r.HttpClient)
}

// CreateReport creates a report and returns the reportID.
func (r *Report) CreateReport(specification CreateReportSpecification) (resp *CreateReportResponse, err error) {
	params := apis.APICall{}
	params.Method = "POST"
	params.APIPath = config.routePrefix() + "/reports"
	params.Body, err = json.Marshal(specification)
	if err != nil {
		return nil, err
	}
	return apis.CallAPIWithResponseType[CreateReportResponse](params, r.HttpClient)
}

// GetReport returns report details (including the reportDocumentId, if available) for the report that you specify.
func (r *Report) GetReport(reportId string) (*ReportModel, error) {
	params := apis.APICall{}
	params.Method = "GET"
	params.APIPath = config.routePrefix() + "/reports/" + reportId
	return apis.CallAPIWithResponseType[ReportModel](params, r.HttpClient)
}

// CancelReport returns report schedule details that match the filters that you specify.
//
// reportTypes is list of report types used to filter report schedules. This is optional can can be nil.
func (r *Report) CancelReport(reportId string) error {
	params := apis.APICall{}
	params.Method = "DELETE"
	params.APIPath = config.routePrefix() + "/reports/" + reportId
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
	params.Method = "GET"
	params.APIPath = config.routePrefix() + "/schedules"
	params.QueryParams.Add("reportTypes", strings.Join(reportTypes, ","))
	return apis.CallAPIWithResponseType[GetReportsResponse](params, r.HttpClient)
}

// CreateReportSchedule creates a report schedule.
// If a report schedule with the same report type and marketplace IDs already exists,
// it will be cancelled and replaced with this one.
func (r *Report) CreateReportSchedule(specification CreateReportScheduleSpecification) (resp *CreateReportScheduleResponse, err error) {
	params := apis.APICall{}
	params.Method = "POST"
	params.APIPath = config.routePrefix() + "/schedules"
	params.Body, err = json.Marshal(specification)
	if err != nil {
		return nil, err
	}
	return apis.CallAPIWithResponseType[CreateReportScheduleResponse](params, r.HttpClient)
}

// GetReportSchedule returns report schedule details for the report schedule that you specify.
func (r *Report) GetReportSchedule(reportScheduleId string) (*ReportSchedule, error) {
	params := apis.APICall{}
	params.Method = "GET"
	params.APIPath = config.routePrefix() + "/schedules/" + reportScheduleId
	return apis.CallAPIWithResponseType[ReportSchedule](params, r.HttpClient)
}

// CancelReportSchedule cancels the report schedule that you specify.
func (r *Report) CancelReportSchedule(reportScheduleId string) error {
	params := apis.APICall{}
	params.Method = "DELETE"
	params.APIPath = config.routePrefix() + "/schedules/" + reportScheduleId
	return apis.CallAPIIgnoreResponse(params, r.HttpClient)
}

// GetReportDocument returns the information required for retrieving a report document's contents.
func (r *Report) GetReportDocument(reportDocumentId string) (*ReportDocument, error) {
	params := apis.APICall{}
	params.Method = "GET"
	params.APIPath = config.routePrefix() + "/documents/" + reportDocumentId
	return apis.CallAPIWithResponseType[ReportDocument](params, r.HttpClient)
}
