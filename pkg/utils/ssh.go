package utils

import (
	"fmt"
	"os/exec"
)

// ExecuteRemoteCommand executes a command on a remote server via SSH
func ExecuteRemoteCommand(ip, user, pass, command string) (string, error) {
	sshCmd := exec.Command("sshpass", "-p", pass, "ssh", "-o", "StrictHostKeyChecking=no",
		fmt.Sprintf("%s@%s", user, ip), command)

	output, err := sshCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %w\nOutput: %s", err, string(output))
	}

	return string(output), nil
}
