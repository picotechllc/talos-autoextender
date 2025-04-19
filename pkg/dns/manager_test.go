package dns

import (
	"testing"
)

func TestNewDNSManager(t *testing.T) {
	manager := NewDNSManager("cloudflare", "example.com")

	if manager == nil {
		t.Fatal("Expected non-nil DNSManager")
	}

	if manager.provider != "cloudflare" {
		t.Errorf("Expected provider to be 'cloudflare', got '%s'", manager.provider)
	}

	if manager.domain != "example.com" {
		t.Errorf("Expected domain to be 'example.com', got '%s'", manager.domain)
	}

	if manager.records == nil {
		t.Error("Expected non-nil records slice")
	}

	if len(manager.records) != 0 {
		t.Errorf("Expected empty records, got %d items", len(manager.records))
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		provider    string
		domain      string
		shouldError bool
	}{
		{
			name:        "valid config",
			provider:    "cloudflare",
			domain:      "example.com",
			shouldError: false,
		},
		{
			name:        "missing provider",
			provider:    "",
			domain:      "example.com",
			shouldError: true,
		},
		{
			name:        "missing domain",
			provider:    "cloudflare",
			domain:      "",
			shouldError: true,
		},
		{
			name:        "missing both",
			provider:    "",
			domain:      "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewDNSManager(tt.provider, tt.domain)
			err := manager.ValidateConfig()
			if (err != nil) != tt.shouldError {
				t.Errorf("ValidateConfig() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}

func TestUpsertRecord(t *testing.T) {
	manager := NewDNSManager("cloudflare", "example.com")

	tests := []struct {
		name        string
		recordName  string
		recordType  string
		content     string
		shouldError bool
	}{
		{
			name:        "valid A record",
			recordName:  "www",
			recordType:  "A",
			content:     "192.168.1.1",
			shouldError: false,
		},
		{
			name:        "valid CNAME record",
			recordName:  "alias",
			recordType:  "CNAME",
			content:     "www.example.com",
			shouldError: false,
		},
		{
			name:        "missing name",
			recordName:  "",
			recordType:  "A",
			content:     "192.168.1.1",
			shouldError: true,
		},
		{
			name:        "missing type",
			recordName:  "www",
			recordType:  "",
			content:     "192.168.1.1",
			shouldError: true,
		},
		{
			name:        "missing content",
			recordName:  "www",
			recordType:  "A",
			content:     "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.UpsertRecord(tt.recordName, tt.recordType, tt.content)
			if (err != nil) != tt.shouldError {
				t.Errorf("UpsertRecord() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}

	// Check that valid records were added
	expectedCount := 2 // Only the first two test cases should succeed
	records, err := manager.ListRecords()
	if err != nil {
		t.Fatalf("ListRecords() error = %v", err)
	}

	if len(records) != expectedCount {
		t.Errorf("Expected %d records, got %d", expectedCount, len(records))
	}
}

func TestListRecords(t *testing.T) {
	manager := NewDNSManager("cloudflare", "example.com")

	// Add some records
	if err := manager.UpsertRecord("www", "A", "192.168.1.1"); err != nil {
		t.Fatalf("Failed to add record: %v", err)
	}

	if err := manager.UpsertRecord("mail", "MX", "mail.example.com"); err != nil {
		t.Fatalf("Failed to add record: %v", err)
	}

	// List records
	records, err := manager.ListRecords()
	if err != nil {
		t.Fatalf("ListRecords() error = %v", err)
	}

	if len(records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(records))
	}

	// Verify record content
	found := make(map[string]bool)
	for _, record := range records {
		if record.Name == "www" && record.Type == "A" && record.Content == "192.168.1.1" {
			found["www"] = true
		}
		if record.Name == "mail" && record.Type == "MX" && record.Content == "mail.example.com" {
			found["mail"] = true
		}
	}

	if !found["www"] {
		t.Error("www record not found or has incorrect values")
	}

	if !found["mail"] {
		t.Error("mail record not found or has incorrect values")
	}
}

// Test for future implementation of DNS provider integration
func TestCloudflareIntegration(t *testing.T) {
	t.Skip("Cloudflare integration not yet implemented")
}

// Test for future implementation of weighted records
func TestWeightedRecords(t *testing.T) {
	t.Skip("Weighted records not yet implemented")
}

// Test for future implementation of health checks
func TestDNSHealthChecks(t *testing.T) {
	t.Skip("DNS health checks not yet implemented")
}

// Test for future implementation of blue/green deployments
func TestBlueGreenDeployment(t *testing.T) {
	t.Skip("Blue/green deployment not yet implemented")
}
