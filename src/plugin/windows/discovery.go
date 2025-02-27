package windows

import (
	"bytes"
	"fmt"
	"github.com/masterzen/winrm"
	"log"
	"time"
)

func Discover(ip string, username string, password string) string {
	// Validate input
	if ip == "" || username == "" || password == "" {
		return `{"error": "Missing required fields: IP, Username, Password"}`
	}

	// Create a WinRM client with a timeout
	endpoint := winrm.NewEndpoint(ip, 5985, false, false, nil, nil, nil, 30*time.Second) // 30s timeout
	client, err := winrm.NewClient(endpoint, username, password)
	if err != nil {
		log.Printf("Failed to create WinRM client: %v", err)
		return fmt.Sprintf(`{"error": "Failed to create WinRM client: %v"}`, err)
	}

	// Execute a command (hostname)
	var stdout, stderr bytes.Buffer
	exitCode, err := client.Run("hostname", &stdout, &stderr)
	if err != nil || exitCode != 0 {
		return fmt.Sprintf(`{"error": "Failed to execute command: %v", "stderr": "%s"}`, err, stderr.String())
	}

	// Return success response
	return fmt.Sprintf(`{"message": "Windows machine discovered successfully", "output": "%s"}`, stdout.String())
}
