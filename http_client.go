package selling_partner_api

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

type Region string
type Endpoint string

const (
	AWSRegionUSEast Region = "us-east-1"
	AWSRegionEUWest Region = "eu-west-1"
	AWSRegionUSWest Region = "us-west-2"

	EndpointNorthAmerica Endpoint = "https://sellingpartnerapi-na.amazon.com"
	EndpointEurope       Endpoint = "https://sellingpartnerapi-eu.amazon.com"
	EndpointFarEast      Endpoint = "https://sellingpartnerapi-fe.amazon.com"
)

type HttpClientConfig struct {
	client             *http.Client
	TokenUpdater       TokenUpdaterInterface
	IAMUserAccessKeyID string
	IAMUserSecretKey   string
	Region             Region
	RoleArn            string
	Endpoint           Endpoint
}

func NewHttpClient(config HttpClientConfig) (client *HttpClient, err error) {
	c := HttpClient{
		client:       config.client,
		tokenUpdater: config.TokenUpdater,
		region:       config.Region,
		roleArn:      config.RoleArn,
		endpoint:     config.Endpoint,
	}
	c.awsSession, err = session.NewSession(
		&aws.Config{Credentials: credentials.NewStaticCredentials(config.IAMUserAccessKeyID, config.IAMUserSecretKey, "")},
	)
	return &c, nil
}

type HttpClient struct {
	client            *http.Client
	endpoint          Endpoint
	tokenUpdater      TokenUpdaterInterface
	region            Region
	roleArn           string
	aws4Signer        *v4.Signer
	awsStsCredentials *sts.Credentials
	awsSession        *session.Session
}

func (h *HttpClient) Do(req *http.Request) (*http.Response, error) {
	h.addAccessTokenToHeader(req)

	if err := h.signRequest(req); err != nil {
		return nil, err
	}

	return h.client.Do(req)
}

func (h *HttpClient) GetEndpoint() string {
	return string(h.endpoint)
}

func (h *HttpClient) addAccessTokenToHeader(req *http.Request) {
	if req.Header.Get("X-Amz-Access-Token") == "" {
		req.Header.Add("X-Amz-Access-Token", h.tokenUpdater.GetAccessToken())
	}
}

func (h *HttpClient) signRequest(r *http.Request) error {

	if h.aws4Signer == nil ||
		h.awsStsCredentials == nil ||
		h.aws4Signer.Credentials.IsExpired() ||
		h.awsStsCredentials.Expiration.IsZero() ||
		h.awsStsCredentials.Expiration.Round(0).Add(-ExpiryDelta).Before(time.Now().UTC()) {
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

	_, err := h.aws4Signer.Sign(r, body, "execute-api", string(h.region), time.Now().UTC())

	return err
}
func (h *HttpClient) RefreshCredentials() error {

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
