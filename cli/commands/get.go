package commands

import (
	"fmt"

	"github.com/owenHochwald/egg-carton/cli/api"
	"github.com/owenHochwald/egg-carton/cli/auth"
	"github.com/owenHochwald/egg-carton/cli/config"
	"github.com/spf13/cobra"
)

// GetCmd represents the get command
var GetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Retrieve a secret",
	Long:  `Decrypt and retrieve a secret from your EggCarton vault.`,
	Args:  cobra.MaximumNArgs(1), // 0 or 1 args - if no key, list all
	RunE:  runGet,
}

func runGet(cmd *cobra.Command, args []string) error {
	// 1. Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// 2. Load tokens (check if logged in)
	tokens, err := cfg.LoadTokens()
	if err != nil {
		return fmt.Errorf("you are not logged in. Please run 'egg login' first: %w", err)
	}

	// 3. Check if token is valid (refresh if needed)
	if !tokens.IsTokenValid() {
		fmt.Println("⏰ Token expired, refreshing...")
		newTokens, err := auth.RefreshAccessToken(cfg.GetTokenURL(), cfg.CognitoConfig.ClientID, tokens.RefreshToken)
		if err != nil {
			return fmt.Errorf("failed to refresh token: %w", err)
		}
		if err := cfg.SaveTokens(newTokens); err != nil {
			return fmt.Errorf("failed to save refreshed tokens: %w", err)
		}
		tokens = newTokens
	}

	// 4. Extract owner from token
	owner, err := cfg.GetOwner()
	if err != nil {
		return fmt.Errorf("failed to extract owner from token: %w", err)
	}

	// 5. Create API client
	client := api.NewClient(cfg.GetAPIBaseURL(), tokens.AccessToken)

	// 6. Call GetEgg to get all eggs
	eggs, err := client.GetEgg(owner)
	if err != nil {
		return fmt.Errorf("failed to get eggs: %w", err)
	}

	// 7. If a specific key was provided, find and print just that one
	if len(args) == 1 {
		key := args[0]
		found := false
		for _, egg := range eggs {
			if egg.SecretID == key {
				fmt.Printf("🥚 Secret: %s\n", key)
				fmt.Printf("Value: %s\n", egg.Plaintext)
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("secret '%s' not found", key)
		}
	} else {
		// No key provided - list all secrets
		if len(eggs) == 0 {
			fmt.Println("No secrets found in your vault.")
			return nil
		}

		fmt.Printf("🥚 Found %d secret(s):\n\n", len(eggs))
		for _, egg := range eggs {
			fmt.Printf("Key: %s\n", egg.SecretID)
			fmt.Printf("Value: %s\n", egg.Plaintext)
			fmt.Printf("Created: %s\n", egg.CreatedAt)
			fmt.Println("---")
		}
	}

	return nil
}
