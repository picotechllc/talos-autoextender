
package dns

import (
	"testing"
)

func TestDNSManagerConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		provider    string
		domain      string
		shouldError bool
	}{
		{
			name:        "valid cloudflare config",
			provider:    "cloudflare",
			domain:     "example.com",
			shouldError: false,
		},
		{
			name:        "missing provider",
			provider:    "",
			domain:     "example.com",
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

func TestDNSRecordManagement(t *testing.T) {
	manager := NewDNSManager("cloudflare", "example.com")

	err := manager.UpsertRecord("test", "A", "192.168.1.1")
	if err != nil {
		t.Errorf("UpsertRecord() error = %v", err)
	}

	records, err := manager.ListRecords()
	if err != nil {
		t.Errorf("ListRecords() error = %v", err)
	}

	if len(records) != 1 {
		t.Errorf("Expected 1 record, got %d", len(records))
	}
}
