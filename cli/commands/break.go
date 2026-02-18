package commands

import (
	"fmt"

	"github.com/owenHochwald/egg-carton/cli/api"
	"github.com/owenHochwald/egg-carton/cli/auth"
	"github.com/owenHochwald/egg-carton/cli/config"
	"github.com/spf13/cobra"
)

// BreakCmd represents the break command
var BreakCmd = &cobra.Command{
	Use:   "break [key]",
	Short: "Delete a secret",
	Long:  `Permanently delete a secret from your EggCarton vault.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runBreak,
}

func runBreak(cmd *cobra.Command, args []string) error {
	key := args[0]

	fmt.Printf("💥 Breaking egg: %s\n", key)

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

	// 6. Call BreakEgg(owner, secretID)
	if err := client.BreakEgg(owner, key); err != nil {
		return fmt.Errorf("failed to break egg: %w", err)
	}

	// 7. Print confirmation message
	fmt.Printf("✅ Successfully deleted secret: %s\n", key)

	return nil
}
