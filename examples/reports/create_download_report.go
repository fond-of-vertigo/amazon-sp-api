package main

import (
	amznsp "github.com/fond-of-vertigo/amazon-sp-api"
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/apis/reports"
	"github.com/fond-of-vertigo/amazon-sp-api/apis/tokens"
	"github.com/fond-of-vertigo/logger"
	"io"
	"net/http"
	"time"
)

const PollingDelay = time.Second * 5

func DownloadReport(log logger.Logger, sp *amznsp.SellingPartnerClient, specification reports.CreateReportSpecification, useRDT bool) ([]byte, error) {
	resp, err := sp.Report.CreateReport(specification)
	if err != nil {
		return nil, err
	}
	log.Infof("Report with ID=%s was queued..", resp.ReportID)

	var rm *reports.ReportModel
	for rm == nil || rm.ProcessingStatus != reports.ProcessingStatusDone {
		rm, err = sp.Report.GetReport(resp.ReportID)
		if err != nil {
			return nil, err
		}
		log.Infof("Report with ID=%s has processingStatus=%s", rm.ReportID, rm.ProcessingStatus)
		log.Infof("Wait %v seconds", PollingDelay.Seconds())
		time.Sleep(PollingDelay)
	}
	var rdt *string
	if useRDT {
		log.Infof("Fetching RDT for %s", rm.GetDocumentAPIPath())
		rr := tokens.CreateRestrictedDataTokenRequest{
			RestrictedResources: []tokens.RestrictedResource{
				{
					Method: http.MethodGet,
					Path:   rm.GetDocumentAPIPath(),
				},
			},
		}
		tokenResp, err := sp.Token.CreateRestrictedDataTokenRequest(rr)
		if err != nil {
			return nil, err
		}
		log.Infof("Fetched RDT=%s", tokenResp.RestrictedDataToken)
		rdt = tokenResp.RestrictedDataToken
	}

	doc, err := sp.Report.GetReportDocument(*rm.ReportDocumentID, rdt)
	if err != nil {
		return nil, err
	}
	log.Infof("Downloading document ID=%s via URL=%s", doc.ReportDocumentID, doc.Url)

	httpResp, err := http.Get(doc.Url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

func main() {
	log := logger.New(logger.LvlDebug)
	c := amznsp.Config{
		ClientID:           "EXAMPLE_CLIENTID",
		ClientSecret:       "EXAMPLE_SECRET",
		RefreshToken:       "EXAMPLE_REFRESHTOKEN",
		IAMUserAccessKeyID: "EXAMPLE_ACCESSKEY",
		IAMUserSecretKey:   "EXAMPLE_SECRET",
		Region:             amznsp.AWSRegionEUWest,
		RoleArn:            "EXAMPLE_ROLE",
		Endpoint:           amznsp.EndpointEurope,
		Log:                log,
	}

	sp, err := amznsp.NewSellingPartnerClient(c)
	if err != nil {
		panic(err)
	}
	defer sp.Close()

	now := time.Now()
	from := now.Add(-24 * time.Hour * 7)
	spec := reports.CreateReportSpecification{
		ReportType:     "GET_AMAZON_FULFILLED_SHIPMENTS_DATA_INVOICING",
		DataStartTime:  apis.JsonTimeISO8601{Time: from},
		DataEndTime:    apis.JsonTimeISO8601{Time: now},
		MarketplaceIDs: []reports.MarketplaceID{reports.MarketplaceIDGermany},
	}
	r, err := DownloadReport(log, sp, spec, true)
	if err != nil {
		log.Errorf("Report could not be downloaded: %w", err)
	}
	log.Infof("Report: %s", r)
}
