package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"github.com/fond-of-vertigo/amazon-sp-api/constants"
	"github.com/fond-of-vertigo/logger"
)

const tokenURL = "https://api.amazon.com/auth/o2/token"

type TokenUpdaterConfig struct {
	RefreshToken string
	ClientID     string
	ClientSecret string
	HTTPClient   HTTPRequester
	Logger       logger.Logger
}

type PeriodicTokenUpdater struct {
	accessToken  atomic.Pointer[string]
	refreshToken string
	clientID     string
	clientSecret string
	hTTPClient   HTTPRequester
	log          logger.Logger
}

type AccessTokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	TokenType        string `json:"token_type"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func newTokenUpdater(config TokenUpdaterConfig) *PeriodicTokenUpdater {
	return &PeriodicTokenUpdater{
		refreshToken: config.RefreshToken,
		clientID:     config.ClientID,
		clientSecret: config.ClientSecret,
		log:          config.Logger,
		hTTPClient:   config.HTTPClient,
	}
}

// GetAccessToken returns the current access-token
func (t *PeriodicTokenUpdater) GetAccessToken() string {
	token := t.accessToken.Load()
	if token == nil {
		return ""
	}
	return *token
}

// RunInBackground starts a goroutine that fetches a new access token periodically
// and stores it in the client. The goroutine is stopped when the returned cancel function is called.
func (t *PeriodicTokenUpdater) RunInBackground() (cancel func(), err error) {
	durationNextFetch, err := t.doInitialFetch()
	if err != nil {
		return nil, err
	}

	ticker := time.NewTicker(durationNextFetch)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				t.log.Infof("Stopped goroutine of token-updater.")
				return
			case <-ticker.C:
				token, err := t.doTokenRequest()
				if err != nil {
					t.log.Errorf("Failed to fetch new access-tokenAPI: %s", err.Error())
					ticker.Reset(constants.DefaultTokenUpdaterBackoffTime)
					continue
				}
				t.accessToken.Store(&token.AccessToken)
				durationToWait := durationBetweenTokenRequests(token)
				ticker.Reset(durationToWait)
			}
		}
	}()

	cancelFunc := func() {
		ticker.Stop()
		done <- true
	}
	return cancelFunc, nil

}

func (t *PeriodicTokenUpdater) doInitialFetch() (time.Duration, error) {
	t.log.Debugf("Fetching first access-tokenAPI")
	token, err := t.doTokenRequest()
	if err != nil {
		return constants.DefaultTokenUpdaterBackoffTime, err
	}
	t.accessToken.Store(&token.AccessToken)
	durationNextFetch := durationBetweenTokenRequests(token)
	return durationNextFetch, nil
}

func durationBetweenTokenRequests(token *AccessTokenResponse) time.Duration {
	secondsToWait := token.ExpiresIn - int(constants.ExpiryDelta/time.Second)
	durationToWait := time.Duration(secondsToWait) * time.Second
	return durationToWait
}

func (t *PeriodicTokenUpdater) doTokenRequest() (*AccessTokenResponse, error) {
	body := makeRequestBody(t.refreshToken, t.clientID, t.clientSecret)
	resp, err := t.hTTPClient.Post(tokenURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			t.log.Errorf(err.Error())
		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tkn, err := t.parseAccessTokenResponse(respBody)
	if err != nil {
		return nil, err
	}

	if tkn.AccessToken == "" {
		return nil, errors.New("refreshToken response did not contain access token")
	}
	return tkn, nil
}

func (t *PeriodicTokenUpdater) parseAccessTokenResponse(body []byte) (*AccessTokenResponse, error) {
	parsedResp := &AccessTokenResponse{}
	if err := json.Unmarshal(body, parsedResp); err != nil {
		return nil, fmt.Errorf("refreshToken response parse failed. Body: %s", string(body))
	}
	return parsedResp, nil
}

func makeRequestBody(refreshToken, clientID, clientSecret string) []byte {
	body, _ := json.Marshal(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"client_id":     clientID,
		"client_secret": clientSecret,
	})
	return body
}
