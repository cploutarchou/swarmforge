package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cydevcloud/infra-cli/pkg/template"
	"github.com/cydevcloud/infra-cli/pkg/types"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy services to the swarm",
	Long: `Deploy services to the Docker Swarm cluster with optional Traefik routing.
Example: infra deploy --name api --type api --lang go --port 8080 --domain example.com --subdomain api`,
}

var deployAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Deploy entire infrastructure",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check swarm status
		swarmCmd := exec.Command("./scripts/check-swarm.sh")
		if err := swarmCmd.Run(); err == nil {
			fmt.Println("Swarm already configured, redeploying services...")
			return deployServices()
		}

		fmt.Println("Setting up new swarm...")
		if err := setupSwarm(); err != nil {
			return fmt.Errorf("failed to setup swarm: %w", err)
		}

		return deployServices()
	},
}

var deployServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Deploy all services",
	RunE: func(cmd *cobra.Command, args []string) error {
		return deployServices()
	},
}

var deployStackCmd = &cobra.Command{
	Use:   "stack [name] [file]",
	Short: "Deploy a stack from a compose file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]
		composeFile := args[1]
		// TODO: Implement stack deployment logic
		fmt.Printf("Deploying stack %s from file %s\n", stackName, composeFile)
		return nil
	},
}

var deployServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Deploy a service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverIP == "" {
			return fmt.Errorf("server IP is required")
		}

		// Generate deployment configuration
		config := types.DeploymentConfig{
			ServiceName: serviceName,
			ImageName:   fmt.Sprintf("%s-%s", serviceName, appType),
			Port:        port,
			Replicas:    replicas,
			Domain:      domain,
			Subdomain:   subdomain,
			UseTraefik:  useTraefik,
			Labels:      make(map[string]string),
			Environment: map[string]string{
				"SERVICE_NAME": serviceName,
				"APP_PORT":     fmt.Sprintf("%d", port),
			},
		}

		// Add Traefik labels if enabled
		if useTraefik {
			if domain == "" || subdomain == "" {
				return fmt.Errorf("domain and subdomain are required when using Traefik")
			}

			config.Labels = map[string]string{
				"traefik.enable": "true",
				"traefik.http.routers." + serviceName + ".rule":                      fmt.Sprintf("Host(`%s.%s`)", subdomain, domain),
				"traefik.http.routers." + serviceName + ".service":                   serviceName,
				"traefik.http.services." + serviceName + ".loadbalancer.server.port": fmt.Sprintf("%d", port),
				"traefik.http.routers." + serviceName + ".entrypoints":               "websecure",
				"traefik.http.routers." + serviceName + ".tls.certresolver":          "letsencrypt",
			}
		}

		// Generate deployment files
		generator := template.NewGenerator()

		// Generate main deployment
		deploymentYAML, err := generator.GenerateDeployment(config)
		if err != nil {
			return fmt.Errorf("failed to generate deployment: %w", err)
		}

		// Generate Traefik config if enabled
		var traefikYAML string
		if useTraefik {
			traefikYAML, err = generator.GenerateTraefikConfig(config)
			if err != nil {
				return fmt.Errorf("failed to generate Traefik config: %w", err)
			}
		}

		// Save deployment files
		deployDir := filepath.Join("/tmp", serviceName)
		if err := os.MkdirAll(deployDir, 0755); err != nil {
			return fmt.Errorf("failed to create deployment directory: %w", err)
		}

		deploymentFile := filepath.Join(deployDir, "deployment.yaml")
		if err := os.WriteFile(deploymentFile, []byte(deploymentYAML), 0644); err != nil {
			return fmt.Errorf("failed to write deployment file: %w", err)
		}

		if useTraefik {
			traefikFile := filepath.Join(deployDir, "traefik.yaml")
			if err := os.WriteFile(traefikFile, []byte(traefikYAML), 0644); err != nil {
				return fmt.Errorf("failed to write Traefik config: %w", err)
			}
		}

		// Copy files to server
		scpCmd := fmt.Sprintf("scp -r %s %s@%s:/tmp/", deployDir, username, serverIP)
		if err := exec.Command("sh", "-c", scpCmd).Run(); err != nil {
			return fmt.Errorf("failed to copy deployment files: %w", err)
		}

		// Deploy the service
		deployCmd := fmt.Sprintf("docker stack deploy -c /tmp/%s/deployment.yaml %s", serviceName, serviceName)
		if useTraefik {
			deployCmd += fmt.Sprintf(" && docker stack deploy -c /tmp/%s/traefik.yaml traefik", serviceName)
		}

		result, err := executeRemoteCommand(serverIP, username, password, deployCmd)
		if err != nil {
			return fmt.Errorf("failed to deploy service: %w", err)
		}

		fmt.Printf("Service %s deployed successfully\n", serviceName)
		if useTraefik {
			fmt.Printf("Service available at: https://%s.%s\n", subdomain, domain)
		}
		fmt.Println(result)

		return nil
	},
}

func deployServices() error {
	if serverIP == "" {
		return fmt.Errorf("server IP is required")
	}

	// Copy service stack files
	stacks := []string{
		"gitlab-stack.yaml",
		"monitoring-stack.yaml",
		"apps-stack.yaml",
	}

	for _, stack := range stacks {
		sourcePath := filepath.Join("deployments", "templates", stack)
		scpCmd := exec.Command("sshpass", "-p", password, "scp",
			"-o", "StrictHostKeyChecking=no",
			sourcePath,
			fmt.Sprintf("%s@%s:/root/", username, serverIP))

		if output, err := scpCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to copy %s: %w\n%s", stack, err, string(output))
		}

		// Deploy stack
		deployCmd := exec.Command("sshpass", "-p", password, "ssh",
			"-o", "StrictHostKeyChecking=no",
			fmt.Sprintf("%s@%s", username, serverIP),
			fmt.Sprintf("cd /root && docker stack deploy -c %s %s", stack, stack[:len(stack)-5]))

		if output, err := deployCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to deploy %s: %w\n%s", stack, err, string(output))
		}
	}

	fmt.Println("Services deployed successfully")
	return nil
}

func setupSwarm() error {
	if serverIP == "" {
		return fmt.Errorf("server IP is required")
	}

	// Initialize swarm
	initCmd := exec.Command("sshpass", "-p", password, "ssh",
		"-o", "StrictHostKeyChecking=no",
		fmt.Sprintf("%s@%s", username, serverIP),
		"docker swarm init")

	if output, err := initCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to initialize swarm: %w\n%s", err, string(output))
	}

	fmt.Println("Swarm initialized successfully")
	return nil
}

func init() {
	// Add subcommands
	deployCmd.AddCommand(deployAllCmd)
	deployCmd.AddCommand(deployServicesCmd)
	deployCmd.AddCommand(deployStackCmd)
	deployCmd.AddCommand(deployServiceCmd)

	// Add to root command
	rootCmd.AddCommand(deployCmd)

	// Add flags
	deployServiceCmd.Flags().StringVar(&serviceName, "name", "", "Service name")
	deployServiceCmd.Flags().StringVar(&appType, "type", "", "Application type (api, standalone)")
	deployServiceCmd.Flags().StringVar(&appLang, "lang", "", "Application language")
	deployServiceCmd.Flags().IntVar(&port, "port", 8080, "Service port")
	deployServiceCmd.Flags().IntVar(&replicas, "replicas", 1, "Number of replicas")
	deployServiceCmd.Flags().StringVar(&domain, "domain", "", "Domain name for Traefik routing")
	deployServiceCmd.Flags().StringVar(&subdomain, "subdomain", "", "Subdomain for Traefik routing")
	deployServiceCmd.Flags().BoolVar(&useTraefik, "use-traefik", false, "Enable Traefik routing")

	deployServiceCmd.MarkFlagRequired("name")
	deployServiceCmd.MarkFlagRequired("type")
	deployServiceCmd.MarkFlagRequired("lang")
}
