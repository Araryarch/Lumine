package port

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

// Manager handles port allocation and detection
type Manager struct {
	usedPorts map[int]bool
}

// NewManager creates a new port manager
func NewManager() *Manager {
	return &Manager{
		usedPorts: make(map[int]bool),
	}
}

// IsPortAvailable checks if a port is available
func (m *Manager) IsPortAvailable(port int) bool {
	// Check if already marked as used
	if m.usedPorts[port] {
		return false
	}

	// Try to listen on the port
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

// FindAvailablePort finds an available port starting from the given port
func (m *Manager) FindAvailablePort(startPort int, maxAttempts int) (int, error) {
	for i := 0; i < maxAttempts; i++ {
		port := startPort + i
		if m.IsPortAvailable(port) {
			m.usedPorts[port] = true
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port found after %d attempts starting from %d", maxAttempts, startPort)
}

// GetAlternativePort returns an alternative port if the original is taken
func (m *Manager) GetAlternativePort(originalPort int) (int, error) {
	// If original port is available, use it
	if m.IsPortAvailable(originalPort) {
		m.usedPorts[originalPort] = true
		return originalPort, nil
	}

	// Try common alternative ports based on service type
	alternatives := m.getAlternativePorts(originalPort)
	for _, altPort := range alternatives {
		if m.IsPortAvailable(altPort) {
			m.usedPorts[altPort] = true
			return altPort, nil
		}
	}

	// If no alternatives work, find any available port
	return m.FindAvailablePort(originalPort+100, 100)
}

// getAlternativePorts returns common alternative ports for a given port
func (m *Manager) getAlternativePorts(port int) []int {
	alternatives := make(map[int][]int)
	
	// MySQL alternatives
	alternatives[3306] = []int{3307, 3308, 3309, 33060, 33061}
	
	// PostgreSQL alternatives
	alternatives[5432] = []int{5433, 5434, 5435, 54320, 54321}
	
	// MongoDB alternatives
	alternatives[27017] = []int{27018, 27019, 27020, 27021}
	
	// Redis alternatives
	alternatives[6379] = []int{6380, 6381, 6382, 6383}
	
	// HTTP alternatives
	alternatives[80] = []int{8080, 8000, 8888, 3000, 5000}
	alternatives[8080] = []int{8081, 8082, 8083, 8084, 8085}
	
	// HTTPS alternatives
	alternatives[443] = []int{8443, 4443, 9443}
	
	// Elasticsearch alternatives
	alternatives[9200] = []int{9201, 9202, 9203}
	
	if alts, ok := alternatives[port]; ok {
		return alts
	}
	
	// Default: try next 5 ports
	return []int{port + 1, port + 2, port + 3, port + 4, port + 5}
}

// ReservePort marks a port as used
func (m *Manager) ReservePort(port int) {
	m.usedPorts[port] = true
}

// ReleasePort marks a port as available
func (m *Manager) ReleasePort(port int) {
	delete(m.usedPorts, port)
}

// GetUsedPorts returns all currently used ports
func (m *Manager) GetUsedPorts() []int {
	ports := make([]int, 0, len(m.usedPorts))
	for port := range m.usedPorts {
		ports = append(ports, port)
	}
	return ports
}

// WaitForPort waits for a port to become available
func (m *Manager) WaitForPort(port int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	
	for time.Now().Before(deadline) {
		if m.IsPortAvailable(port) {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	
	return fmt.Errorf("timeout waiting for port %d to become available", port)
}

// GetPortInfo returns information about what's using a port
func GetPortInfo(port int) (string, error) {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Sprintf("Port %d is in use", port), nil
	}
	listener.Close()
	return fmt.Sprintf("Port %d is available", port), nil
}

// ScanPortRange scans a range of ports and returns available ones
func (m *Manager) ScanPortRange(startPort, endPort int) []int {
	available := []int{}
	for port := startPort; port <= endPort; port++ {
		if m.IsPortAvailable(port) {
			available = append(available, port)
		}
	}
	return available
}

// GetServicePort returns the appropriate port for a service type
func (m *Manager) GetServicePort(serviceType string, preferredPort int) (int, error) {
	// Default ports for common services
	defaultPorts := map[string]int{
		"mysql":         3306,
		"mariadb":       3307,
		"postgres":      5432,
		"postgresql":    5432,
		"mongodb":       27017,
		"mongo":         27017,
		"redis":         6379,
		"nginx":         80,
		"apache":        80,
		"httpd":         80,
		"elasticsearch": 9200,
		"rabbitmq":      5672,
		"memcached":     11211,
		"phpmyadmin":    8080,
		"adminer":       8081,
	}

	// Use preferred port if specified
	if preferredPort > 0 {
		return m.GetAlternativePort(preferredPort)
	}

	// Use default port for service type
	if defaultPort, ok := defaultPorts[serviceType]; ok {
		return m.GetAlternativePort(defaultPort)
	}

	// Find any available port in common range
	return m.FindAvailablePort(8000, 100)
}

// ValidatePortRange checks if a port is in valid range
func ValidatePortRange(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("port %d is out of valid range (1-65535)", port)
	}
	if port < 1024 {
		return fmt.Errorf("port %d is in privileged range (requires root)", port)
	}
	return nil
}

// FormatPortMapping returns a formatted port mapping string
func FormatPortMapping(hostPort, containerPort int) string {
	return fmt.Sprintf("%d:%d", hostPort, containerPort)
}

// ParsePortMapping parses a port mapping string
func ParsePortMapping(mapping string) (hostPort, containerPort int, err error) {
	var host, container string
	_, err = fmt.Sscanf(mapping, "%s:%s", &host, &container)
	if err != nil {
		return 0, 0, err
	}
	
	hostPort, err = strconv.Atoi(host)
	if err != nil {
		return 0, 0, err
	}
	
	containerPort, err = strconv.Atoi(container)
	if err != nil {
		return 0, 0, err
	}
	
	return hostPort, containerPort, nil
}
