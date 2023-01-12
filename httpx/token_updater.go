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

const tokenURL = "https://api.amazon.com/auth/o2/token"

type TokenUpdaterConfig struct {
	RefreshToken string
	ClientID     string
	ClientSecret string
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
	RefreshToken string
	ClientID     string
	ClientSecret string
	Log          logger.Logger
}

type AccessTokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	TokenType        string `json:"token_type"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func makeTokenUpdater(config TokenUpdaterConfig) tokenUpdater {
	return &tokenUpdaterData{
		RefreshToken: config.RefreshToken,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Log:          config.Logger,
	}
}

func (t *tokenUpdaterData) RunInBackground() error {
	t.Log.Debugf("Fetching first access-tokenAPI")
	return t.fetchNewToken()
}

func (t *tokenUpdaterData) Stop() {
	timer := t.timerPtr.Load()
	if timer != nil {
		timer.Stop()
	}
}

func (t *tokenUpdaterData) GetAccessToken() string {
	token := t.accessToken.Load()
	if token == nil {
		return ""
	}
	return *token
}

func (t *tokenUpdaterData) fetchNewToken() error {
	reqBody, _ := json.Marshal(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": t.RefreshToken,
		"client_id":     t.ClientID,
		"client_secret": t.ClientSecret,
	})

	resp, err := http.Post(tokenURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			t.Log.Errorf(err.Error())
		}
	}(resp.Body)

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	parsedResp := &AccessTokenResponse{}
	if err = json.Unmarshal(respBodyBytes, parsedResp); err != nil {
		return fmt.Errorf("RefreshToken response parse failed. Body: %s", string(respBodyBytes))
	}

	if parsedResp.AccessToken != "" {
		t.accessToken.Store(&parsedResp.AccessToken)
		t.Log.Debugf("Successfully refreshed access token")

		secondsToWait := parsedResp.ExpiresIn - int(constants.ExpiryDelta/time.Second)
		durationToWait := time.Duration(secondsToWait) * time.Second
		t.timerPtr.Store(time.AfterFunc(durationToWait, func() {
			if err := t.fetchNewToken(); err != nil {
				t.Log.Errorf(err.Error())
			}
		}))
	}

	return nil
}
