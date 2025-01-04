package cmd

import (
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/cploutarchou/swarmforge/pkg/auth"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication credentials",
	Long: `Commands for managing authentication credentials for servers.
The credentials are stored securely in an encrypted SQLite database.`,
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Store server credentials",
	Long: `Store server credentials in the encrypted database.
	
Example:
  infra auth login --server 192.168.1.10 --user root --role manager`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverIP == "" || username == "" {
			return fmt.Errorf("server and username are required")
		}

		fmt.Print("Enter password: ")
		passBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		fmt.Println()

		fmt.Print("Enter master key for encryption: ")
		masterBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read master key: %w", err)
		}
		fmt.Println()

		store, err := auth.NewCredentialStore(string(masterBytes))
		if err != nil {
			return fmt.Errorf("failed to initialize credential store: %w", err)
		}
		defer store.Close()

		creds := auth.Credentials{
			Server:   serverIP,
			Username: username,
			Password: string(passBytes),
			Role:     serverRole,
		}

		if err := store.SaveCredentials(creds); err != nil {
			return fmt.Errorf("failed to save credentials: %w", err)
		}

		fmt.Printf("Credentials stored successfully for %s@%s\n", username, serverIP)
		return nil
	},
}

var listCredsCmd = &cobra.Command{
	Use:   "list",
	Short: "List stored credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print("Enter master key for decryption: ")
		masterBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read master key: %w", err)
		}
		fmt.Println()

		store, err := auth.NewCredentialStore(string(masterBytes))
		if err != nil {
			return fmt.Errorf("failed to initialize credential store: %w", err)
		}
		defer store.Close()

		creds, err := store.ListCredentials()
		if err != nil {
			return fmt.Errorf("failed to list credentials: %w", err)
		}

		fmt.Println("\nStored credentials:")
		for _, cred := range creds {
			fmt.Printf("Server: %s, Username: %s, Role: %s\n", cred.Server, cred.Username, cred.Role)
		}

		return nil
	},
}

var deleteCredsCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete stored credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverIP == "" || username == "" {
			return fmt.Errorf("server and username are required")
		}

		fmt.Print("Enter master key for decryption: ")
		masterBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read master key: %w", err)
		}
		fmt.Println()

		store, err := auth.NewCredentialStore(string(masterBytes))
		if err != nil {
			return fmt.Errorf("failed to initialize credential store: %w", err)
		}
		defer store.Close()

		if err := store.DeleteCredentials(serverIP, username); err != nil {
			return fmt.Errorf("failed to delete credentials: %w", err)
		}

		fmt.Printf("Credentials deleted successfully for %s@%s\n", username, serverIP)
		return nil
	},
}

func init() {
	// Add subcommands
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(listCredsCmd)
	authCmd.AddCommand(deleteCredsCmd)

	// Add to root command
	rootCmd.AddCommand(authCmd)

	// Add flags
	authCmd.PersistentFlags().StringVar(&serverIP, "server", "", "Server IP address")
	authCmd.PersistentFlags().StringVar(&username, "user", "", "Username")
	authCmd.PersistentFlags().StringVar(&serverRole, "role", "", "Server role")
}
