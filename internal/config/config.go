package config

import (
	"fmt"
	"os"
	"strings"
)

type Credentials struct {
	ClientID     string
	ClientSecret string
	TokenURL     string
	APIAddress   string
	CommunityID  string
}

const (
	ClientIDEnv     = "MIXI2_CLIENT_ID"
	ClientSecretEnv = "MIXI2_CLIENT_SECRET"
	TokenURLEnv     = "MIXI2_TOKEN_URL"
	APIAddressEnv   = "MIXI2_API_ADDRESS"
	CommunityIDEnv  = "MIXI2_COMMUNITY_ID"
)

func LoadCredentials() (Credentials, error) {
	creds := Credentials{
		ClientID:     os.Getenv(ClientIDEnv),
		ClientSecret: os.Getenv(ClientSecretEnv),
		TokenURL:     os.Getenv(TokenURLEnv),
		APIAddress:   os.Getenv(APIAddressEnv),
		CommunityID:  os.Getenv(CommunityIDEnv),
	}

	var missing []string
	if creds.ClientID == "" {
		missing = append(missing, ClientIDEnv)
	}
	if creds.ClientSecret == "" {
		missing = append(missing, ClientSecretEnv)
	}
	if creds.TokenURL == "" {
		missing = append(missing, TokenURLEnv)
	}
	if creds.APIAddress == "" {
		missing = append(missing, APIAddressEnv)
	}
	if creds.CommunityID == "" {
		missing = append(missing, CommunityIDEnv)
	}
	if len(missing) > 0 {
		return Credentials{}, fmt.Errorf("missing environment variables: %s", strings.Join(missing, ", "))
	}
	return creds, nil
}
