package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "System management commands",
	Long:  `Commands for managing system configuration, updates, and security.`,
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update system packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverIP == "" {
			return fmt.Errorf("server IP is required")
		}

		// Update system packages
		updateCmd := `apt-get update && apt-get upgrade -y`
		command := exec.Command("sshpass", "-p", password, "ssh",
			"-o", "StrictHostKeyChecking=no",
			fmt.Sprintf("%s@%s", username, serverIP),
			updateCmd)

		output, err := command.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to update system: %w\n%s", err, string(output))
		}

		fmt.Printf("System updated successfully on %s\n", serverIP)
		return nil
	},
}

var hardenCmd = &cobra.Command{
	Use:   "harden",
	Short: "Apply security hardening",
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverIP == "" {
			return fmt.Errorf("server IP is required")
		}

		// Basic security hardening commands
		hardenCmds := []string{
			"ufw enable",
			"ufw allow ssh",
			"ufw allow http",
			"ufw allow https",
			"sed -i 's/PermitRootLogin yes/PermitRootLogin prohibit-password/' /etc/ssh/sshd_config",
			"systemctl restart sshd",
		}

		for _, cmd := range hardenCmds {
			command := exec.Command("sshpass", "-p", password, "ssh",
				"-o", "StrictHostKeyChecking=no",
				fmt.Sprintf("%s@%s", username, serverIP),
				cmd)

			output, err := command.CombinedOutput()
			if err != nil {
				return fmt.Errorf("failed to execute command '%s': %w\n%s", cmd, err, string(output))
			}
		}

		fmt.Printf("Security hardening applied successfully on %s\n", serverIP)
		return nil
	},
}

func init() {
	// Add subcommands
	systemCmd.AddCommand(updateCmd)
	systemCmd.AddCommand(hardenCmd)

	// Add to root command
	rootCmd.AddCommand(systemCmd)
}
