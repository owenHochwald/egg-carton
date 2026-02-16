package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// CallbackResult holds the result of the OAuth callback
type CallbackResult struct {
	Code  string
	Error string
}

// TODO: Implement StartCallbackServer
// Should start a local HTTP server to receive the OAuth callback
func StartCallbackServer(ctx context.Context) (string, error) {
	// TODO:
	// 1. Create a channel to receive the authorization code
	// 2. Create HTTP handler for /callback route
	// 3. Extract 'code' or 'error' from query parameters
	// 4. Send code through channel
	// 5. Respond with success/error HTML page
	// 6. Start server on localhost:8080
	// 7. Use context to gracefully shutdown after receiving code
	// 8. Return the authorization code

	resultChan := make(chan CallbackResult, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Extract code from query params
		// r.URL.Query().Get("code")
		// r.URL.Query().Get("error")

		// TODO: Send HTML response
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<h1>Authentication successful!</h1><p>You can close this window.</p>")

		// TODO: Send result through channel
		resultChan <- CallbackResult{Code: "REPLACE_WITH_ACTUAL_CODE"}
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// TODO: Start server in goroutine
	// TODO: Wait for result or timeout
	// TODO: Shutdown server gracefully

	_ = server
	_ = time.After(5 * time.Minute) // 5 minute timeout

	return "", fmt.Errorf("not implemented")
}
