package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
)

// PKCEChallenge holds the PKCE code verifier and challenge
type PKCEChallenge struct {
	Verifier  string
	Challenge string
}

// Should generate a random code verifier and compute the SHA256 challenge
func GeneratePKCEChallenge() (*PKCEChallenge, error) {
	// generate random bytes
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	// encode as an unpadded string
	verifier := base64.RawURLEncoding.EncodeToString(randomBytes)
	// SHA256 hash the verifier
	hasher := sha256.New()
	hasher.Write([]byte(verifier))
	// create challenge
	challenge := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))

	return &PKCEChallenge{
		Verifier:  verifier,
		Challenge: challenge,
	}, nil

}

// Should build the complete OAuth authorization URL with PKCE parameters
func BuildAuthorizationURL(authURL, clientID, redirectURI, codeChallenge string) string {

	// "https://eggcarton-auth-uqhqvdut.auth.us-west-1.amazoncognito.com/oauth2/
	// authorize?client_id=1vccvf2hh5amna78lurbn9bjhi&response_type=token&scope=email+openid+profile&redirect_uri=https://oauth.pstmn.io/v1/callback"
	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("response_type", "code")
	params.Set("scope", "openid email profile")
	params.Set("redirect_uri", redirectURI)
	params.Set("code_challenge", codeChallenge)
	params.Set("code_challenge_method", "S256")

	// Hint: Use url.Values or fmt.Sprintf

	// return ""
	return fmt.Sprintf("%s?%s", authURL, params.Encode())
}
