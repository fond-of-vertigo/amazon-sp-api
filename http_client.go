package selling_partner_api

import "net/http"

type httpClient struct {
	HttpClient   *http.Client
	TokenUpdater TokenUpdaterInterface
}

func (h *httpClient) Do(req *http.Request) (*http.Response, error) {
	h.addAccessToken(req)
	return h.HttpClient.Do(req)
}

func (h *httpClient) addAccessToken(req *http.Request) {
	if req.Header.Get("X-Amz-Access-Token") == "" {
		req.Header.Add("X-Amz-Access-Token", h.TokenUpdater.GetAccessToken())
	}
}
