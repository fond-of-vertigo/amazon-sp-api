package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	sp_api "github.com/fond-of-vertigo/amazon-sp-api"
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/apis/reports"
	"github.com/fond-of-vertigo/amazon-sp-api/apis/tokens"
	"github.com/fond-of-vertigo/amazon-sp-api/constants"
	"github.com/fond-of-vertigo/logger"
)

const PollingDelay = time.Second * 5

func main() {
	log := logger.New(logger.LvlDebug)
	c := sp_api.Config{
		ClientID:     mustGetenv("AMZN_CLIENT_ID"),
		ClientSecret: mustGetenv("AMZN_CLIENT_SECRET"),
		RefreshToken: mustGetenv("AMZN_REFRESH_TOKEN"),
		Endpoint:     constants.Europe,
		Log:          log,
	}

	client, err := sp_api.NewClient(c)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	now := time.Now()
	from := now.Add(-24 * time.Hour * 7)
	spec := &reports.CreateReportSpecification{
		ReportType:     "GET_SALES_AND_TRAFFIC_REPORT", // reports.FBAAmazonFulfilledShipmentsInvoicing,
		DataStartTime:  apis.JsonTimeISO8601{Time: from},
		DataEndTime:    apis.JsonTimeISO8601{Time: now},
		MarketplaceIDs: []constants.MarketplaceID{constants.Germany},
	}
	reportID, err := RequestReport(log, client, spec)
	if err != nil {
		log.Errorf("Report could not be requested: %w - %v", err, err)
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

func RequestReport(log logger.Logger, client *sp_api.Client, specification *reports.CreateReportSpecification) (string, error) {
	createdReportResp, err := client.ReportsAPI.CreateReport(specification)
	if err != nil {
		return "", err
	}
	if createdReportResp.IsError() {
		return "", fmt.Errorf("creating report failed with status %v. ErrorList: %v", createdReportResp.Status, createdReportResp.ErrorList)
	}
	log.Infof("API with ID=%s was queued..", createdReportResp.ResponseBody.ReportID)
	return createdReportResp.ResponseBody.ReportID, nil
}

func WaitForReport(log logger.Logger, client *sp_api.Client, reportID string) (*reports.GetReportResponse, error) {
	var getReportResp *apis.CallResponse[reports.GetReportResponse]
	var err error
	for getReportResp == nil || !getReportResp.ResponseBody.ProcessingStatus.IsDone() {
		getReportResp, err = client.ReportsAPI.GetReport(reportID)
		if err != nil {
			return nil, err
		}
		if getReportResp.IsError() {
			return nil, fmt.Errorf("waiting for report(id: %v) failed with status %v. ErrorList: %v", reportID, getReportResp.Status, getReportResp.ErrorList)
		}
		log.Infof("API with ID=%s has processingStatus=%s", getReportResp.ResponseBody.ReportID, getReportResp.ResponseBody.ProcessingStatus)
		log.Infof("Wait %v seconds", PollingDelay.Seconds())
		time.Sleep(PollingDelay)
	}
	return getReportResp.ResponseBody, nil
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
		if tokenResp.IsError() {
			return nil, fmt.Errorf("create RestrictedDataToken failed with status %v. ErrorList: %v", tokenResp.Status, tokenResp.ErrorList)
		}
		log.Infof("Fetched RDT=%s", tokenResp.ResponseBody.RestrictedDataToken)
		rdt = tokenResp.ResponseBody.RestrictedDataToken
	}

	getRepDocResp, err := client.ReportsAPI.GetReportDocument(*getReport.ReportDocumentID, rdt)
	if err != nil {
		return nil, err
	}
	if getRepDocResp.IsError() {
		return nil, fmt.Errorf("create GetReportDocument request failed with status %v. ErrorList: %v", getRepDocResp.Status, getRepDocResp.ErrorList)
	}
	log.Infof("Downloading document ID=%s via URL=%s", getRepDocResp.ResponseBody.ReportDocumentID, getRepDocResp.ResponseBody.Url)

	httpResp, httpErr := http.Get(getRepDocResp.ResponseBody.Url)
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

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("missing env var %s", key))
	}
	return v
}
