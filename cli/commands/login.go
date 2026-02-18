package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/owenHochwald/egg-carton/cli/auth"
	"github.com/owenHochwald/egg-carton/cli/config"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

// LoginCmd represents the login command
var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Cognito via OAuth",
	Long: `Opens your browser to authenticate with AWS Cognito.
	
Uses PKCE flow for secure authentication without client secrets.
Tokens are stored locally in ~/.eggcarton/credentials.json`,
	RunE: runLogin,
}

func runLogin(cmd *cobra.Command, args []string) error {
	fmt.Println("🔐 Starting authentication flow...")

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	existingTokens, _ := cfg.LoadTokens()
	if existingTokens != nil && existingTokens.IsTokenValid() {
		fmt.Println("You are already logged in!")
		fmt.Println("Your session is still valid. Use --force to re-authenticate.")
		return nil
	}

	fmt.Println("Generating PKCE challenge...")
	pkce, err := auth.GeneratePKCEChallenge()
	if err != nil {
		return fmt.Errorf("failed to generate PKCE: %w", err)
	}

	authURL := auth.BuildAuthorizationURL(
		cfg.GetAuthorizationURL(),
		cfg.CognitoConfig.ClientID,
		cfg.GetRedirectURI(),
		pkce.Challenge,
	)

	fmt.Printf("If browser doesn't open, visit:\n   %s\n\n", authURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Start server in background
	serverErrChan := make(chan error, 1)
	var authCode string
	go func() {
		code, err := auth.StartCallbackServer(ctx)
		if err != nil {
			serverErrChan <- err
			return
		}
		authCode = code
		serverErrChan <- nil
	}()

	time.Sleep(500 * time.Millisecond) // Give server time to start
	if err := browser.OpenURL(authURL); err != nil {
		fmt.Printf("Failed to open browser automatically: %v\n", err)
		fmt.Printf("Please open this URL manually:\n%s\n", authURL)
	}

	if err := <-serverErrChan; err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	fmt.Println("Authorization code received!")

	fmt.Println("Exchanging code for tokens...")
	tokens, err := auth.ExchangeCodeForTokens(
		cfg.GetTokenURL(),
		cfg.CognitoConfig.ClientID,
		authCode,
		cfg.GetRedirectURI(),
		pkce.Verifier,
	)
	if err != nil {
		return fmt.Errorf("failed to exchange code for tokens: %w", err)
	}

	if err := cfg.SaveTokens(tokens); err != nil {
		return fmt.Errorf("failed to save tokens: %w", err)
	}

	fmt.Println("\n🎉 Login successful!")

	return nil
}
