package config

import (
	"strings"
	"testing"
)

func TestEnvPrefix(t *testing.T) {
	if got, want := EnvPrefix("math.CT"), "MIXI2_MATH_CT"; got != want {
		t.Fatalf("EnvPrefix() = %q, want %q", got, want)
	}
}

func TestLoadCredentials(t *testing.T) {
	t.Setenv("MIXI2_MATH_CT_CLIENT_ID", "client-id")
	t.Setenv("MIXI2_MATH_CT_CLIENT_SECRET", "client-secret")
	t.Setenv("MIXI2_MATH_CT_COMMUNITY_ID", "community-id")
	t.Setenv(TokenURLEnv, "https://token.example")
	t.Setenv(APIAddressEnv, "api.example:443")

	creds, err := LoadCredentials("math.CT")
	if err != nil {
		t.Fatalf("LoadCredentials() error = %v", err)
	}
	if creds.ClientID != "client-id" || creds.ClientSecret != "client-secret" || creds.TokenURL != "https://token.example" || creds.APIAddress != "api.example:443" || creds.CommunityID != "community-id" {
		t.Fatalf("unexpected credentials: %+v", creds)
	}
}

func TestLoadCredentialsMissing(t *testing.T) {
	t.Setenv("MIXI2_MATH_AG_CLIENT_ID", "client-id")
	t.Setenv("MIXI2_MATH_AG_CLIENT_SECRET", "client-secret")
	t.Setenv(TokenURLEnv, "https://token.example")
	t.Setenv(APIAddressEnv, "api.example:443")

	_, err := LoadCredentials("math.AG")
	if err == nil {
		t.Fatal("LoadCredentials() error = nil, want missing env error")
	}
	if !strings.Contains(err.Error(), "MIXI2_MATH_AG_COMMUNITY_ID") {
		t.Fatalf("error %q does not include missing community ID env name", err)
	}
}

func TestLoadApplicationCredentialsDoesNotRequireCommunityID(t *testing.T) {
	t.Setenv("MIXI2_MATH_CT_CLIENT_ID", "client-id")
	t.Setenv("MIXI2_MATH_CT_CLIENT_SECRET", "client-secret")
	t.Setenv(TokenURLEnv, "https://token.example")
	t.Setenv(APIAddressEnv, "api.example:443")

	creds, err := LoadApplicationCredentials("math.CT")
	if err != nil {
		t.Fatalf("LoadApplicationCredentials() error = %v", err)
	}
	if creds.ClientID != "client-id" || creds.ClientSecret != "client-secret" || creds.TokenURL != "https://token.example" || creds.APIAddress != "api.example:443" || creds.CommunityID != "" {
		t.Fatalf("unexpected credentials: %+v", creds)
	}
}

func TestLoadApplicationCredentialsMissing(t *testing.T) {
	_, err := LoadApplicationCredentials("math.AG")
	if err == nil {
		t.Fatal("LoadApplicationCredentials() error = nil, want missing env error")
	}
	for _, name := range []string{"MIXI2_MATH_AG_CLIENT_ID", TokenURLEnv, APIAddressEnv} {
		if !strings.Contains(err.Error(), name) {
			t.Fatalf("error %q does not include missing env name %s", err, name)
		}
	}
}
