package commands

import (
	"fmt"

	"github.com/owenHochwald/egg-carton/cli/api"
	"github.com/owenHochwald/egg-carton/cli/auth"
	"github.com/owenHochwald/egg-carton/cli/config"
	"github.com/spf13/cobra"
)

// AddCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add [key] [value]",
	Short: "Store a secret",
	Long:  `Encrypt and store a secret in your EggCarton vault.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runAdd,
}

func runAdd(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	fmt.Printf("🥚 Adding secret: %s\n", key)

	config, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	tokens, err := config.LoadTokens()
	if err != nil {
		return fmt.Errorf("failed to load tokens: %w", err)
	}
	// 3. Check if token is valid (refresh if needed)
	if tokens == nil {
		return fmt.Errorf("you are not logged in. Please run 'egg login' first")
	}
	if !tokens.IsTokenValid() {
		fmt.Println("⏰ Token expired, refreshing...")
		newTokens, err := auth.RefreshAccessToken(config.GetTokenURL(), config.CognitoConfig.ClientID, tokens.RefreshToken)
		if err != nil {
			return fmt.Errorf("failed to refresh token: %w", err)
		}
		if err := config.SaveTokens(newTokens); err != nil {
			return fmt.Errorf("failed to save refreshed tokens: %w", err)
		}
		tokens = newTokens
	}

	// 4. Extract owner from token
	owner, err := config.GetOwner()
	if err != nil {
		return fmt.Errorf("failed to extract owner from token: %w", err)
	}

	// 5. Create API client
	client := api.NewClient(config.GetAPIBaseURL(), tokens.AccessToken)

	// 6. Call PutEgg(owner, key, value)
	if err := client.PutEgg(owner, key, value); err != nil {
		return fmt.Errorf("failed to put egg: %w", err)
	}

	// 7. Print success message
	fmt.Printf("✅ Successfully added secret: %s\n", key)

	return nil
}
