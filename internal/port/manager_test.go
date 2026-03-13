package port

import (
	"testing"
	"time"
)

func TestIsPortAvailable(t *testing.T) {
	mgr := NewManager()
	
	// Port 0 should always be available (OS assigns)
	if !mgr.IsPortAvailable(0) {
		t.Error("Port 0 should be available")
	}
	
	// Very high port should be available
	if !mgr.IsPortAvailable(65000) {
		t.Error("Port 65000 should be available")
	}
}

func TestFindAvailablePort(t *testing.T) {
	mgr := NewManager()
	
	port, err := mgr.FindAvailablePort(8000, 10)
	if err != nil {
		t.Fatalf("Failed to find available port: %v", err)
	}
	
	if port < 8000 || port >= 8010 {
		t.Errorf("Port %d is out of expected range", port)
	}
}

func TestGetAlternativePort(t *testing.T) {
	mgr := NewManager()
	
	// Test MySQL alternatives
	port, err := mgr.GetAlternativePort(3306)
	if err != nil {
		t.Fatalf("Failed to get alternative port: %v", err)
	}
	
	if port != 3306 && port != 3307 && port != 3308 {
		t.Errorf("Unexpected alternative port: %d", port)
	}
}

func TestReserveAndReleasePort(t *testing.T) {
	mgr := NewManager()
	
	port := 9999
	mgr.ReservePort(port)
	
	if mgr.IsPortAvailable(port) {
		t.Error("Reserved port should not be available")
	}
	
	mgr.ReleasePort(port)
	
	if !mgr.IsPortAvailable(port) {
		t.Error("Released port should be available")
	}
}

func TestGetServicePort(t *testing.T) {
	mgr := NewManager()
	
	tests := []struct {
		serviceType string
		expected    int
	}{
		{"mysql", 3306},
		{"postgres", 5432},
		{"redis", 6379},
		{"mongodb", 27017},
	}
	
	for _, tt := range tests {
		port, err := mgr.GetServicePort(tt.serviceType, 0)
		if err != nil {
			t.Errorf("Failed to get port for %s: %v", tt.serviceType, err)
		}
		
		// Port should be the default or an alternative
		if port < tt.expected || port > tt.expected+10 {
			t.Errorf("Port %d for %s is out of expected range", port, tt.serviceType)
		}
	}
}

func TestValidatePortRange(t *testing.T) {
	tests := []struct {
		port      int
		shouldErr bool
	}{
		{0, true},      // Invalid
		{80, true},     // Privileged
		{1024, false},  // Valid
		{8080, false},  // Valid
		{65535, false}, // Valid
		{65536, true},  // Invalid
	}
	
	for _, tt := range tests {
		err := ValidatePortRange(tt.port)
		if (err != nil) != tt.shouldErr {
			t.Errorf("ValidatePortRange(%d) error = %v, shouldErr = %v", tt.port, err, tt.shouldErr)
		}
	}
}

func TestScanPortRange(t *testing.T) {
	mgr := NewManager()
	
	available := mgr.ScanPortRange(9000, 9010)
	
	if len(available) == 0 {
		t.Error("Should find at least some available ports")
	}
	
	for _, port := range available {
		if port < 9000 || port > 9010 {
			t.Errorf("Port %d is out of scan range", port)
		}
	}
}

func TestWaitForPort(t *testing.T) {
	mgr := NewManager()
	
	// Should succeed immediately for available port
	err := mgr.WaitForPort(9876, 1*time.Second)
	if err != nil {
		t.Errorf("WaitForPort failed: %v", err)
	}
}

func TestFormatPortMapping(t *testing.T) {
	mapping := FormatPortMapping(8080, 80)
	expected := "8080:80"
	
	if mapping != expected {
		t.Errorf("FormatPortMapping() = %s, want %s", mapping, expected)
	}
}

func TestParsePortMapping(t *testing.T) {
	host, container, err := ParsePortMapping("8080:80")
	if err != nil {
		t.Fatalf("ParsePortMapping failed: %v", err)
	}
	
	if host != 8080 {
		t.Errorf("Host port = %d, want 8080", host)
	}
	
	if container != 80 {
		t.Errorf("Container port = %d, want 80", container)
	}
}
