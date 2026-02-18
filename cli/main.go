package main

import (
	"fmt"
	"os"

	"github.com/owenHochwald/egg-carton/cli/commands"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "egg",
	Short: "🥚 EggCarton - Secure secret management CLI",
	Long: `EggCarton is a secure CLI tool for managing secrets with an egg theme!

Commands:
  🔐 login           - Authenticate with OAuth
  🐔 lay (add)       - Store a secret (lay an egg)
  🥚 get             - Retrieve secrets from your vault
  🐣 hatch (run)     - Inject secrets and run a command (hatch your eggs)
  💥 break           - Delete a secret from your vault

It uses AWS Lambda, DynamoDB, and KMS for encryption,
with Cognito authentication via OAuth PKCE flow.`,
}

func main() {
	// Add all subcommands
	rootCmd.AddCommand(commands.LoginCmd)
	rootCmd.AddCommand(commands.AddCmd)
	rootCmd.AddCommand(commands.GetCmd)
	rootCmd.AddCommand(commands.BreakCmd)
	rootCmd.AddCommand(commands.RunCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
