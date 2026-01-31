package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/thepsadmin/mathemcli/internal/api"
	"github.com/thepsadmin/mathemcli/internal/config"
	"golang.org/x/term"
)

var (
	loginEmail    string
	loginPassword string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Mathem",
	Long:  `Authenticate with your Mathem account using email and password.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		email := loginEmail
		password := loginPassword

		// Prompt for email if not provided
		if email == "" {
			fmt.Print("Email: ")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read email: %w", err)
			}
			email = strings.TrimSpace(input)
		}

		// Prompt for password if not provided
		if password == "" {
			fmt.Print("Password: ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}
			fmt.Println() // Add newline after password input
			password = string(bytePassword)
		}

		// Create client and attempt login
		c := api.NewClient()
		if err := c.Login(email, password); err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		// Save session
		session := &config.Session{
			SessionID: c.SessionID(),
			CSRFToken: c.CSRFToken(),
			Email:     email,
		}
		if err := config.SaveSession(session); err != nil {
			return fmt.Errorf("failed to save session: %w", err)
		}

		fmt.Printf("Successfully logged in as %s\n", email)
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Mathem",
	Long:  `Clear the saved session and logout.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.ClearSession(); err != nil {
			return fmt.Errorf("failed to clear session: %w", err)
		}
		fmt.Println("Logged out successfully")
		return nil
	},
}

func init() {
	loginCmd.Flags().StringVarP(&loginEmail, "email", "e", "", "Email address")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "Password (not recommended, use prompt instead)")
}
