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
	TokenURLEnv   = "MIXI2_TOKEN_URL"
	APIAddressEnv = "MIXI2_API_ADDRESS"
)

func EnvPrefix(category string) string {
	replacer := strings.NewReplacer(".", "_", "-", "_")
	return "MIXI2_" + strings.ToUpper(replacer.Replace(category))
}

func LoadCredentials(category string) (Credentials, error) {
	prefix := EnvPrefix(category)
	creds := Credentials{
		ClientID:     os.Getenv(prefix + "_CLIENT_ID"),
		ClientSecret: os.Getenv(prefix + "_CLIENT_SECRET"),
		TokenURL:     os.Getenv(TokenURLEnv),
		APIAddress:   os.Getenv(APIAddressEnv),
		CommunityID:  os.Getenv(prefix + "_COMMUNITY_ID"),
	}

	var missing []string
	if creds.ClientID == "" {
		missing = append(missing, prefix+"_CLIENT_ID")
	}
	if creds.ClientSecret == "" {
		missing = append(missing, prefix+"_CLIENT_SECRET")
	}
	if creds.TokenURL == "" {
		missing = append(missing, TokenURLEnv)
	}
	if creds.APIAddress == "" {
		missing = append(missing, APIAddressEnv)
	}
	if creds.CommunityID == "" {
		missing = append(missing, prefix+"_COMMUNITY_ID")
	}
	if len(missing) > 0 {
		return Credentials{}, fmt.Errorf("missing environment variables for %s: %s", category, strings.Join(missing, ", "))
	}
	return creds, nil
}
