package selling_partner_api

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

type TokenUpdaterInterface interface {
	GetAccessToken() string
}

type TokenUpdater struct {
	accessToken     *atomic.Value
	ExpireTimestamp *atomic.Int64
	refreshToken    string
	clientID        string
	clientSecret    string
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

func NewTokenUpdater(config TokenUpdaterConfig) *TokenUpdater {
	t := TokenUpdater{
		refreshToken: config.RefreshToken,
		clientID:     config.ClientID,
		clientSecret: config.ClientSecret,
		log:          config.Logger,
		quitSignal:   config.QuitSignal,
	}
	return &t
}

func (t *TokenUpdater) RunInBackground() error {
	t.ExpireTimestamp = &atomic.Int64{}
	t.accessToken = &atomic.Value{}
	t.log.Debugf("Fetching first access-token")
	if err := t.fetchNewToken(); err != nil {
		return err
	}

	go t.checkAccessToken()
	return nil
}

func (t *TokenUpdater) checkAccessToken() {
	for {
		select {
		case <-t.quitSignal:
			t.log.Infof("Received signal to stop access-token updates.")
			return
		default:
			secondsToWait := secondsUntilExpired(t.ExpireTimestamp.Load())
			if secondsToWait <= int64(constants.ExpiryDelta.Seconds()) {
				if err := t.fetchNewToken(); err != nil {
					t.log.Errorf(err.Error())
				}
			} else {
				time.Sleep(time.Duration(secondsToWait-int64(constants.ExpiryDelta.Seconds())) * time.Second)
			}
		}
	}
}

func (t *TokenUpdater) GetAccessToken() string {
	return fmt.Sprintf("%v", t.accessToken.Load())
}

func (t *TokenUpdater) fetchNewToken() error {
	reqBody, _ := json.Marshal(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": t.refreshToken,
		"client_id":     t.clientID,
		"client_secret": t.clientSecret,
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
