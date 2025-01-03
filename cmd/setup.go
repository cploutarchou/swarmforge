package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cydevcloud/infra-cli/pkg/template"
	"github.com/cydevcloud/infra-cli/pkg/types"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up infrastructure components",
	Long:  `Commands for setting up infrastructure components like Traefik, monitoring, etc.`,
}

var setupTraefikCmd = &cobra.Command{
	Use:   "traefik",
	Short: "Setup Traefik reverse proxy",
	Long: `Setup Traefik reverse proxy with automatic SSL certificate management.
	
Example:
  infra setup traefik --ip 192.168.1.10 --domain example.com --email admin@example.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverIP == "" || domain == "" || email == "" {
			return fmt.Errorf("server IP, domain, and email are required")
		}

		// Generate Traefik configuration
		generator := template.NewGenerator()
		config := types.DeploymentConfig{
			Domain: domain,
			Email:  email,
		}

		traefikYAML, err := generator.GenerateTraefikConfig(config)
		if err != nil {
			return fmt.Errorf("failed to generate Traefik config: %w", err)
		}

		// Save configuration
		traefikFile := "/tmp/traefik-init.yaml"
		if err := os.WriteFile(traefikFile, []byte(traefikYAML), 0644); err != nil {
			return fmt.Errorf("failed to write Traefik config: %w", err)
		}

		// Copy configuration to server
		scpCmd := fmt.Sprintf("scp %s %s@%s:/tmp/", traefikFile, username, serverIP)
		if err := exec.Command("sh", "-c", scpCmd).Run(); err != nil {
			return fmt.Errorf("failed to copy Traefik config: %w", err)
		}

		// Create Traefik network
		networkCmd := "docker network create --driver=overlay traefik-public"
		if _, err := executeRemoteCommand(serverIP, username, password, networkCmd); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return fmt.Errorf("failed to create Traefik network: %w", err)
			}
		}

		// Deploy Traefik
		deployCmd := "docker stack deploy -c /tmp/traefik-init.yaml traefik"
		result, err := executeRemoteCommand(serverIP, username, password, deployCmd)
		if err != nil {
			return fmt.Errorf("failed to deploy Traefik: %w", err)
		}

		fmt.Println("Traefik setup completed successfully")
		fmt.Printf("Traefik dashboard available at: https://traefik.%s\n", domain)
		fmt.Println(result)

		return nil
	},
}

var setupServersCmd = &cobra.Command{
	Use:   "servers",
	Short: "Setup all servers",
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverIP == "" {
			return fmt.Errorf("server IP is required")
		}

		// Update system and install dependencies
		setupCmds := []string{
			"apt-get update",
			"apt-get upgrade -y",
			"apt-get install -y curl wget git",
			"curl -fsSL https://get.docker.com | sh",
			"systemctl enable docker",
			"systemctl start docker",
		}

		for _, cmd := range setupCmds {
			command := exec.Command("sshpass", "-p", password, "ssh",
				"-o", "StrictHostKeyChecking=no",
				fmt.Sprintf("%s@%s", username, serverIP),
				cmd)

			if output, err := command.CombinedOutput(); err != nil {
				return fmt.Errorf("failed to execute command '%s': %w\n%s", cmd, err, string(output))
			}
		}

		fmt.Printf("Server %s setup completed successfully\n", serverIP)
		return nil
	},
}

var setupFirewallCmd = &cobra.Command{
	Use:   "firewall",
	Short: "Setup firewall rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverIP == "" {
			return fmt.Errorf("server IP is required")
		}

		// Setup firewall rules
		firewallCmds := []string{
			"ufw allow OpenSSH",
			"ufw allow 80/tcp",   // HTTP
			"ufw allow 443/tcp",  // HTTPS
			"ufw allow 2377/tcp", // Docker Swarm cluster management
			"ufw allow 7946/tcp", // Container network discovery
			"ufw allow 7946/udp",
			"ufw allow 4789/udp", // Container overlay network
			"ufw --force enable",
		}

		for _, cmd := range firewallCmds {
			command := exec.Command("sshpass", "-p", password, "ssh",
				"-o", "StrictHostKeyChecking=no",
				fmt.Sprintf("%s@%s", username, serverIP),
				cmd)

			if output, err := command.CombinedOutput(); err != nil {
				return fmt.Errorf("failed to execute command '%s': %w\n%s", cmd, err, string(output))
			}
		}

		fmt.Printf("Firewall setup completed successfully on %s\n", serverIP)
		return nil
	},
}

func init() {
	// Add subcommands
	setupCmd.AddCommand(setupServersCmd)
	setupCmd.AddCommand(setupFirewallCmd)
	setupCmd.AddCommand(setupTraefikCmd)

	// Add to root command
	rootCmd.AddCommand(setupCmd)

	// Add flags
	setupTraefikCmd.Flags().StringVar(&email, "email", "", "Email address for Let's Encrypt")
	setupTraefikCmd.Flags().StringVar(&domain, "domain", "", "Domain name for Traefik dashboard")
	setupTraefikCmd.MarkFlagRequired("email")
	setupTraefikCmd.MarkFlagRequired("domain")
}
