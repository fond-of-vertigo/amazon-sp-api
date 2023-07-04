package tokens

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/httpx"
)

const pathPrefix = "/tokens/2021-03-01"

type API struct {
	httpClient *httpx.Client
}

func NewAPI(httpClient *httpx.Client) *API {
	return &API{
		httpClient: httpClient,
	}
}

// CreateRestrictedDataTokenRequest returns a Restricted Data Token (RDT) for one or more restricted resources that you specify.
func (t *API) CreateRestrictedDataTokenRequest(restrictedResources *CreateRestrictedDataTokenRequest) (*apis.CallResponse[CreateRestrictedDataTokenResponse], error) {
	body, err := json.Marshal(restrictedResources)
	if err != nil {
		return nil, err
	}
	return apis.NewCall[CreateRestrictedDataTokenResponse](http.MethodPost, pathPrefix+"/restrictedDataToken").
		WithBody(body).
		WithRateLimit(1.0, time.Second).
		WithParseErrorListOnError(true).
		Execute(t.httpClient)
}
