
package dns

import (
	"fmt"
)

type DNSManager struct {
	provider string
	domain   string
	records  []Record
}

type Record struct {
	Name    string
	Type    string
	Content string
}

func NewDNSManager(provider, domain string) *DNSManager {
	return &DNSManager{
		provider: provider,
		domain:   domain,
		records:  make([]Record, 0),
	}
}

func (d *DNSManager) ValidateConfig() error {
	if d.provider == "" {
		return fmt.Errorf("DNS provider is required")
	}
	if d.domain == "" {
		return fmt.Errorf("domain is required")
	}
	return nil
}

func (d *DNSManager) UpsertRecord(name, recordType, content string) error {
	if name == "" || recordType == "" || content == "" {
		return fmt.Errorf("name, type, and content are required")
	}

	record := Record{
		Name:    name,
		Type:    recordType,
		Content: content,
	}

	d.records = append(d.records, record)
	return nil
}

func (d *DNSManager) ListRecords() ([]Record, error) {
	return d.records, nil
}
