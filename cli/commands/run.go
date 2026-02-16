package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// RunCmd represents the run command
var RunCmd = &cobra.Command{
	Use:   "run -- [command]",
	Short: "Inject secrets and run a command",
	Long: `Fetch all secrets, set them as environment variables, and execute a command.
	
Example:
  egg run -- go run main.go
  egg run -- npm start
  egg run -- ./my-script.sh`,
	RunE: runRun,
	// DisableFlagParsing allows passing flags to the subprocess
	DisableFlagParsing: true,
}

func runRun(cmd *cobra.Command, args []string) error {
	// TODO: Implement run logic
	// 1. Load config
	// 2. Load tokens (check if logged in)
	// 3. Check if token is valid (refresh if needed)
	// 4. Create API client
	// 5. Fetch ALL secrets for the user (you may need to add a ListEggs API method)
	// 6. Parse secrets into key-value map
	// 7. Get current environment variables
	// 8. Merge secrets into environment
	// 9. Find the "--" separator in args
	// 10. Extract command and arguments after "--"
	// 11. Create exec.Command with custom environment
	// 12. Wire up stdin/stdout/stderr
	// 13. Run command and wait
	// 14. Exit with same code as subprocess

	fmt.Println("🚀 Running command with injected secrets...")
	fmt.Printf("Args: %v\n", args)

	return fmt.Errorf("not implemented - see api/client.go for fetching secrets")
}
