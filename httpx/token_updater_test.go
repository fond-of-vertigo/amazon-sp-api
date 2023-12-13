package httpx

import (
	"encoding/json"
	"errors"
	"github.com/fond-of-vertigo/logger"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockHTTPClient struct {
	testing.TB
	URL              string
	BodyType         string
	Body             []byte
	PostCallCount    int
	MockResponseBody []byte
}

func (m *mockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	return nil, nil
}

func (m *mockHTTPClient) Post(url string, bodyType string, body io.Reader) (*http.Response, error) {
	m.PostCallCount++
	assert.Equal(m, m.URL, url)
	assert.Equal(m, m.BodyType, bodyType)

	assert.NotNil(m, body)
	acutalBody, err := io.ReadAll(body)
	assert.NoError(m, err)
	assert.Equal(m, m.Body, acutalBody)

	resp := httptest.NewRecorder()
	_, err = resp.Write(m.MockResponseBody)
	assert.NoError(m, err)
	return resp.Result(), nil
}

func TestPeriodicTokenUpdater_RunInBackground(t *testing.T) {
	type args struct {
		RefreshToken      string
		ClientID          string
		ClientSecret      string
		ExpectedPostCount int
		MockTokenResponse AccessTokenResponse
	}
	tests := []struct {
		name      string
		WantError error

		args args
	}{
		{
			name:      "Simple",
			WantError: nil,
			args: args{
				RefreshToken:      "refreshToken",
				ClientID:          "clientID",
				ClientSecret:      "clientSecret",
				ExpectedPostCount: 2,
				MockTokenResponse: AccessTokenResponse{AccessToken: "accessToken", ExpiresIn: 1},
			},
		},
		{
			name:      "Fail",
			WantError: errors.New("refreshToken response did not contain access token"),
			args: args{
				RefreshToken:      "refreshToken",
				ClientID:          "clientID",
				ClientSecret:      "clientSecret",
				ExpectedPostCount: 1,
				MockTokenResponse: AccessTokenResponse{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			respBody, err := json.Marshal(tt.args.MockTokenResponse)
			assert.NoError(t, err)

			tu := newTokenUpdater(TokenUpdaterConfig{
				RefreshToken: tt.args.RefreshToken,
				ClientID:     tt.args.ClientID,
				ClientSecret: tt.args.ClientSecret,
				HTTPClient: &mockHTTPClient{
					TB:               t,
					URL:              tokenURL,
					BodyType:         "application/json",
					Body:             makeRequestBody(tt.args.RefreshToken, tt.args.ClientID, tt.args.ClientSecret),
					MockResponseBody: respBody,
				},
				Logger: logger.New(logger.LvlTrace),
			})

			//  when
			cancel, err := tu.RunInBackground()

			// then
			if tt.WantError != nil {
				assert.Error(t, err, tt.WantError)
			}

			assert.Equal(t, tt.args.MockTokenResponse.AccessToken, tu.GetAccessToken())
			time.Sleep(time.Duration(tt.args.MockTokenResponse.ExpiresIn) * time.Second)
			assert.Equal(t, tt.args.MockTokenResponse.AccessToken, tu.GetAccessToken())

			// wait for the next update
			time.Sleep((time.Duration(tt.args.MockTokenResponse.ExpiresIn) * time.Second) - time.Duration(500)*time.Millisecond)
			cancel()

			assert.Equal(t, tt.args.MockTokenResponse.AccessToken, tu.GetAccessToken())
			time.Sleep(time.Duration(tt.args.MockTokenResponse.ExpiresIn) * time.Second)
			assert.Equal(t, tt.args.MockTokenResponse.AccessToken, tu.GetAccessToken())
			assert.Equal(t, tt.args.ExpectedPostCount, tu.httpClient.(*mockHTTPClient).PostCallCount)
		})
	}
}
