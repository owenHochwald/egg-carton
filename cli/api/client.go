package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

// GetEggsResponse represents the response containing multiple secrets
type GetEggsResponse struct {
	Eggs []GetEggResponse `json:"eggs"`
}

// Should call PUT /eggs endpoint to store a secret
func (c *Client) PutEgg(owner, key, value string) error {
	request := PutEggRequest{
		Owner:     owner,
		SecretID:  key,
		Plaintext: value,
	}

	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := c.doRequest("PUT", "/eggs", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	defer req.Body.Close()

	if req.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(req.Body)
		return fmt.Errorf("failed to put egg (status %d) %s", req.StatusCode, body)

	}
	return nil
}

// GetEgg retrieves all secrets for an owner
func (c *Client) GetEgg(owner string) ([]GetEggResponse, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/eggs/%s", owner), nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get egg (status %d) %s", resp.StatusCode, body)
	}

	var response GetEggsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Eggs, nil
}

// Should call DELETE /eggs/{owner} endpoint to delete a secret
func (c *Client) BreakEgg(owner string) error {
	resp, err := c.doRequest("DELETE", fmt.Sprintf("/eggs/%s", owner), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to break egg (status %d) %s", resp.StatusCode, body)
	}

	return nil
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

// Should decode JWT and extract the 'sub' claim (user ID)
func ExtractOwnerFromToken(accessToken string) (string, error) {
	parts := strings.Split(accessToken, ".")

	if len(parts) != 3 {
		return "", fmt.Errorf("invalid JWT token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])

	if err != nil {
		return "", fmt.Errorf("failed to decode JWT payload: %w", err)
	}

	var claims struct {
		Sub string `json:"sub"`
	}

	if err = json.Unmarshal(payload, &claims); err != nil {
		return "", fmt.Errorf("failed to parse JWT claims: %w", err)
	}

	if claims.Sub == "" {
		return "", fmt.Errorf("sub claim not found in token")
	}

	return claims.Sub, nil
}

// function to make authenticated requests
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
