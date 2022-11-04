package tokens

import (
	"encoding/json"
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"net/http"
)

const pathPrefix = "/tokens/2021-03-01"

type Token struct {
	HttpClient apis.HttpRequestDoer
}

// CreateRestrictedDataTokenRequest returns a Restricted Data Token (RDT) for one or more restricted resources that you specify.
func (t *Token) CreateRestrictedDataTokenRequest(restrictedResources []RestrictedResource) (resp *CreateRestrictedDataTokenResponse, err error) {
	params := apis.APICall{}
	params.Method = http.MethodPost
	params.APIPath = pathPrefix + "/restrictedDataToken"
	params.Body, err = json.Marshal(restrictedResources)
	if err != nil {
		return nil, err
	}
	return apis.CallAPIWithResponseType[CreateRestrictedDataTokenResponse](params, t.HttpClient)
}
