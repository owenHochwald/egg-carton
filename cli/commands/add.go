package commands

import (
	"fmt"

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

	// TODO: Implement add logic
	// 1. Load config
	// 2. Load tokens (check if logged in)
	// 3. Check if token is valid (refresh if needed)
	// 4. Create API client
	// 5. Call PutEgg(key, value)
	// 6. Print success message

	_ = key
	_ = value

	return fmt.Errorf("not implemented - see api/client.go")
}
