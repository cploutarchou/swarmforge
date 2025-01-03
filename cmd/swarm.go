package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/cydevcloud/infra-cli/pkg/types"
	"github.com/spf13/cobra"
)

var (
	workerToken   string
	managerToken  string
	advertiseAddr string
)

var swarmCmd = &cobra.Command{
	Use:   "swarm",
	Short: "Manage Docker Swarm cluster",
	Long:  `Commands for managing Docker Swarm cluster, including initialization, joining nodes, and viewing statistics.`,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new swarm",
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverIP == "" {
			return fmt.Errorf("server IP is required")
		}

		// Verify this is a manager node
		if serverRole != string(types.ManagerServer) {
			return fmt.Errorf("swarm can only be initialized on a manager node")
		}

		// Initialize swarm
		initCmd := fmt.Sprintf("docker swarm init --advertise-addr %s", serverIP)
		result, err := executeRemoteCommand(serverIP, username, password, initCmd)
		if err != nil {
			return fmt.Errorf("failed to initialize swarm: %w", err)
		}

		fmt.Println("Swarm initialized successfully")
		fmt.Println(result)

		// Get tokens
		workerToken, err = getSwarmToken(serverIP, username, password, "worker")
		if err != nil {
			return fmt.Errorf("failed to get worker token: %w", err)
		}

		managerToken, err = getSwarmToken(serverIP, username, password, "manager")
		if err != nil {
			return fmt.Errorf("failed to get manager token: %w", err)
		}

		// Apply manager node labels
		labels := types.GetServerLabels(types.ManagerServer)
		for key, value := range labels {
			labelCmd := fmt.Sprintf("docker node update --label-add %s=%s $(docker node ls --format '{{.ID}}')", key, value)
			if _, err := executeRemoteCommand(serverIP, username, password, labelCmd); err != nil {
				return fmt.Errorf("failed to apply labels: %w", err)
			}
		}

		fmt.Printf("Worker join token: %s\n", workerToken)
		fmt.Printf("Manager join token: %s\n", managerToken)
		return nil
	},
}

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join a node to the swarm",
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverIP == "" || managerIP == "" {
			return fmt.Errorf("server IP and manager IP are required")
		}

		if !types.IsValidServerRole(serverRole) {
			return fmt.Errorf("invalid server role. Valid roles are: %v", types.ValidServerRoles())
		}

		// Get join token based on role
		var token string
		var err error
		if serverRole == string(types.ManagerServer) {
			token, err = getSwarmToken(managerIP, username, password, "manager")
		} else {
			token, err = getSwarmToken(managerIP, username, password, "worker")
		}
		if err != nil {
			return fmt.Errorf("failed to get join token: %w", err)
		}

		// Join swarm
		joinCmd := fmt.Sprintf("docker swarm join --token %s %s:2377", token, managerIP)
		result, err := executeRemoteCommand(serverIP, username, password, joinCmd)
		if err != nil {
			return fmt.Errorf("failed to join swarm: %w", err)
		}

		fmt.Printf("Node %s joined the swarm successfully\n", serverIP)
		fmt.Println(result)

		// Apply role-specific labels
		role := types.ServerRole(serverRole)
		labels := types.GetServerLabels(role)
		for key, value := range labels {
			labelCmd := fmt.Sprintf("docker node update --label-add %s=%s %s", key, value, serverIP)
			if _, err := executeRemoteCommand(managerIP, username, password, labelCmd); err != nil {
				return fmt.Errorf("failed to apply labels: %w", err)
			}
		}

		// Setup role-specific configurations
		switch role {
		case types.GitlabServer:
			if err := setupGitlabNode(serverIP, username, password); err != nil {
				return fmt.Errorf("failed to setup Gitlab node: %w", err)
			}
		case types.MonitorServer:
			if err := setupMonitorNode(serverIP, username, password); err != nil {
				return fmt.Errorf("failed to setup Monitor node: %w", err)
			}
		case types.AppsServer:
			if err := setupAppsNode(serverIP, username, password); err != nil {
				return fmt.Errorf("failed to setup Apps node: %w", err)
			}
		}

		fmt.Printf("Node labeled and configured with role: %s\n", serverRole)
		return nil
	},
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show swarm statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverIP == "" {
			return fmt.Errorf("server IP is required")
		}

		// Get nodes info with role labels
		nodesCmd := "docker node ls --format '{{.ID}}\t{{.Hostname}}\t{{.Status}}\t{{.Availability}}\t{{.ManagerStatus}}\t{{.Labels}}'"
		nodesResult, err := executeRemoteCommand(serverIP, username, password, nodesCmd)
		if err != nil {
			return fmt.Errorf("failed to get nodes info: %w", err)
		}

		// Get services info
		servicesCmd := "docker service ls"
		servicesResult, err := executeRemoteCommand(serverIP, username, password, servicesCmd)
		if err != nil {
			return fmt.Errorf("failed to get services info: %w", err)
		}

		fmt.Println("\nSwarm Nodes:")
		fmt.Println("ID\tHOSTNAME\tSTATUS\tAVAILABILITY\tMANAGER STATUS\tLABELS")
		fmt.Println(nodesResult)
		fmt.Println("\nSwarm Services:")
		fmt.Println(servicesResult)

		return nil
	},
}

func setupGitlabNode(ip, user, pass string) error {
	// Setup Gitlab-specific requirements
	cmds := []string{
		"apt-get update",
		"apt-get install -y ca-certificates curl openssh-server",
		"mkdir -p /etc/gitlab/config",
		"mkdir -p /var/log/gitlab",
		"mkdir -p /var/opt/gitlab",
	}

	for _, cmd := range cmds {
		if _, err := executeRemoteCommand(ip, user, pass, cmd); err != nil {
			return err
		}
	}
	return nil
}

func setupMonitorNode(ip, user, pass string) error {
	// Setup monitoring-specific requirements
	cmds := []string{
		"apt-get update",
		"apt-get install -y prometheus-node-exporter",
		"mkdir -p /etc/prometheus",
		"mkdir -p /var/lib/prometheus",
	}

	for _, cmd := range cmds {
		if _, err := executeRemoteCommand(ip, user, pass, cmd); err != nil {
			return err
		}
	}
	return nil
}

func setupAppsNode(ip, user, pass string) error {
	// Setup application node requirements
	cmds := []string{
		"apt-get update",
		"apt-get install -y docker-compose-plugin",
		"mkdir -p /app/data",
		"mkdir -p /app/config",
	}

	for _, cmd := range cmds {
		if _, err := executeRemoteCommand(ip, user, pass, cmd); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	// Add subcommands
	swarmCmd.AddCommand(initCmd)
	swarmCmd.AddCommand(joinCmd)
	swarmCmd.AddCommand(statsCmd)

	// Add to root command
	rootCmd.AddCommand(swarmCmd)

	// Add flags
	joinCmd.Flags().StringVar(&managerIP, "manager-ip", "", "Manager node IP address")
	joinCmd.Flags().StringVar(&advertiseAddr, "advertise-addr", "", "Advertise address (format: <ip|interface>[:port])")
}

func getSwarmToken(ip, user, pass, role string) (string, error) {
	cmd := fmt.Sprintf("docker swarm join-token -q %s", role)
	result, err := executeRemoteCommand(ip, user, pass, cmd)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result), nil
}

func executeRemoteCommand(ip, user, pass, command string) (string, error) {
	sshCmd := exec.Command("sshpass", "-p", pass, "ssh", "-o", "StrictHostKeyChecking=no",
		fmt.Sprintf("%s@%s", user, ip), command)

	output, err := sshCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %w\nOutput: %s", err, string(output))
	}

	return string(output), nil
}
