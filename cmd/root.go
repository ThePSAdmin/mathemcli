package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thepsadmin/mathemcli/internal/api"
	"github.com/thepsadmin/mathemcli/internal/config"
)

var (
	client  *api.Client
	rootCmd = &cobra.Command{
		Use:   "mathemcli",
		Short: "CLI for interacting with the Mathem grocery API",
		Long: `mathemcli is a command-line tool for searching products
and managing your shopping cart on Mathem.se.

Before using most commands, you need to login:
  mathemcli login`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip client setup for login and help commands
			if cmd.Name() == "login" || cmd.Name() == "help" || cmd.Name() == "version" {
				return nil
			}

			// Load saved session
			session, err := config.LoadSession()
			if err != nil {
				return fmt.Errorf("failed to load session: %w", err)
			}

			if session == nil || session.SessionID == "" {
				return fmt.Errorf("not logged in. Run 'mathemcli login' first")
			}

			client = api.NewClientWithSession(session.SessionID, session.CSRFToken)
			return nil
		},
	}
)

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(cartCmd)
	rootCmd.AddCommand(versionCmd)
}
