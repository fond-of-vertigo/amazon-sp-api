package selling_partner_api

import (
	"github.com/fond-of-vertigo/amazon-sp-api/apis/reports"
	"github.com/fond-of-vertigo/logger"
	"net/http"
	"net/url"
)

type Config struct {
	ClientID           string
	ClientSecret       string
	RefreshToken       string
	IAMUserAccessKeyID string
	IAMUserSecretKey   string
	Region             string
	RoleArn            string
	Endpoint           url.URL
	Log                logger.Logger
}

type SellingPartnerClient struct {
	tokenUpdater *TokenUpdater
	quitSignal   chan bool
	Report       reports.Report
}

// Close stops the TokenUpdater thread
func (s *SellingPartnerClient) Close() {
	s.quitSignal <- true
}

func NewSellingPartnerClient(config Config) (*SellingPartnerClient, error) {
	quitSignal := make(chan bool)

	t, err := NewTokenUpdater(TokenUpdaterConfig{
		RefreshToken: config.RefreshToken,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Logger:       config.Log,
	})
	if err != nil {
		return nil, err
	}
	t.RunInBackground()

	httpClient := httpClient{HttpClient: &http.Client{}, TokenUpdater: t}
	return &SellingPartnerClient{
		tokenUpdater: t,
		quitSignal:   quitSignal,
		Report:       reports.Report{HttpClient: &httpClient},
	}, nil
}
