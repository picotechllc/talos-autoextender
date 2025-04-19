
package dns

import (
	"testing"
	"time"
)

func TestDNSFailoverPatterns(t *testing.T) {
	tests := []struct {
		name         string
		provider     string
		domain       string
		records      []Record
		failoverSpec FailoverSpec
		shouldError  bool
	}{
		{
			name:     "gradual traffic shift",
			provider: "cloudflare",
			domain:   "example.com",
			records: []Record{
				{Name: "api", Type: "A", Content: "192.168.1.1", Weight: 100},
				{Name: "api", Type: "A", Content: "192.168.1.2", Weight: 0},
			},
			failoverSpec: FailoverSpec{
				Type:           "weighted",
				StepPercentage: 20,
				StepInterval:   time.Second * 30,
			},
			shouldError: false,
		},
		{
			name:     "instant failover",
			provider: "cloudflare",
			domain:   "example.com",
			records: []Record{
				{Name: "api", Type: "A", Content: "192.168.1.1"},
				{Name: "api", Type: "A", Content: "192.168.1.2"},
			},
			failoverSpec: FailoverSpec{
				Type: "immediate",
			},
			shouldError: false,
		},
		{
			name:     "geolocation-based failover",
			provider: "cloudflare",
			domain:   "example.com",
			records: []Record{
				{Name: "api", Type: "A", Content: "192.168.1.1", GeoLocation: "US"},
				{Name: "api", Type: "A", Content: "192.168.1.2", GeoLocation: "EU"},
			},
			failoverSpec: FailoverSpec{
				Type:      "geo",
				Regions:   []string{"US", "EU"},
				Priority: map[string]int{"US": 1, "EU": 2},
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewDNSManager(tt.provider, tt.domain)

			// Setup initial records
			for _, record := range tt.records {
				err := manager.UpsertRecord(record)
				if err != nil {
					t.Fatalf("Failed to create initial record: %v", err)
				}
			}

			// Execute failover
			err := manager.ExecuteFailover(tt.failoverSpec)
			if (err != nil) != tt.shouldError {
				t.Errorf("ExecuteFailover() error = %v, shouldError %v", err, tt.shouldError)
			}

			// Verify record states
			records, err := manager.GetRecords("api")
			if err != nil {
				t.Fatalf("Failed to get records: %v", err)
			}

			switch tt.failoverSpec.Type {
			case "weighted":
				verifyWeightedFailover(t, records, tt.failoverSpec)
			case "immediate":
				verifyImmediateFailover(t, records)
			case "geo":
				verifyGeoFailover(t, records, tt.failoverSpec)
			}
		})
	}
}

func verifyWeightedFailover(t *testing.T, records []Record, spec FailoverSpec) {
	var totalWeight int
	for _, record := range records {
		totalWeight += record.Weight
	}
	if totalWeight != 100 {
		t.Errorf("Total weight should be 100, got %d", totalWeight)
	}
}

func verifyImmediateFailover(t *testing.T, records []Record) {
	activeCount := 0
	for _, record := range records {
		if record.Active {
			activeCount++
		}
	}
	if activeCount != 1 {
		t.Errorf("Expected exactly one active record, got %d", activeCount)
	}
}

func verifyGeoFailover(t *testing.T, records []Record, spec FailoverSpec) {
	for _, record := range records {
		if record.GeoLocation == "" {
			t.Errorf("Record missing geolocation")
		}
		if _, exists := spec.Priority[record.GeoLocation]; !exists {
			t.Errorf("Record has invalid geolocation: %s", record.GeoLocation)
		}
	}
}

func TestDNSHealthCheck(t *testing.T) {
	tests := []struct {
		name       string
		endpoints  []string
		threshold  int
		interval   time.Duration
		shouldFail bool
	}{
		{
			name:       "detect unhealthy endpoint",
			endpoints:  []string{"https://api1.example.com", "https://api2.example.com"},
			threshold:  3,
			interval:   time.Second * 5,
			shouldFail: true,
		},
		{
			name:       "maintain healthy endpoints",
			endpoints:  []string{"https://healthy1.example.com", "https://healthy2.example.com"},
			threshold:  3,
			interval:   time.Second * 5,
			shouldFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monitor := NewHealthMonitor(tt.threshold, tt.interval)
			
			for _, endpoint := range tt.endpoints {
				err := monitor.AddEndpoint(endpoint)
				if err != nil {
					t.Fatalf("Failed to add endpoint: %v", err)
				}
			}

			status := monitor.CheckHealth()
			if status.HasFailures != tt.shouldFail {
				t.Errorf("Health check result %v != expected %v", status.HasFailures, tt.shouldFail)
			}
		})
	}
}
