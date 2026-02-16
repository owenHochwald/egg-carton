package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client represents the API client for Lambda functions
type Client struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewClient creates a new API client
func NewClient(baseURL, accessToken string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   accessToken,
		client:  &http.Client{},
	}
}

// PutEggRequest represents the request body for storing a secret
type PutEggRequest struct {
	Owner     string `json:"owner"`
	SecretID  string `json:"secret_id"`
	Plaintext string `json:"plaintext"`
}

// GetEggResponse represents the response from getting a secret
type GetEggResponse struct {
	Owner     string `json:"owner"`
	SecretID  string `json:"secret_id"`
	Plaintext string `json:"plaintext"`
	CreatedAt string `json:"created_at"`
}

// TODO: Implement PutEgg
// Should call PUT /eggs endpoint to store a secret
func (c *Client) PutEgg(owner, key, value string) error {
	// TODO:
	// 1. Build request body with owner, secret_id, plaintext
	// 2. Marshal to JSON
	// 3. Create PUT request to {baseURL}/eggs
	// 4. Add Authorization: Bearer {token} header
	// 5. Execute request
	// 6. Check status code (200 = success)
	// 7. Return error if failed

	reqBody := PutEggRequest{
		Owner:     owner,
		SecretID:  key,
		Plaintext: value,
	}

	_ = reqBody
	_ = json.Marshal
	_ = bytes.NewBuffer

	return fmt.Errorf("not implemented")
}

// TODO: Implement GetEgg
// Should call GET /eggs/{owner} endpoint to retrieve a secret
func (c *Client) GetEgg(owner string) (string, error) {
	// TODO:
	// 1. Create GET request to {baseURL}/eggs/{owner}
	// 2. Add Authorization: Bearer {token} header
	// 3. Execute request
	// 4. Check status code (200 = success, 404 = not found)
	// 5. Parse JSON response into GetEggResponse
	// 6. Return the plaintext value

	_ = io.ReadAll
	_ = json.Unmarshal

	return "", fmt.Errorf("not implemented")
}

// TODO: Implement BreakEgg
// Should call DELETE /eggs/{owner} endpoint to delete a secret
func (c *Client) BreakEgg(owner string) error {
	// TODO:
	// 1. Create DELETE request to {baseURL}/eggs/{owner}
	// 2. Add Authorization: Bearer {token} header
	// 3. Execute request
	// 4. Check status code (200 = success, 404 = not found)
	// 5. Return error if failed

	return fmt.Errorf("not implemented")
}

// TODO (Optional for EC-14): Implement ListEggs
// Should call GET /eggs endpoint to list all secrets for a user
// You may need to add a new Lambda function for this
func (c *Client) ListEggs(owner string) (map[string]string, error) {
	// TODO:
	// This might require a new Lambda function that scans DynamoDB for all eggs belonging to owner
	// For now, return not implemented

	return nil, fmt.Errorf("not implemented - may need new Lambda endpoint")
}

// TODO: Implement ExtractOwnerFromToken
// Should decode JWT and extract the 'sub' claim (user ID)
func ExtractOwnerFromToken(accessToken string) (string, error) {
	// TODO:
	// 1. JWT has 3 parts separated by '.'
	// 2. Second part is the payload (base64url encoded JSON)
	// 3. Decode it and extract 'sub' claim

	// Hint: Use strings.Split and base64.RawURLEncoding

	return "", fmt.Errorf("not implemented")
}

// Helper function to make authenticated requests
func (c *Client) doRequest(method, path string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}
