package migration

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cploutarchou/swarmforge/pkg/utils"
)

type TraefikConfig struct {
	Certificates map[string]interface{} `json:"certificates"`
	HTTPRouters  map[string]interface{} `json:"http"`
	TCPRouters   map[string]interface{} `json:"tcp"`
	Middlewares  map[string]interface{} `json:"middlewares"`
}

func MigrateTraefik(sourceIP, targetIP, username, password string) error {
	// Check if Traefik is running on source
	output, err := utils.ExecuteRemoteCommand(sourceIP, username, password,
		"docker service ls --filter name=traefik --format '{{.Name}}'")
	if err != nil {
		return fmt.Errorf("failed to check Traefik service: %w", err)
	}
	if !strings.Contains(output, "traefik") {
		return fmt.Errorf("Traefik service not found on source node")
	}

	// Backup current Traefik configuration
	backupTime := time.Now().Format("20060102150405")
	backupDir := fmt.Sprintf("/tmp/traefik_backup_%s", backupTime)

	backupCmd := fmt.Sprintf(`
		mkdir -p %s &&
		docker service inspect traefik > %s/traefik_service.json &&
		cp -r /etc/traefik %s/config &&
		cp /var/lib/docker/volumes/traefik-certs/_data/acme.json %s/certs/
	`, backupDir, backupDir, backupDir, backupDir)
	if _, err := utils.ExecuteRemoteCommand(sourceIP, username, password, backupCmd); err != nil {
		return fmt.Errorf("failed to backup Traefik: %w", err)
	}

	// Get current Traefik configuration
	cmd := "docker service inspect traefik"
	output, err = utils.ExecuteRemoteCommand(sourceIP, username, password, cmd)
	if err != nil {
		return fmt.Errorf("failed to inspect Traefik service: %w", err)
	}

	var traefikService []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &traefikService); err != nil {
		return fmt.Errorf("failed to parse Traefik service: %w", err)
	}

	// Extract important configurations
	var config TraefikConfig
	if err := extractTraefikConfig(traefikService[0], &config); err != nil {
		return fmt.Errorf("failed to extract Traefik config: %w", err)
	}

	// Prepare target node
	setupCmds := []string{
		"mkdir -p /etc/traefik/config",
		"mkdir -p /etc/traefik/certs",
	}

	for _, cmd := range setupCmds {
		if _, err := utils.ExecuteRemoteCommand(targetIP, username, password, cmd); err != nil {
			return fmt.Errorf("failed to setup target node: %w", err)
		}
	}

	// Copy configurations to target node
	copyCmds := []string{
		fmt.Sprintf("scp -r %s/config/* %s@%s:/etc/traefik/config/", backupDir, username, targetIP),
		fmt.Sprintf("scp %s/certs/acme.json %s@%s:/etc/traefik/certs/", backupDir, username, targetIP),
	}

	for _, cmd := range copyCmds {
		if _, err := utils.ExecuteRemoteCommand(sourceIP, username, password, cmd); err != nil {
			return fmt.Errorf("failed to copy configurations: %w", err)
		}
	}

	// Stop Traefik on source
	if err := stopTraefik(sourceIP, username, password); err != nil {
		return fmt.Errorf("failed to stop Traefik on source: %w", err)
	}

	// Start Traefik on target
	if err := startTraefik(targetIP, username, password); err != nil {
		// Rollback if failed
		if rollbackErr := startTraefik(sourceIP, username, password); rollbackErr != nil {
			return fmt.Errorf("failed to start Traefik on target and rollback failed: %v, rollback error: %v", err, rollbackErr)
		}
		return fmt.Errorf("failed to start Traefik on target: %w", err)
	}

	// Verify Traefik is running on target
	if err := verifyTraefik(targetIP, username, password); err != nil {
		return fmt.Errorf("Traefik verification failed on target: %w", err)
	}

	return nil
}

func extractTraefikConfig(service map[string]interface{}, config *TraefikConfig) error {
	// Implementation to extract configuration from service inspection
	return nil
}

func stopTraefik(ip, username, password string) error {
	cmd := "docker service rm traefik"
	_, err := utils.ExecuteRemoteCommand(ip, username, password, cmd)
	return err
}

func startTraefik(ip, username, password string) error {
	cmd := `docker service create \
		--name traefik \
		--publish 80:80 \
		--publish 443:443 \
		--mount type=bind,source=/etc/traefik,target=/etc/traefik \
		--mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock \
		--network traefik-public \
		traefik:v2.10`

	_, err := utils.ExecuteRemoteCommand(ip, username, password, cmd)
	return err
}

func verifyTraefik(ip, username, password string) error {
	// Wait for service to be running
	maxAttempts := 30
	for i := 0; i < maxAttempts; i++ {
		cmd := "docker service ls --filter name=traefik --format '{{.Name}}\t{{.Replicas}}'"
		output, err := utils.ExecuteRemoteCommand(ip, username, password, cmd)
		if err == nil && strings.Contains(output, "1/1") {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("Traefik service failed to start properly")
}
