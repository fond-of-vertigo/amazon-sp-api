package sp_api

import (
	"github.com/fond-of-vertigo/amazon-sp-api/apis/orders"
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

type Client struct {
	quitSignal chan bool
	OrdersAPI  orders.API
	ReportsAPI reports.API
	TokenAPI   tokens.API
}

// Close stops the TokenUpdater thread
func (s *Client) Close() {
	s.quitSignal <- true
}

func NewClient(config Config) (*Client, error) {
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

	return &Client{
		quitSignal: quitSignal,
		OrdersAPI:  orders.NewAPI(httpClient),
		ReportsAPI: reports.NewAPI(httpClient),
		TokenAPI:   tokens.NewAPI(httpClient),
	}, nil
}
