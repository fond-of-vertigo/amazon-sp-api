package selling_partner_api

import (
	"github.com/fond-of-vertigo/amazon-sp-api/apis/reports"
	"github.com/fond-of-vertigo/amazon-sp-api/apis/tokens"
	"github.com/fond-of-vertigo/amazon-sp-api/constants"
	"github.com/fond-of-vertigo/logger"
	"net/http"
)

type Config struct {
	ClientID           string
	ClientSecret       string
	RefreshToken       string
	IAMUserAccessKeyID string
	IAMUserSecretKey   string
	Region             constants.Region
	RoleArn            string
	Endpoint           constants.Endpoint
	Log                logger.Logger
}

type SellingPartnerClient struct {
	quitSignal chan bool
	Report     reports.Report
	Token      tokens.Token
}

// Close stops the TokenUpdater thread
func (s *SellingPartnerClient) Close() {
	s.quitSignal <- true
}

func NewSellingPartnerClient(config Config) (*SellingPartnerClient, error) {
	quitSignal := make(chan bool)

	t := NewTokenUpdater(TokenUpdaterConfig{
		RefreshToken: config.RefreshToken,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Logger:       config.Log,
	})
	if err := t.RunInBackground(); err != nil {
		return nil, err
	}

	h := HttpClientConfig{
		client:             &http.Client{},
		Endpoint:           config.Endpoint,
		TokenUpdater:       t,
		IAMUserAccessKeyID: config.IAMUserAccessKeyID,
		IAMUserSecretKey:   config.IAMUserSecretKey,
		Region:             config.Region,
		RoleArn:            config.RoleArn,
	}
	httpClient, err := NewHttpClient(h)
	if err != nil {
		return nil, err
	}

	return &SellingPartnerClient{
		quitSignal: quitSignal,
		Report:     reports.Report{HttpClient: httpClient},
		Token:      tokens.Token{HttpClient: httpClient},
	}, nil
}
