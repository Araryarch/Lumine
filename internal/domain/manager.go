package domain

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Manager struct {
	hostsFile string
}

func NewManager() *Manager {
	hostsFile := "/etc/hosts"
	if runtime.GOOS == "windows" {
		hostsFile = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	}

	return &Manager{
		hostsFile: hostsFile,
	}
}

// AddDomain adds a .test domain to /etc/hosts
func (m *Manager) AddDomain(domain string) error {
	if !strings.HasSuffix(domain, ".test") {
		domain = domain + ".test"
	}

	entry := fmt.Sprintf("127.0.0.1 %s\n", domain)

	// Check if domain already exists
	content, err := os.ReadFile(m.hostsFile)
	if err != nil {
		return err
	}

	if strings.Contains(string(content), domain) {
		return nil // Already exists
	}

	// Append to hosts file (requires sudo)
	cmd := exec.Command("sudo", "sh", "-c", fmt.Sprintf("echo '%s' >> %s", entry, m.hostsFile))
	return cmd.Run()
}

// RemoveDomain removes a .test domain from /etc/hosts
func (m *Manager) RemoveDomain(domain string) error {
	if !strings.HasSuffix(domain, ".test") {
		domain = domain + ".test"
	}

	cmd := exec.Command("sudo", "sed", "-i", fmt.Sprintf("/%s/d", domain), m.hostsFile)
	return cmd.Run()
}

// ListDomains lists all .test domains
func (m *Manager) ListDomains() ([]string, error) {
	content, err := os.ReadFile(m.hostsFile)
	if err != nil {
		return nil, err
	}

	var domains []string
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.Contains(line, ".test") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				domains = append(domains, parts[1])
			}
		}
	}

	return domains, nil
}
