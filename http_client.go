package selling_partner_api

import "net/http"

type httpClient struct {
	HttpClient     *http.Client
	TokenRefresher *TokenRefresher
}

func (h *httpClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-Amz-Access-Token", h.TokenRefresher.GetAccessToken())
	return h.HttpClient.Do(req)
}
