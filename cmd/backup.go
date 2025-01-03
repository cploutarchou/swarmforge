package cmd

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup management commands",
	Long:  `Commands for managing backups of services like GitLab.`,
}

var (
	backupPath string
)

var createBackupCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new backup",
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverIP == "" {
			return fmt.Errorf("server IP is required")
		}

		timestamp := time.Now().Format("20060102150405")
		backupCmd := fmt.Sprintf("docker exec gitlab gitlab-backup create BACKUP=%s", timestamp)
		
		command := exec.Command("sshpass", "-p", password, "ssh",
			"-o", "StrictHostKeyChecking=no",
			fmt.Sprintf("%s@%s", username, serverIP),
			backupCmd)

		output, err := command.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to create backup: %w\n%s", err, string(output))
		}

		fmt.Printf("Backup created successfully with timestamp %s\n", timestamp)
		return nil
	},
}

var restoreBackupCmd = &cobra.Command{
	Use:   "restore [timestamp]",
	Short: "Restore from backup",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverIP == "" {
			return fmt.Errorf("server IP is required")
		}

		timestamp := args[0]
		restoreCmd := fmt.Sprintf("docker exec gitlab gitlab-backup restore BACKUP=%s", timestamp)
		
		command := exec.Command("sshpass", "-p", password, "ssh",
			"-o", "StrictHostKeyChecking=no",
			fmt.Sprintf("%s@%s", username, serverIP),
			restoreCmd)

		output, err := command.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to restore backup: %w\n%s", err, string(output))
		}

		fmt.Printf("Backup %s restored successfully\n", timestamp)
		return nil
	},
}

var listBackupsCmd = &cobra.Command{
	Use:   "list",
	Short: "List available backups",
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverIP == "" {
			return fmt.Errorf("server IP is required")
		}

		listCmd := "ls -l /var/opt/gitlab/backups/"
		command := exec.Command("sshpass", "-p", password, "ssh",
			"-o", "StrictHostKeyChecking=no",
			fmt.Sprintf("%s@%s", username, serverIP),
			listCmd)

		output, err := command.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to list backups: %w\n%s", err, string(output))
		}

		fmt.Printf("Available backups:\n%s", string(output))
		return nil
	},
}

func init() {
	// Add flags
	backupCmd.PersistentFlags().StringVar(&backupPath, "path", "/var/opt/gitlab/backups", "Backup directory path")

	// Add subcommands
	backupCmd.AddCommand(createBackupCmd)
	backupCmd.AddCommand(restoreBackupCmd)
	backupCmd.AddCommand(listBackupsCmd)

	// Add to root command
	rootCmd.AddCommand(backupCmd)
}
