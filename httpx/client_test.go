package httpx

import (
	"bytes"
	"github.com/fond-of-vertigo/amazon-sp-api/constants"
	"net/http"
	"testing"
)

type mockTokenUpdater struct {
	ReturnAccessToken string
}

func (m *mockTokenUpdater) GetAccessToken() string {
	return m.ReturnAccessToken
}
func (m *mockTokenUpdater) RunInBackground() error {
	return nil
}

func Test_httpClient_addAccessToken(t *testing.T) {
	reqWithRDT, _ := http.NewRequest(http.MethodGet, "example.com", bytes.NewBufferString("example"))
	reqWithoutRDT, _ := http.NewRequest(http.MethodGet, "example.com", bytes.NewBufferString("example"))
	reqWithRDT.Header.Add(constants.AccessTokenHeader, "EXISTING-RDT")

	type fields struct {
		HttpClient   *http.Client
		TokenUpdater TokenUpdater
	}
	tests := []struct {
		name            string
		fields          fields
		request         *http.Request
		wantAccessToken string
	}{
		{
			name: "AccessToken should not replace an existing (e.g. RestrictedDataToken)",
			fields: fields{
				HttpClient:   nil,
				TokenUpdater: &mockTokenUpdater{ReturnAccessToken: "ACCESS-TOKEN-XY"},
			},
			request:         reqWithRDT,
			wantAccessToken: "EXISTING-RDT",
		},
		{
			name: "AccessToken should be inserted if no RestrictedDataToken is set",
			fields: fields{
				HttpClient:   nil,
				TokenUpdater: &mockTokenUpdater{ReturnAccessToken: "ACCESS-TOKEN-XY"},
			},
			request:         reqWithoutRDT,
			wantAccessToken: "ACCESS-TOKEN-XY",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &client{
				HttpClient:   tt.fields.HttpClient,
				TokenUpdater: tt.fields.TokenUpdater,
			}
			h.addAccessTokenToHeader(tt.request)
			if tt.request.Header.Get(constants.AccessTokenHeader) != tt.wantAccessToken {
				t.Fatalf("Token %s != %s", tt.request.Header.Get(constants.AccessTokenHeader), tt.wantAccessToken)
			}
		})
	}
}
