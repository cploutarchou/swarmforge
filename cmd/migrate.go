package cmd

import (
	"fmt"
	"strings"

	"github.com/cydevcloud/infra-cli/pkg/migration"
	"github.com/cydevcloud/infra-cli/pkg/types"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate infrastructure components",
	Long:  `Commands for migrating infrastructure components like nodes and services.`,
}

var migrateNodeCmd = &cobra.Command{
	Use:   "node [server-ip]",
	Short: "Migrate a node to a new role",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverIP := args[0]
		if serverIP == "" {
			return fmt.Errorf("server IP is required")
		}

		if username == "" {
			return fmt.Errorf("username is required")
		}

		if password == "" {
			return fmt.Errorf("password is required")
		}

		if serverRole == "" {
			return fmt.Errorf("server role is required")
		}

		// Convert string to ServerRole type
		role := types.ServerRole(serverRole)

		// Validate role
		switch role {
		case types.ManagerServer, types.GitlabServer, types.MonitorServer, types.AppsServer:
			// Valid role
		default:
			return fmt.Errorf("invalid server role: %s", serverRole)
		}

		if !force {
			fmt.Printf("WARNING: This will migrate node %s to role %s. Continue? [y/N] ", serverIP, role)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				return fmt.Errorf("migration aborted")
			}
		}

		fmt.Printf("Migrating node %s to role %s...\n", serverIP, role)
		return migration.MigrateNode(serverIP, username, password, role)
	},
}

var setupManagerCmd = &cobra.Command{
	Use:   "setup-manager",
	Short: "Set up a new manager node",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !force {
			fmt.Printf("WARNING: This will set up %s as a new manager node. Continue? [y/N] ", targetIP)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				return fmt.Errorf("setup aborted")
			}
		}

		// Get existing manager info
		if sourceIP != "" {
			output, err := executeRemoteCommand(sourceIP, username, password, "docker node ls --format '{{.Hostname}} {{.Role}}'")
			if err != nil {
				return fmt.Errorf("failed to get node list: %w", err)
			}
			fmt.Printf("Current nodes:\n%s\n", output)
		}

		// Set up new manager
		if err := migration.SetupManager(targetIP, username, password); err != nil {
			return fmt.Errorf("failed to set up manager: %w", err)
		}

		// Migrate Traefik if requested
		if migrateTraefik && sourceIP != "" {
			if err := migration.MigrateTraefik(sourceIP, targetIP, username, password); err != nil {
				return fmt.Errorf("failed to migrate Traefik: %w", err)
			}
		}

		fmt.Printf("Manager node %s successfully set up\n", targetIP)
		return nil
	},
}

var migrateTraefikCmd = &cobra.Command{
	Use:   "traefik [source-ip] [target-ip]",
	Short: "Migrate Traefik from one node to another",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceIP := args[0]
		targetIP := args[1]

		if username == "" {
			return fmt.Errorf("username is required")
		}

		if password == "" {
			return fmt.Errorf("password is required")
		}

		fmt.Printf("Migrating Traefik from %s to %s...\n", sourceIP, targetIP)
		return migration.MigrateTraefik(sourceIP, targetIP, username, password)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateNodeCmd)
	migrateCmd.AddCommand(setupManagerCmd)
	migrateCmd.AddCommand(migrateTraefikCmd)

	// Add flags for migrate node command
	migrateNodeCmd.Flags().StringVar(&serverRole, "role", "", "New role for the server (manager, gitlab, monitor, apps)")
	migrateNodeCmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")
	migrateNodeCmd.MarkFlagRequired("role")

	// Add flags for setup manager command
	setupManagerCmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")
	setupManagerCmd.Flags().StringVar(&sourceIP, "source-ip", "", "Source manager IP (for migration)")
	setupManagerCmd.Flags().StringVar(&targetIP, "target-ip", "", "Target manager IP")
	setupManagerCmd.Flags().BoolVar(&migrateTraefik, "migrate-traefik", false, "Migrate Traefik from source to target")

	// Add flags for migrate traefik command
	migrateTraefikCmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")
	migrateTraefikCmd.Flags().BoolVar(&skipBackup, "skip-backup", false, "Skip backup creation")
	migrateTraefikCmd.Flags().StringVar(&sourceIP, "source-ip", "", "Source node IP")
	migrateTraefikCmd.Flags().StringVar(&targetIP, "target-ip", "", "Target node IP")
}
