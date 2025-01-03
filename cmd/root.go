package cmd

import (
	"github.com/spf13/cobra"
)

// Common variables used across commands
var (
	// Server configuration
	serverIP   string
	username   string
	password   string
	serverRole string

	// Service configuration
	serviceName string
	appType     string
	appLang     string
	port        int
	replicas    int

	// Domain configuration
	domain    string
	subdomain string
	email     string

	// Migration configuration
	sourceIP       string
	targetIP       string
	managerIP      string
	migrateTraefik bool

	// Common flags
	force      bool
	skipBackup bool
	useTraefik bool
)

var rootCmd = &cobra.Command{
	Use:   "infra",
	Short: "Infrastructure management CLI",
	Long: `A CLI tool for managing infrastructure services, deployments, and Docker Swarm operations.
Complete documentation is available at https://github.com/yourusername/infrastructure-setup`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add persistent flags that will be available to all commands
	rootCmd.PersistentFlags().StringVar(&serverIP, "ip", "", "Server IP address")
	rootCmd.PersistentFlags().StringVar(&username, "user", "root", "SSH username")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "SSH password")
	rootCmd.PersistentFlags().StringVar(&serverRole, "role", "", "Server role (manager, gitlab, monitor, apps)")
}
