package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// PKCEChallenge holds the PKCE code verifier and challenge
type PKCEChallenge struct {
	Verifier  string
	Challenge string
}

// TODO: Implement GeneratePKCEChallenge
// Should generate a random code verifier and compute the SHA256 challenge
func GeneratePKCEChallenge() (*PKCEChallenge, error) {
	// TODO:
	// 1. Generate 32 random bytes using crypto/rand
	// 2. Base64URL encode them (no padding) to create verifier
	// 3. SHA256 hash the verifier
	// 4. Base64URL encode the hash (no padding) to create challenge
	// 5. Return PKCEChallenge struct

	// Hint: Use base64.RawURLEncoding (no padding)
	// Verifier should be 43-128 characters

	verifier := make([]byte, 32)
	if _, err := rand.Read(verifier); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// TODO: Implement the rest
	_ = sha256.New()
	_ = base64.RawURLEncoding

	return nil, fmt.Errorf("not implemented")
}

// TODO: Implement BuildAuthorizationURL
// Should build the complete OAuth authorization URL with PKCE parameters
func BuildAuthorizationURL(authURL, clientID, redirectURI, codeChallenge string) string {
	// TODO:
	// Build URL with these query parameters:
	// - client_id
	// - response_type=code
	// - scope=openid email profile
	// - redirect_uri
	// - code_challenge
	// - code_challenge_method=S256

	// Hint: Use url.Values or fmt.Sprintf

	return ""
}
