package tokens

import (
	"encoding/json"
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/httpx"
	"net/http"
)

const pathPrefix = "/tokens/2021-03-01"

type API interface {
	// CreateRestrictedDataTokenRequest returns a Restricted Data Token (RDT) for one or more restricted resources that you specify.
	CreateRestrictedDataTokenRequest(*CreateRestrictedDataTokenRequest) (*apis.CallResponse[CreateRestrictedDataTokenResponse], error)
}

func NewAPI(httpClient httpx.Client) API {
	return &api{
		HttpClient: httpClient,
	}
}

type api struct {
	HttpClient httpx.Client
}

func (t *api) CreateRestrictedDataTokenRequest(restrictedResources *CreateRestrictedDataTokenRequest) (*apis.CallResponse[CreateRestrictedDataTokenResponse], error) {
	body, err := json.Marshal(restrictedResources)
	if err != nil {
		return nil, err
	}
	return apis.NewCall[CreateRestrictedDataTokenResponse](http.MethodPost, pathPrefix+"/restrictedDataToken").
		WithBody(body).
		Execute(t.HttpClient)
}
