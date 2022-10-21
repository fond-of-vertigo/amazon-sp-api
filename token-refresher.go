package selling_partner_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fond-of-vertigo/logger"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

type TokenRefresherConfig struct {
	RefreshToken string
	ClientID     string
	ClientSecret string
	Logger       logger.Logger
	QuitSignal   chan bool
}

type TokenRefresher struct {
	accessToken     *atomic.Value
	ExpireTimestamp *atomic.Int64
	config          TokenRefresherConfig
	log             logger.Logger
	quitSignal      chan bool
}
type AccessTokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	TokenType        string `json:"token_type"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func NewTokenRefresher(config TokenRefresherConfig) (*TokenRefresher, error) {
	tokenRefresher := TokenRefresher{
		config:     config,
		log:        config.Logger,
		quitSignal: config.QuitSignal,
	}
	if err := tokenRefresher.fetchNewToken(); err != nil {
		return nil, fmt.Errorf("accesstoken could not be fetched: %w", err)
	}
	return &tokenRefresher, nil
}

func (t *TokenRefresher) StartUpdatesInBackground() {
	go t.checkAccessToken()
}

func (t *TokenRefresher) checkAccessToken() {
	for {
		select {
		case <-t.quitSignal:
			t.log.Infof("Received signal to stop token updates.")
			return
		default:
			secondsToWait := secondsUntilExpired(t.ExpireTimestamp.Load())
			if secondsToWait <= 5 {
				if err := t.fetchNewToken(); err != nil {
					t.log.Errorf(err.Error())
				}
			} else {
				time.Sleep(time.Duration(secondsToWait-5) * time.Second)
			}
		}
	}
}

func (t *TokenRefresher) GetAccessToken() string {
	return fmt.Sprintf("%v", t.accessToken.Load())
}

func (t *TokenRefresher) fetchNewToken() error {
	reqBody, _ := json.Marshal(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": t.config.RefreshToken,
		"client_id":     t.config.ClientID,
		"client_secret": t.config.ClientSecret,
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
			t.log.Errorf(err.Error())
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
		t.accessToken.Swap(parsedResp.AccessToken)

		expireTimestamp := time.Now().UTC().Add(time.Duration(parsedResp.ExpiresIn) * time.Second)
		t.ExpireTimestamp.Swap(expireTimestamp.Unix())
	}
	return nil
}

func secondsUntilExpired(unixTimestamp int64) int64 {
	currentTimestamp := time.Now().Unix()
	return unixTimestamp - currentTimestamp
}
