package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fond-of-vertigo/amazon-sp-api/constants"
	"github.com/fond-of-vertigo/logger"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

type TokenUpdaterConfig struct {
	RefreshToken string
	ClientID     string
	ClientSecret string
	Logger       logger.Logger
	QuitSignal   chan bool
}

type TokenUpdater interface {
	GetAccessToken() string
	RunInBackground() error
}

type tokenUpdater struct {
	AccessToken     *atomic.Value
	ExpireTimestamp *atomic.Int64
	RefreshToken    string
	ClientID        string
	ClientSecret    string
	Log             logger.Logger
	QuitSignal      chan bool
}
type AccessTokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	TokenType        string `json:"token_type"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func NewTokenUpdater(config TokenUpdaterConfig) TokenUpdater {
	t := tokenUpdater{
		RefreshToken: config.RefreshToken,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Log:          config.Logger,
		QuitSignal:   config.QuitSignal,
	}
	return &t
}

func (t *tokenUpdater) RunInBackground() error {
	t.ExpireTimestamp = &atomic.Int64{}
	t.AccessToken = &atomic.Value{}
	t.Log.Debugf("Fetching first access-tokenAPI")
	if err := t.fetchNewToken(); err != nil {
		return err
	}

	go t.checkAccessToken()
	return nil
}

func (t *tokenUpdater) checkAccessToken() {
	for {
		select {
		case <-t.QuitSignal:
			t.Log.Infof("Received signal to stop access-tokenAPI updates.")
			return
		default:
			secondsToWait := secondsUntilExpired(t.ExpireTimestamp.Load())
			if secondsToWait <= int64(constants.ExpiryDelta.Seconds()) {
				if err := t.fetchNewToken(); err != nil {
					t.Log.Errorf(err.Error())
				}
			} else {
				time.Sleep(time.Duration(secondsToWait-int64(constants.ExpiryDelta.Seconds())) * time.Second)
			}
		}
	}
}

func (t *tokenUpdater) GetAccessToken() string {
	return fmt.Sprintf("%v", t.AccessToken.Load())
}

func (t *tokenUpdater) fetchNewToken() error {
	reqBody, _ := json.Marshal(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": t.RefreshToken,
		"client_id":     t.ClientID,
		"client_secret": t.ClientSecret,
	})

	resp, err := http.Post(
		"https://api.amazon.com/auth/o2/token",
		"application/json",
		bytes.NewBuffer(reqBody))

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
		t.AccessToken.Swap(parsedResp.AccessToken)

		expireTimestamp := time.Now().UTC().Add(time.Duration(parsedResp.ExpiresIn) * time.Second)
		t.ExpireTimestamp.Swap(expireTimestamp.Unix())
	}
	return nil
}

func secondsUntilExpired(unixTimestamp int64) int64 {
	currentTimestamp := time.Now().Unix()
	return unixTimestamp - currentTimestamp
}
