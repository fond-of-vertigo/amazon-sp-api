package httpx

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/fond-of-vertigo/amazon-sp-api/constants"
)

type ClientConfig struct {
	HTTPClient         HTTPRequester
	TokenUpdaterConfig TokenUpdaterConfig
	Endpoint           constants.Endpoint
	UnmarshalFunc      UnmarshalFunc
}

func NewClient(config ClientConfig) (c *Client, err error) {
	c = &Client{
		httpClient:    config.HTTPClient,
		endpoint:      config.Endpoint,
		unmarshalFunc: config.UnmarshalFunc,
	}

	if c.unmarshalFunc == nil {
		c.unmarshalFunc = unmarshalBody
	}

	c.tokenUpdater = newTokenUpdater(config.TokenUpdaterConfig)
	if c.tokenUpdaterCancelFunc, err = c.tokenUpdater.RunInBackground(); err != nil {
		return nil, err
	}

	return c, nil
}

type UnmarshalFunc func(resp *http.Response, into any) error

type Client struct {
	tokenUpdater           tokenUpdater
	tokenUpdaterCancelFunc func()
	httpClient             HTTPRequester
	endpoint               constants.Endpoint
	unmarshalFunc          UnmarshalFunc
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

func (h *Client) Unmarshal(resp *http.Response, into any) error {
	return h.unmarshalFunc(resp, into)
}

func unmarshalBody(resp *http.Response, into any) error {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(bodyBytes) == 0 {
		return nil
	}
	return json.Unmarshal(bodyBytes, into)
}

func (h *Client) addAccessTokenToHeader(req *http.Request) {
	if req.Header.Get(constants.AccessTokenHeader) == "" {
		req.Header.Add(constants.AccessTokenHeader, h.tokenUpdater.GetAccessToken())
	}
}
