package selling_partner_api

import (
	"github.com/fond-of-vertigo/logger"
	"net/url"
)

type Config struct {
	ClientID           string
	ClientSecret       string
	RefreshToken       string
	IAMUserAccessKeyID string
	IAMUserSecretKey   string
	Region             string
	RoleArn            string
	Endpoint           url.URL
	Log                logger.Logger
}

type SellingPartnerClient struct {
	tokenRefresher *TokenRefresher
}

func NewSellingPartnerClient(config Config) (*SellingPartnerClient, error) {
	tokenRefresher, err := NewTokenRefresher(TokenRefresherConfig{
		RefreshToken: config.RefreshToken,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Logger:       config.Log,
	})
	if err != nil {
		return nil, err
	}

	return &SellingPartnerClient{tokenRefresher: tokenRefresher}, nil
}
