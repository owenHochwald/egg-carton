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

// Should start a local HTTP server to receive the OAuth callback
func StartCallbackServer(ctx context.Context) (string, error) {
	resultChan := make(chan CallbackResult, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// Extract code and error from query params
		code := r.URL.Query().Get("code")
		errParam := r.URL.Query().Get("error")

		// Send friendly HTML response
		w.Header().Set("Content-Type", "text/html")
		if code != "" {
			fmt.Fprintf(w, `
				<html>
				<head><title>Authentication Successful</title></head>
				<body style="font-family: Arial; text-align: center; padding: 50px;">
					<h1>✅ Authentication Successful!</h1>
					<p>You can close this window and return to the terminal.</p>
				</body>
				</html>
			`)
		} else {
			fmt.Fprintf(w, `
				<html>
				<head><title>Authentication Failed</title></head>
				<body style="font-family: Arial; text-align: center; padding: 50px;">
					<h1>❌ Authentication Failed</h1>
					<p>Error: %s</p>
					<p>Please try again.</p>
				</body>
				</html>
			`, errParam)
		}

		// Send result through channel (non-blocking because it's buffered)
		resultChan <- CallbackResult{Code: code, Error: errParam}
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start server in a goroutine so it doesn't block
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	var result CallbackResult
	select {
	case result = <-resultChan:
		// Got callback! Shutdown server gracefully
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Printf("Server shutdown error: %v\n", err)
		}

	case <-ctx.Done():
		// Timeout or cancellation
		server.Close() // Force close
		return "", fmt.Errorf("authentication timeout or cancelled")
	}

	// Check if we got an error from OAuth
	if result.Error != "" {
		return "", fmt.Errorf("authentication failed: %s", result.Error)
	}

	// Check if we got a code
	if result.Code == "" {
		return "", fmt.Errorf("no authorization code received")
	}

	return result.Code, nil
}
