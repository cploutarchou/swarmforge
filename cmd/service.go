package cmd

import (
	"fmt"

	"github.com/cydevcloud/infra-cli/pkg/types"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage services",
	Long:  `Commands for managing services in the swarm.`,
}

var createServiceCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("service name is required")
		}

		serviceName := args[0]
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

		// Add service-specific configuration here
		fmt.Printf("Creating service %s with config: %+v\n", serviceName, config)
		return nil
	},
}

var listServiceCmd = &cobra.Command{
	Use:   "list",
	Short: "List services",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Listing services...")
		return nil
	},
}

var deleteServiceCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("service name is required")
		}

		serviceName := args[0]
		fmt.Printf("Deleting service %s...\n", serviceName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(createServiceCmd, listServiceCmd, deleteServiceCmd)
}
