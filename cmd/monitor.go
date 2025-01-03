package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"path/filepath"

	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitoring commands",
	Long:  `Commands for monitoring system health, services, and resources.`,
}

var checkHealthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check system health",
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverIP == "" {
			return fmt.Errorf("server IP is required")
		}

		// Commands to check system health
		healthChecks := map[string]string{
			"CPU Usage":     "top -bn1 | grep 'Cpu(s)' | awk '{print $2}'",
			"Memory Usage":  "free -m | awk 'NR==2{printf \"%.2f%%\", $3*100/$2}'",
			"Disk Usage":    "df -h / | awk 'NR==2{print $5}'",
			"Docker Status": "systemctl status docker | grep Active",
		}

		for check, command := range healthChecks {
			fmt.Printf("\nChecking %s...\n", check)
			sshCmd := exec.Command("sshpass", "-p", password, "ssh",
				"-o", "StrictHostKeyChecking=no",
				fmt.Sprintf("%s@%s", username, serverIP),
				command)
			
			output, err := sshCmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Error checking %s: %v\n", check, err)
				continue
			}
			fmt.Printf("%s: %s\n", check, strings.TrimSpace(string(output)))
		}

		return nil
	},
}

var checkServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Check services status",
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverIP == "" {
			return fmt.Errorf("server IP is required")
		}

		// Commands to check services
		serviceChecks := map[string]string{
			"Docker Services": "docker service ls",
			"Running Containers": "docker ps",
			"Service Logs": "docker service logs $(docker service ls -q) --tail 10",
		}

		for check, command := range serviceChecks {
			fmt.Printf("\n=== %s ===\n", check)
			sshCmd := exec.Command("sshpass", "-p", password, "ssh",
				"-o", "StrictHostKeyChecking=no",
				fmt.Sprintf("%s@%s", username, serverIP),
				command)
			
			output, err := sshCmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Error checking %s: %v\n", check, err)
				continue
			}
			fmt.Println(string(output))
		}

		return nil
	},
}

var setupMonitoringCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup monitoring stack",
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverIP == "" {
			return fmt.Errorf("server IP is required")
		}

		// Copy monitoring stack files
		stackDir := filepath.Join("stacks", "monitoring")
		copyCmd := exec.Command("sshpass", "-p", password, "scp", "-r",
			"-o", "StrictHostKeyChecking=no",
			stackDir,
			fmt.Sprintf("%s@%s:/root/", username, serverIP))
		
		if output, err := copyCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to copy monitoring stack: %w\n%s", err, string(output))
		}

		// Deploy monitoring stack
		deployCmd := exec.Command("sshpass", "-p", password, "ssh",
			"-o", "StrictHostKeyChecking=no",
			fmt.Sprintf("%s@%s", username, serverIP),
			"cd /root/monitoring && docker stack deploy -c docker-compose.yml monitoring")
		
		if output, err := deployCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to deploy monitoring stack: %w\n%s", err, string(output))
		}

		fmt.Println("Monitoring stack deployed successfully")
		return nil
	},
}

func init() {
	// Add subcommands
	monitorCmd.AddCommand(checkHealthCmd)
	monitorCmd.AddCommand(checkServicesCmd)
	monitorCmd.AddCommand(setupMonitoringCmd)

	// Add to root command
	rootCmd.AddCommand(monitorCmd)
}
