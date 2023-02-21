package httpx

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/fond-of-vertigo/amazon-sp-api/constants"
	"github.com/google/uuid"
)

type ClientConfig struct {
	HttpClient         *http.Client
	TokenUpdaterConfig TokenUpdaterConfig
	IAMUserAccessKeyID string
	IAMUserSecretKey   string
	Region             constants.Region
	RoleArn            string
	Endpoint           constants.Endpoint
}

func NewClient(config ClientConfig) (*Client, error) {
	c := &Client{
		httpClient: config.HttpClient,
		region:     config.Region,
		roleArn:    config.RoleArn,
		endpoint:   config.Endpoint,
	}

	c.tokenUpdater = newTokenUpdater(config.TokenUpdaterConfig)
	if err := c.tokenUpdater.RunInBackground(); err != nil {
		return nil, err
	}

	var err error
	if c.awsSession, err = session.NewSession(
		&aws.Config{Credentials: credentials.NewStaticCredentials(config.IAMUserAccessKeyID, config.IAMUserSecretKey, "")},
	); err != nil {
		return nil, err
	}
	return c, nil
}

type Client struct {
	tokenUpdater      tokenUpdater
	httpClient        *http.Client
	endpoint          constants.Endpoint
	region            constants.Region
	roleArn           string
	aws4Signer        *v4.Signer
	awsStsCredentials *sts.Credentials
	awsSession        *session.Session
}

func (h *Client) Do(req *http.Request) (*http.Response, error) {
	h.addAccessTokenToHeader(req)

	if err := h.signRequest(req); err != nil {
		return nil, err
	}

	return h.httpClient.Do(req)
}

func (h *Client) GetEndpoint() constants.Endpoint {
	return h.endpoint
}

func (h *Client) Close() {
	h.tokenUpdater.Stop()
}

func (h *Client) addAccessTokenToHeader(req *http.Request) {
	if req.Header.Get(constants.AccessTokenHeader) == "" {
		req.Header.Add(constants.AccessTokenHeader, h.tokenUpdater.GetAccessToken())
	}
}

func (h *Client) signRequest(r *http.Request) error {

	if h.aws4Signer == nil ||
		h.awsStsCredentials == nil ||
		h.aws4Signer.Credentials.IsExpired() ||
		h.awsStsCredentials.Expiration.IsZero() ||
		h.awsStsCredentials.Expiration.Round(0).Add(-constants.ExpiryDelta).Before(time.Now().UTC()) {
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

	_, err := h.aws4Signer.Sign(r, body, constants.ServiceExecuteAPI, string(h.region), time.Now().UTC())

	return err
}
func (h *Client) RefreshCredentials() error {

	roleSessionName := uuid.New().String()

	role, err := sts.New(h.awsSession).AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(h.roleArn),
		RoleSessionName: aws.String(roleSessionName),
	})

	if err != nil {
		return fmt.Errorf("RefreshCredentials call failed with error %w", err)
	}

	if role == nil || role.Credentials == nil {
		return fmt.Errorf("AssumeRole call failed in return")
	}

	h.awsStsCredentials = role.Credentials

	h.aws4Signer = v4.NewSigner(credentials.NewStaticCredentials(
		*role.Credentials.AccessKeyId,
		*role.Credentials.SecretAccessKey,
		*role.Credentials.SessionToken),
		func(s *v4.Signer) {
			s.DisableURIPathEscaping = true
		},
	)

	return nil
}
