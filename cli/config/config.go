package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the CLI configuration
type Config struct {
	APIEndpoint   string        `json:"api_endpoint"`
	CognitoConfig CognitoConfig `json:"cognito"`
	TokenPath     string        `json:"-"` // Not serialized
}

// CognitoConfig holds Cognito-specific configuration
type CognitoConfig struct {
	UserPoolID string `json:"user_pool_id"`
	ClientID   string `json:"client_id"`
	Domain     string `json:"domain"`
	Region     string `json:"region"`
}

// TokenData holds the OAuth tokens
type TokenData struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	IssuedAt     int64  `json:"issued_at"` // Unix timestamp when token was received
}

// TODO: Implement LoadConfig
// Should load from environment variables or a default config
// For now, you can hardcode the values from terraform output
func LoadConfig() (*Config, error) {
	// TODO: Get these from environment variables or config file
	// For now, hardcode from your terraform output:
	config := &Config{
		APIEndpoint: "https://z7ha1j4xr9.execute-api.us-west-1.amazonaws.com/dev",
		CognitoConfig: CognitoConfig{
			UserPoolID: "us-west-1_2fzGdmIah",
			ClientID:   "1vccvf2hh5amna78lurbn9bjhi",
			Domain:     "eggcarton-auth-uqhqvdut.auth.us-west-1.amazoncognito.com",
			Region:     "us-west-1",
		},
	}

	// Set token path
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	config.TokenPath = filepath.Join(home, ".eggcarton", "credentials.json")

	return config, nil
}

// TODO: Implement SaveTokens
// Should save tokens to ~/.eggcarton/credentials.json with 0600 permissions
func (c *Config) SaveTokens(tokens *TokenData) error {
	// TODO:
	// 1. Create ~/.eggcarton directory if it doesn't exist
	// 2. Marshal tokens to JSON
	// 3. Write to file with 0600 permissions (os.FileMode(0600))
	// 4. Return error if anything fails

	return fmt.Errorf("not implemented")
}

// TODO: Implement LoadTokens
// Should load tokens from ~/.eggcarton/credentials.json
func (c *Config) LoadTokens() (*TokenData, error) {
	// TODO:
	// 1. Read file from c.TokenPath
	// 2. Unmarshal JSON into TokenData
	// 3. Return error if file doesn't exist or parsing fails

	return nil, fmt.Errorf("not implemented")
}

// TODO: Implement IsTokenValid
// Should check if access token is still valid (not expired)
func (t *TokenData) IsTokenValid() bool {
	// TODO:
	// 1. Get current time
	// 2. Compare with IssuedAt + ExpiresIn
	// 3. Add a small buffer (e.g., 5 minutes) to refresh before actual expiry

	return false
}

// TODO: Implement GetRedirectURI
// Returns the OAuth redirect URI for the callback server
func (c *Config) GetRedirectURI() string {
	return "http://localhost:8080/callback"
}

// TODO: Implement GetAuthorizationURL
// Returns the full authorization URL for Cognito
func (c *Config) GetAuthorizationURL() string {
	return fmt.Sprintf("https://%s/oauth2/authorize", c.CognitoConfig.Domain)
}

// TODO: Implement GetTokenURL
// Returns the token exchange endpoint
func (c *Config) GetTokenURL() string {
	return fmt.Sprintf("https://%s/oauth2/token", c.CognitoConfig.Domain)
}
