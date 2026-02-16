package main_test

// This file contains example tests you can write as you develop each component

import (
	"testing"
	// Import your packages as you implement them
)

// Phase 1 Tests - Config
func TestConfigLoadTokens(t *testing.T) {
	// TODO: Test loading tokens from file
	t.Skip("Implement this after config.LoadTokens()")
}

func TestConfigSaveTokens(t *testing.T) {
	// TODO: Test saving tokens with correct permissions
	t.Skip("Implement this after config.SaveTokens()")
}

func TestTokenIsValid(t *testing.T) {
	// TODO: Test token expiration logic
	t.Skip("Implement this after config.IsTokenValid()")
}

// Phase 2 Tests - Auth
func TestGeneratePKCEChallenge(t *testing.T) {
	// TODO: Test PKCE generation
	// Verify verifier length (43 chars)
	// Verify challenge is SHA256 hash of verifier
	t.Skip("Implement this after auth.GeneratePKCEChallenge()")
}

func TestBuildAuthorizationURL(t *testing.T) {
	// TODO: Test URL building
	// Verify all required parameters are present
	t.Skip("Implement this after auth.BuildAuthorizationURL()")
}

// Phase 3 Tests - API
func TestExtractOwnerFromToken(t *testing.T) {
	// TODO: Test JWT decoding
	// Use a test JWT token
	t.Skip("Implement this after api.ExtractOwnerFromToken()")
}

func TestAPIClientPutEgg(t *testing.T) {
	// TODO: Test API client (use httptest for mock server)
	t.Skip("Implement this after api.PutEgg()")
}

// Example of how to test the full flow
func TestFullLoginFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// TODO: This would be a full integration test
	// 1. Generate PKCE
	// 2. Build auth URL
	// 3. Mock the callback (don't actually open browser)
	// 4. Exchange code for tokens
	// 5. Save tokens
	// 6. Verify tokens are saved

	t.Skip("Implement this as an integration test")
}

// Run tests with:
// go test -v
// go test -v -short  (skip integration tests)
