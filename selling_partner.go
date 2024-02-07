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
	ClientID      string
	ClientSecret  string
	RefreshToken  string
	Endpoint      constants.Endpoint
	Log           logger.Logger
	HTTPClient    *http.Client
	UnmarshalFunc httpx.UnmarshalFunc
}

type Client struct {
	httpClient  *httpx.Client
	FinancesAPI *finances.API
	FeedsAPI    *feeds.API
	OrdersAPI   *orders.API
	ReportsAPI  *reports.API
	TokenAPI    *tokens.API
}

// Close stops the TokenUpdater thread
func (s *Client) Close() {
	s.httpClient.Close()
}

func NewClient(config Config) (*Client, error) {
	hc := config.HTTPClient
	if config.HTTPClient == nil {
		hc = http.DefaultClient
	}

	clientConfig := httpx.ClientConfig{
		HTTPClient:    hc,
		Endpoint:      config.Endpoint,
		UnmarshalFunc: config.UnmarshalFunc,
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
		httpClient:  httpxClient,
		FinancesAPI: finances.NewAPI(httpxClient),
		FeedsAPI:    feeds.NewAPI(httpxClient),
		OrdersAPI:   orders.NewAPI(httpxClient),
		ReportsAPI:  reports.NewAPI(httpxClient),
		TokenAPI:    tokens.NewAPI(httpxClient),
	}, nil
}
