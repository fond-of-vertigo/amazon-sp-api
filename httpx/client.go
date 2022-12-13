package httpx

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/fond-of-vertigo/amazon-sp-api/constants"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

type Client interface {
	Do(*http.Request) (*http.Response, error)
	GetEndpoint() constants.Endpoint
}

type ClientConfig struct {
	HttpClient         *http.Client
	TokenUpdater       TokenUpdater
	IAMUserAccessKeyID string
	IAMUserSecretKey   string
	Region             constants.Region
	RoleArn            string
	Endpoint           constants.Endpoint
}

func NewClient(config ClientConfig) (Client, error) {
	c := &client{
		HttpClient:   config.HttpClient,
		TokenUpdater: config.TokenUpdater,
		Region:       config.Region,
		RoleArn:      config.RoleArn,
		Endpoint:     config.Endpoint,
	}
	var err error
	if c.AWSSession, err = session.NewSession(
		&aws.Config{Credentials: credentials.NewStaticCredentials(config.IAMUserAccessKeyID, config.IAMUserSecretKey, "")},
	); err != nil {
		return nil, err
	}
	return c, nil
}

type client struct {
	HttpClient        *http.Client
	Endpoint          constants.Endpoint
	TokenUpdater      TokenUpdater
	Region            constants.Region
	RoleArn           string
	AWS4Signer        *v4.Signer
	AWSStsCredentials *sts.Credentials
	AWSSession        *session.Session
}

func (h *client) Do(req *http.Request) (*http.Response, error) {
	h.addAccessTokenToHeader(req)

	if err := h.signRequest(req); err != nil {
		return nil, err
	}

	return h.HttpClient.Do(req)
}

func (h *client) GetEndpoint() constants.Endpoint {
	return h.Endpoint
}

func (h *client) addAccessTokenToHeader(req *http.Request) {
	if req.Header.Get(constants.AccessTokenHeader) == "" {
		req.Header.Add(constants.AccessTokenHeader, h.TokenUpdater.GetAccessToken())
	}
}

func (h *client) signRequest(r *http.Request) error {

	if h.AWS4Signer == nil ||
		h.AWSStsCredentials == nil ||
		h.AWS4Signer.Credentials.IsExpired() ||
		h.AWSStsCredentials.Expiration.IsZero() ||
		h.AWSStsCredentials.Expiration.Round(0).Add(-constants.ExpiryDelta).Before(time.Now().UTC()) {
		if err := h.RefreshCredentials(); err != nil {
			return fmt.Errorf("cannot refresh role credentials. Error: %w", err)
		}
	}

	var body io.ReadSeeker
	if r.Body != nil {
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		r.Body = io.NopCloser(bytes.NewReader(payload))
		body = bytes.NewReader(payload)
	}

	_, err := h.AWS4Signer.Sign(r, body, constants.ServiceExecuteAPI, string(h.Region), time.Now().UTC())

	return err
}
func (h *client) RefreshCredentials() error {

	roleSessionName := uuid.New().String()

	role, err := sts.New(h.AWSSession).AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(h.RoleArn),
		RoleSessionName: aws.String(roleSessionName),
	})

	if err != nil {
		return fmt.Errorf("RefreshCredentials call failed with error %w", err)
	}

	if role == nil || role.Credentials == nil {
		return fmt.Errorf("AssumeRole call failed in return")
	}

	h.AWSStsCredentials = role.Credentials

	h.AWS4Signer = v4.NewSigner(credentials.NewStaticCredentials(
		*role.Credentials.AccessKeyId,
		*role.Credentials.SecretAccessKey,
		*role.Credentials.SessionToken),
		func(s *v4.Signer) {
			s.DisableURIPathEscaping = true
		},
	)

	return nil
}
