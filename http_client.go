package selling_partner_api

import "net/http"

type httpClient struct {
	HttpClient   *http.Client
	TokenUpdater *TokenUpdater
}

func (h *httpClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-Amz-Access-Token", h.TokenUpdater.GetAccessToken())
	return h.HttpClient.Do(req)
}
