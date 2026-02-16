package commands

import (
	"fmt"

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

	// TODO: Implement break logic
	// 1. Load config
	// 2. Load tokens (check if logged in)
	// 3. Check if token is valid (refresh if needed)
	// 4. Create API client
	// 5. Call BreakEgg(key)
	// 6. Print confirmation message

	_ = key

	return fmt.Errorf("not implemented - see api/client.go")
}
