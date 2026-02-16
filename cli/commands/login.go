package commands

import (
	"fmt"

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

	// TODO: Implement login flow
	// 1. Load config
	// 2. Check if already logged in (optional: ask to re-login)
	// 3. Generate PKCE code verifier and challenge
	// 4. Build authorization URL
	// 5. Start local callback server
	// 6. Open browser to authorization URL
	// 7. Wait for callback with authorization code
	// 8. Exchange code for tokens
	// 9. Save tokens
	// 10. Print success message

	return fmt.Errorf("not implemented - see auth/login.go, auth/server.go, auth/token.go")
}
