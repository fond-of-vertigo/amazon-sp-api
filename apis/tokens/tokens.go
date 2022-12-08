package tokens

import (
	"encoding/json"
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"net/http"
)

const pathPrefix = "/tokens/2021-03-01"

type API interface {
	// CreateRestrictedDataTokenRequest returns a Restricted Data Token (RDT) for one or more restricted resources that you specify.
	CreateRestrictedDataTokenRequest(*CreateRestrictedDataTokenRequest) (*CreateRestrictedDataTokenResponse, apis.CallError)
}

func NewAPI(httpClient apis.HttpRequestDoer) API {
	return &api{
		HttpClient: httpClient,
	}
}

type api struct {
	HttpClient apis.HttpRequestDoer
}

func (t *api) CreateRestrictedDataTokenRequest(restrictedResources *CreateRestrictedDataTokenRequest) (*CreateRestrictedDataTokenResponse, apis.CallError) {
	body, err := json.Marshal(restrictedResources)
	if err != nil {
		return nil, apis.NewError(err)
	}
	return apis.NewCall[CreateRestrictedDataTokenResponse](http.MethodPost, pathPrefix+"/restrictedDataToken").
		WithBody(body).
		Execute(t.HttpClient)
}
