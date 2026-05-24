package config

import (
	"strings"
	"testing"
)

func TestLoadCredentials(t *testing.T) {
	t.Setenv(ClientIDEnv, "client-id")
	t.Setenv(ClientSecretEnv, "client-secret")
	t.Setenv(TokenURLEnv, "https://token.example")
	t.Setenv(APIAddressEnv, "api.example:443")
	t.Setenv(CommunityIDEnv, "community-id")

	creds, err := LoadCredentials()
	if err != nil {
		t.Fatalf("LoadCredentials() error = %v", err)
	}
	if creds.ClientID != "client-id" || creds.ClientSecret != "client-secret" || creds.TokenURL != "https://token.example" || creds.APIAddress != "api.example:443" || creds.CommunityID != "community-id" {
		t.Fatalf("unexpected credentials: %+v", creds)
	}
}

func TestLoadCredentialsMissing(t *testing.T) {
	_, err := LoadCredentials()
	if err == nil {
		t.Fatal("LoadCredentials() error = nil, want missing env error")
	}
	if !strings.Contains(err.Error(), ClientIDEnv) {
		t.Fatalf("error %q does not include missing env name", err)
	}
	if !strings.Contains(err.Error(), TokenURLEnv) {
		t.Fatalf("error %q does not include missing shared token URL env name", err)
	}
	if !strings.Contains(err.Error(), APIAddressEnv) {
		t.Fatalf("error %q does not include missing shared API address env name", err)
	}
	if !strings.Contains(err.Error(), CommunityIDEnv) {
		t.Fatalf("error %q does not include missing community ID env name", err)
	}
}
