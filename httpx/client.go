package httpx

import (
	"github.com/fond-of-vertigo/amazon-sp-api/constants"
	"io"
	"net/http"
)

type ClientConfig struct {
	HTTPClient         HTTPRequester
	TokenUpdaterConfig TokenUpdaterConfig
	Endpoint           constants.Endpoint
}

func NewClient(config ClientConfig) (c *Client, err error) {
	c = &Client{
		httpClient: config.HTTPClient,
		endpoint:   config.Endpoint,
	}

	c.tokenUpdater = newTokenUpdater(config.TokenUpdaterConfig)
	if c.tokenUpdaterCancelFunc, err = c.tokenUpdater.RunInBackground(); err != nil {
		return nil, err
	}

	return c, nil
}

type Client struct {
	tokenUpdater           tokenUpdater
	tokenUpdaterCancelFunc func()
	httpClient             HTTPRequester
	endpoint               constants.Endpoint
}

type HTTPRequester interface {
	Do(req *http.Request) (*http.Response, error)
	Post(url string, bodyType string, body io.Reader) (*http.Response, error)
}

type tokenUpdater interface {
	GetAccessToken() string
	RunInBackground() (cancel func(), err error)
}

func (h *Client) Do(req *http.Request) (*http.Response, error) {
	h.addAccessTokenToHeader(req)

	return h.httpClient.Do(req)
}

func (h *Client) GetEndpoint() constants.Endpoint {
	return h.endpoint
}

func (h *Client) Close() {
	h.tokenUpdaterCancelFunc()
}

func (h *Client) addAccessTokenToHeader(req *http.Request) {
	if req.Header.Get(constants.AccessTokenHeader) == "" {
		req.Header.Add(constants.AccessTokenHeader, h.tokenUpdater.GetAccessToken())
	}
}
