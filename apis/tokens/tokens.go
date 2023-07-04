package tokens

import (
	"encoding/json"
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/httpx"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

const pathPrefix = "/tokens/2021-03-01"

type API interface {
	// CreateRestrictedDataTokenRequest returns a Restricted Data Token (RDT) for one or more restricted resources that you specify.
	CreateRestrictedDataTokenRequest(*CreateRestrictedDataTokenRequest) (*apis.CallResponse[CreateRestrictedDataTokenResponse], error)
}

func NewAPI(httpClient httpx.Client) API {
	return &api{
		HttpClient:                         httpClient,
		RateLimitCreateRestrictedDataToken: rate.NewLimiter(rate.Every(time.Second), 10),
	}
}

type api struct {
	HttpClient                         httpx.Client
	RateLimitCreateRestrictedDataToken *rate.Limiter
}

func (t *api) CreateRestrictedDataTokenRequest(restrictedResources *CreateRestrictedDataTokenRequest) (*apis.CallResponse[CreateRestrictedDataTokenResponse], error) {
	body, err := json.Marshal(restrictedResources)
	if err != nil {
		return nil, err
	}
	return apis.NewCall[CreateRestrictedDataTokenResponse](http.MethodPost, pathPrefix+"/restrictedDataToken").
		WithBody(body).
		WithRateLimiter(t.RateLimitCreateRestrictedDataToken).
		WithParseErrorListOnError(true).
		Execute(t.HttpClient)
}
