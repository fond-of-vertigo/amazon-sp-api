package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/fond-of-vertigo/amazon-sp-api/constants"
	"github.com/fond-of-vertigo/logger"
)

const tokenURL string = "https://api.amazon.com/auth/o2/token"
const retryOnErrorDuration time.Duration = 10 * time.Second

type TokenUpdaterConfig struct {
	RefreshToken string
	ClientID     string
	ClientSecret string
	HTTPClient   *http.Client
	Logger       logger.Logger
}

type tokenUpdater interface {
	GetAccessToken() string
	RunInBackground() error
	Stop()
}

type tokenUpdaterData struct {
	timerPtr     atomic.Pointer[time.Timer]
	accessToken  atomic.Pointer[string]
	client       *http.Client
	RefreshToken string
	ClientID     string
	ClientSecret string
	Log          logger.Logger
}

type accessTokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	TokenType        string `json:"token_type"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type parsedToken struct {
	token        string
	waitDuration time.Duration
}

func makeTokenUpdater(config TokenUpdaterConfig, client *http.Client) tokenUpdater {
	return &tokenUpdaterData{
		client:       client,
		RefreshToken: config.RefreshToken,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Log:          config.Logger,
	}
}

// RunInBackground fetches a new token and then starts the token updater.
// It directly fetches a new token and starts the token updater only if
// the first token call was successful.
func (t *tokenUpdaterData) RunInBackground() error {
	t.Log.Debugf("Fetching first access-tokenAPI")
	resp, err := t.fetchNewToken()
	if err != nil {
		return err
	}

	t.accessToken.Store(&resp.token)
	t.timerPtr.Store(time.AfterFunc(resp.waitDuration, t.endlesslyRefetchToken))
	return nil
}

func (t *tokenUpdaterData) Stop() {
	timer := t.timerPtr.Load()
	if timer != nil {
		timer.Stop()
		t.timerPtr.Store(nil)
	}
}

func (t *tokenUpdaterData) GetAccessToken() string {
	token := t.accessToken.Load()
	if token == nil {
		return ""
	}
	return *token
}

func (t *tokenUpdaterData) endlesslyRefetchToken() {
	resp, err := t.fetchNewToken()
	if err != nil {
		t.Log.Errorf("fetchNewToken failed: %s", err)
		t.Log.Infof("retrying to fetch token after %s", retryOnErrorDuration)
		t.timerPtr.Store(time.AfterFunc(retryOnErrorDuration, t.endlesslyRefetchToken))
		return
	}

	t.accessToken.Store(&resp.token)
	t.timerPtr.Store(time.AfterFunc(resp.waitDuration, t.endlesslyRefetchToken))
	t.Log.Debugf("sucessfully updated access token")
}

func (t *tokenUpdaterData) fetchNewToken() (*parsedToken, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": t.RefreshToken,
		"client_id":     t.ClientID,
		"client_secret": t.ClientSecret,
	})

	resp, err := t.client.Post(tokenURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			t.Log.Errorf(err.Error())
		}
	}(resp.Body)

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tokenResp := &accessTokenResponse{}
	if err = json.Unmarshal(respBodyBytes, tokenResp); err != nil {
		return nil, fmt.Errorf("RefreshToken response parse failed. Body: %s", string(respBodyBytes))
	}

	if tokenResp.AccessToken != "" {
		secondsToWait := tokenResp.ExpiresIn - int(constants.ExpiryDelta/time.Second)
		durationToWait := time.Duration(secondsToWait) * time.Second
		return &parsedToken{
			token:        tokenResp.AccessToken,
			waitDuration: durationToWait,
		}, nil
	}

	return nil, fmt.Errorf("received an empty token")
}
