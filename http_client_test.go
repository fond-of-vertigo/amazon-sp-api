package selling_partner_api

import (
	"bytes"
	"net/http"
	"testing"
)

type mockTokenUpdater struct {
	ReturnAccessToken string
}

func (m *mockTokenUpdater) GetAccessToken() string {
	return m.ReturnAccessToken
}

func Test_httpClient_addAccessToken(t *testing.T) {
	reqWithRDT, _ := http.NewRequest("GET", "example.com", bytes.NewBufferString("example"))
	reqWithoutRDT, _ := http.NewRequest("GET", "example.com", bytes.NewBufferString("example"))
	reqWithRDT.Header.Add("X-Amz-Access-Token", "EXISTING-RDT")

	type fields struct {
		HttpClient   *http.Client
		TokenUpdater TokenUpdaterInterface
	}
	tests := []struct {
		name            string
		fields          fields
		req             *http.Request
		wantAccessToken string
	}{
		{
			name: "AccessToken should not replace an existing (e.g. RestrictedDataToken)",
			fields: fields{
				HttpClient:   nil,
				TokenUpdater: &mockTokenUpdater{ReturnAccessToken: "ACCESS-TOKEN-XY"},
			},
			req:             reqWithRDT,
			wantAccessToken: "EXISTING-RDT",
		},
		{
			name: "AccessToken should be inserted if no RestrictedDataToken is set",
			fields: fields{
				HttpClient:   nil,
				TokenUpdater: &mockTokenUpdater{ReturnAccessToken: "ACCESS-TOKEN-XY"},
			},
			req:             reqWithoutRDT,
			wantAccessToken: "ACCESS-TOKEN-XY",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &httpClient{
				HttpClient:   tt.fields.HttpClient,
				TokenUpdater: tt.fields.TokenUpdater,
			}
			h.addAccessToken(tt.req)
			if tt.req.Header.Get("X-Amz-Access-Token") != tt.wantAccessToken {
				t.Fatalf("Token %s != %s", tt.req.Header.Get("X-Amz-Access-Token"), tt.wantAccessToken)
			}
		})
	}
}
