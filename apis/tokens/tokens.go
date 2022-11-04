package tokens

import (
	"encoding/json"
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
)

type Token struct {
	HttpClient apis.HttpRequestDoer
}

// CreateRestrictedDataTokenRequest returns a Restricted Data Token (RDT) for one or more restricted resources that you specify.
func (t *Token) CreateRestrictedDataTokenRequest(restrictedResources []RestrictedResource) (resp *CreateRestrictedDataTokenResponse, err error) {
	params := apis.APICall{}
	params.Method = "POST"
	params.APIPath = config.pathPrefix() + "/restrictedDataToken"
	params.Body, err = json.Marshal(restrictedResources)
	if err != nil {
		return nil, err
	}
	return apis.CallAPIWithResponseType[CreateRestrictedDataTokenResponse](params, t.HttpClient)
}
