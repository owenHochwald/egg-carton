package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// GetCmd represents the get command
var GetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Retrieve a secret",
	Long:  `Decrypt and retrieve a secret from your EggCarton vault.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runGet,
}

func runGet(cmd *cobra.Command, args []string) error {
	key := args[0]

	// TODO: Implement get logic
	// 1. Load config
	// 2. Load tokens (check if logged in)
	// 3. Check if token is valid (refresh if needed)
	// 4. Create API client
	// 5. Call GetEgg(key)
	// 6. Print the decrypted value to stdout

	_ = key

	return fmt.Errorf("not implemented - see api/client.go")
}
