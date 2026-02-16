package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/owenHochwald/egg-carton/cli/config"
)

// TokenResponse represents the OAuth token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// TODO: Implement ExchangeCodeForTokens
// Should exchange the authorization code for JWT tokens
func ExchangeCodeForTokens(tokenURL, clientID, code, redirectURI, codeVerifier string) (*config.TokenData, error) {
	// TODO:
	// 1. Build form data with:
	//    - grant_type=authorization_code
	//    - client_id
	//    - code (authorization code from callback)
	//    - redirect_uri
	//    - code_verifier (from PKCE)
	// 2. Make POST request to token endpoint
	// 3. Parse JSON response into TokenResponse
	// 4. Convert to config.TokenData with IssuedAt timestamp
	// 5. Return error if token exchange fails

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	// TODO: Set other fields

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// TODO: Execute request
	// TODO: Parse response

	_ = io.ReadAll
	_ = json.Unmarshal
	_ = time.Now().Unix()

	return nil, fmt.Errorf("not implemented")
}

// TODO: Implement RefreshAccessToken
// Should refresh an expired access token using the refresh token
func RefreshAccessToken(tokenURL, clientID, refreshToken string) (*config.TokenData, error) {
	// TODO:
	// 1. Build form data with:
	//    - grant_type=refresh_token
	//    - client_id
	//    - refresh_token
	// 2. Make POST request to token endpoint
	// 3. Parse JSON response
	// 4. Return new TokenData

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	// TODO: Set other fields

	return nil, fmt.Errorf("not implemented")
}
