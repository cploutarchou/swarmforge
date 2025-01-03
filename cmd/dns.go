package cmd

import (
	"fmt"
	"os/exec"

	"github.com/cydevcloud/infra-cli/pkg/dns"
	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Manage DNS records",
	Long:  `Commands for managing DNS records for your infrastructure.`,
}

var updateDNSCmd = &cobra.Command{
	Use:   "update",
	Short: "Update DNS records",
	RunE: func(cmd *cobra.Command, args []string) error {
		return dns.UpdateDNSRecord(domain, subdomain, serverIP)
	},
}

var deleteDNSCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete DNS records",
	RunE: func(cmd *cobra.Command, args []string) error {
		return dns.DeleteDNSRecord(domain, subdomain)
	},
}

var listDNSCmd = &cobra.Command{
	Use:   "list",
	Short: "List DNS records",
	RunE: func(cmd *cobra.Command, args []string) error {
		return dns.ListDNSRecords(domain)
	},
}

var verifyDNSCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify DNS records",
	RunE: func(cmd *cobra.Command, args []string) error {
		if domain == "" {
			return fmt.Errorf("domain is required")
		}

		// Verify DNS records
		verifyCmd := fmt.Sprintf("dig +short %s", domain)
		command := exec.Command("bash", "-c", verifyCmd)

		output, err := command.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to verify DNS: %w\n%s", err, string(output))
		}

		fmt.Printf("DNS records for %s:\n%s", domain, string(output))
		return nil
	},
}

var updateHostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "Update hosts file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverIP == "" || domain == "" {
			return fmt.Errorf("server IP and domain are required")
		}

		// Update /etc/hosts
		updateCmd := fmt.Sprintf(`
			grep -v "%s" /etc/hosts > /tmp/hosts
			echo "%s %s" >> /tmp/hosts
			sudo mv /tmp/hosts /etc/hosts
		`, domain, serverIP, domain)

		command := exec.Command("bash", "-c", updateCmd)
		if output, err := command.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to update hosts file: %w\n%s", err, string(output))
		}

		fmt.Printf("Hosts file updated for %s\n", domain)
		return nil
	},
}

var (
	zoneID   string
	apiToken string
)

func init() {
	// Add flags
	dnsCmd.PersistentFlags().StringVar(&domain, "domain", "", "Domain name")
	updateDNSCmd.Flags().StringVar(&subdomain, "subdomain", "", "Subdomain name")
	updateDNSCmd.Flags().StringVar(&zoneID, "zone-id", "", "Cloudflare Zone ID")
	updateDNSCmd.Flags().StringVar(&apiToken, "api-token", "", "Cloudflare API Token")
	updateHostsCmd.Flags().StringVar(&serverIP, "server-ip", "", "Server IP")

	// Add subcommands
	dnsCmd.AddCommand(updateDNSCmd, deleteDNSCmd, listDNSCmd, verifyDNSCmd, updateHostsCmd)

	// Add to root command
	rootCmd.AddCommand(dnsCmd)
}
