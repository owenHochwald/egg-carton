package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
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

// Should save tokens to ~/.eggcarton/credentials.json with 0600 permissions
func (c *Config) SaveTokens(tokens *TokenData) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(c.TokenPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal tokens to JSON
	b, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tokens: %w", err)
	}

	// Write to file with secure permissions
	if err := os.WriteFile(c.TokenPath, b, 0600); err != nil {
		return fmt.Errorf("failed to write tokens to file: %w", err)
	}

	return nil
}

// Should load tokens from ~/.eggcarton/credentials.json
func (c *Config) LoadTokens() (*TokenData, error) {
	data, err := os.ReadFile(c.TokenPath)
	if err != nil {
		return nil, err
	}

	var tokens TokenData
	err = json.Unmarshal(data, &tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tokens: %w", err)
	}

	return &tokens, nil
}

// Should check if access token is still valid (not expired)
func (t *TokenData) IsTokenValid() bool {
	now := time.Now().Unix()
	if now > t.IssuedAt+int64(t.ExpiresIn)-300 { // 5 minute buffer
		return false
	}
	return true
}

// Returns the OAuth redirect URI for the callback server
func (c *Config) GetRedirectURI() string {
	return "http://localhost:8080/callback"
}

// Returns the full authorization URL for Cognito
func (c *Config) GetAuthorizationURL() string {

	return fmt.Sprintf("https://%s/oauth2/authorize", c.CognitoConfig.Domain)
}

// Returns the token exchange endpoint
func (c *Config) GetTokenURL() string {
	return fmt.Sprintf("https://%s/oauth2/token", c.CognitoConfig.Domain)
}
