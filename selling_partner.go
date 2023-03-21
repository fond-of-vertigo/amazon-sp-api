package sp_api

import (
	"net/http"

	"github.com/fond-of-vertigo/amazon-sp-api/apis/feeds"
	"github.com/fond-of-vertigo/amazon-sp-api/apis/finances"
	"github.com/fond-of-vertigo/amazon-sp-api/apis/orders"
	"github.com/fond-of-vertigo/amazon-sp-api/apis/reports"
	"github.com/fond-of-vertigo/amazon-sp-api/apis/tokens"
	"github.com/fond-of-vertigo/amazon-sp-api/constants"
	"github.com/fond-of-vertigo/amazon-sp-api/httpx"
	"github.com/fond-of-vertigo/logger"
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
	HTTPClient         *http.Client
}

type Client struct {
	hTTPClient  *httpx.Client
	FinancesAPI *finances.API
	FeedsAPI    *feeds.API
	OrdersAPI   *orders.API
	ReportsAPI  *reports.API
	TokenAPI    *tokens.API
}

// Close stops the TokenUpdater thread
func (s *Client) Close() {
	s.hTTPClient.Close()
}

func NewClient(config Config) (*Client, error) {
	hc := config.HTTPClient
	if config.HTTPClient == nil {
		hc = http.DefaultClient
	}

	clientConfig := httpx.ClientConfig{
		HTTPClient:         hc,
		Endpoint:           config.Endpoint,
		IAMUserAccessKeyID: config.IAMUserAccessKeyID,
		IAMUserSecretKey:   config.IAMUserSecretKey,
		Region:             config.Region,
		RoleArn:            config.RoleArn,
		TokenUpdaterConfig: httpx.TokenUpdaterConfig{
			RefreshToken: config.RefreshToken,
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			HTTPClient:   hc,
			Logger:       config.Log,
		},
	}

	httpxClient, err := httpx.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	return &Client{
		hTTPClient:  httpxClient,
		FinancesAPI: finances.NewAPI(httpxClient),
		FeedsAPI:    feeds.NewAPI(httpxClient),
		OrdersAPI:   orders.NewAPI(httpxClient),
		ReportsAPI:  reports.NewAPI(httpxClient),
		TokenAPI:    tokens.NewAPI(httpxClient),
	}, nil
}
