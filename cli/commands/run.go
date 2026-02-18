package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/owenHochwald/egg-carton/cli/api"
	"github.com/owenHochwald/egg-carton/cli/auth"
	"github.com/owenHochwald/egg-carton/cli/config"
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

	// 5. Create API client and fetch ALL secrets
	client := api.NewClient(cfg.GetAPIBaseURL(), tokens.AccessToken)
	eggs, err := client.GetEgg(owner)
	if err != nil {
		return fmt.Errorf("failed to get eggs: %w", err)
	}

	// 6. Parse secrets into environment variables
	secretEnvVars := make(map[string]string)
	for _, egg := range eggs {
		// Convert secret_id to uppercase env var format (e.g., api_key -> API_KEY)
		envVarName := strings.ToUpper(egg.SecretID)
		secretEnvVars[envVarName] = egg.Plaintext
	}

	// 7. Find the "--" separator in args
	dashIndex := -1
	for i, arg := range args {
		if arg == "--" {
			dashIndex = i
			break
		}
	}

	if dashIndex == -1 || dashIndex == len(args)-1 {
		return fmt.Errorf("usage: egg run -- <command> [args...]")
	}

	// 8. Extract command and arguments after "--"
	commandArgs := args[dashIndex+1:]
	if len(commandArgs) == 0 {
		return fmt.Errorf("no command specified after '--'")
	}

	commandName := commandArgs[0]
	commandArguments := commandArgs[1:]

	// 9. Get current environment variables
	currentEnv := os.Environ()

	// 10. Merge secrets into environment
	mergedEnv := append([]string{}, currentEnv...)
	for key, value := range secretEnvVars {
		mergedEnv = append(mergedEnv, fmt.Sprintf("%s=%s", key, value))
	}

	fmt.Printf("🚀 Running command with %d secret(s) injected...\n", len(secretEnvVars))
	for key := range secretEnvVars {
		fmt.Printf("   ✓ %s\n", key)
	}
	fmt.Println()

	// 11. Create exec.Command with custom environment
	command := exec.Command(commandName, commandArguments...)
	command.Env = mergedEnv
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// 12. Run command and wait
	if err := command.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// 13. Exit with same code as subprocess
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("failed to run command: %w", err)
	}

	return nil
}
