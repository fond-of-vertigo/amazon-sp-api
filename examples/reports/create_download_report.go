package main

import (
	sp_api "github.com/fond-of-vertigo/amazon-sp-api"
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/apis/reports"
	"github.com/fond-of-vertigo/amazon-sp-api/apis/tokens"
	"github.com/fond-of-vertigo/amazon-sp-api/constants"
	"github.com/fond-of-vertigo/logger"
	"io"
	"net/http"
	"time"
)

const PollingDelay = time.Second * 5

func main() {
	log := logger.New(logger.LvlDebug)
	c := sp_api.Config{
		ClientID:           "EXAMPLE_CLIENTID",
		ClientSecret:       "EXAMPLE_SECRET",
		RefreshToken:       "EXAMPLE_REFRESHTOKEN",
		IAMUserAccessKeyID: "EXAMPLE_ACCESSKEY",
		IAMUserSecretKey:   "EXAMPLE_SECRET",
		Region:             constants.EUWest,
		RoleArn:            "EXAMPLE_ROLE",
		Endpoint:           constants.Europe,
		Log:                log,
	}

	client, err := sp_api.NewClient(c)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	now := time.Now()
	from := now.Add(-24 * time.Hour * 7)
	spec := &reports.CreateReportSpecification{
		ReportType:     reports.FBAAmazonFulfilledShipmentsInvoicing,
		DataStartTime:  apis.JsonTimeISO8601{Time: from},
		DataEndTime:    apis.JsonTimeISO8601{Time: now},
		MarketplaceIDs: []constants.MarketplaceID{constants.Germany},
	}
	reportID, callErr := RequestReport(log, client, spec)
	if callErr != nil {
		log.Errorf("Report could not be requested: %w - %v", callErr, callErr.ErrorList())
		return
	}
	getReport, err := WaitForReport(log, client, reportID)
	if err != nil {
		log.Errorf("Report could not be requested: %w", err)
		log.Errorf("Error while waiting for report(%s): %w", reportID, err)
		return
	}
	r, err := DownloadReport(log, client, getReport, true)
	if err != nil {
		log.Errorf("Report could not be downloaded: %w", err)
		return
	}
	log.Infof("Report data: %s", r)
}

func RequestReport(log logger.Logger, client *sp_api.Client, specification *reports.CreateReportSpecification) (string, apis.CallError) {
	createdReport, err := client.ReportsAPI.CreateReport(specification)
	if err != nil {
		return "", err
	}
	log.Infof("API with ID=%s was queued..", createdReport.ReportID)
	return createdReport.ReportID, nil
}
func WaitForReport(log logger.Logger, client *sp_api.Client, reportID string) (*reports.GetReportResponse, error) {
	var getReport *reports.GetReportResponse
	var err error
	for getReport == nil || !getReport.ProcessingStatus.IsDone() {
		getReport, err = client.ReportsAPI.GetReport(reportID)
		if err != nil {
			return nil, err
		}
		log.Infof("API with ID=%s has processingStatus=%s", getReport.ReportID, getReport.ProcessingStatus)
		log.Infof("Wait %v seconds", PollingDelay.Seconds())
		time.Sleep(PollingDelay)
	}
	return getReport, nil
}
func DownloadReport(log logger.Logger, client *sp_api.Client, getReport *reports.GetReportResponse, useRDT bool) ([]byte, error) {
	var rdt *string
	if useRDT {
		log.Infof("Fetching RDT for %s", getReport.GetDocumentAPIPath())
		rr := &tokens.CreateRestrictedDataTokenRequest{
			RestrictedResources: []tokens.RestrictedResource{
				{
					Method: http.MethodGet,
					Path:   getReport.GetDocumentAPIPath(),
				},
			},
		}
		tokenResp, err := client.TokenAPI.CreateRestrictedDataTokenRequest(rr)
		if err != nil {
			return nil, err
		}
		log.Infof("Fetched RDT=%s", tokenResp.RestrictedDataToken)
		rdt = tokenResp.RestrictedDataToken
	}

	doc, err := client.ReportsAPI.GetReportDocument(*getReport.ReportDocumentID, rdt)
	if err != nil {
		return nil, err
	}
	log.Infof("Downloading document ID=%s via URL=%s", doc.ReportDocumentID, doc.Url)

	httpResp, httpErr := http.Get(doc.Url)
	if httpErr != nil {
		return nil, httpErr
	}
	defer httpResp.Body.Close()

	bodyBytes, httpErr := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, httpErr
	}
	return bodyBytes, nil
}
