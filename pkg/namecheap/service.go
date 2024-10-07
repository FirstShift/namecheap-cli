package namecheap

import (
	"errors"
	"fmt"

	nc "github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

type NamecheapService struct {
	Client *nc.Client
}

type NamecheapDNSRecord struct {
	Host  string `json:"host"`
	Type  string `json:"record_type"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}

func New(client *nc.Client) *NamecheapService {
	return &NamecheapService{
		Client: client,
	}
}

func GetCurrentRecords(client *nc.Client, domain string) (*[]nc.DomainsDNSHostRecordDetailed, error) {
	// Get the current records for the domain
	response, err := client.DomainsDNS.GetHosts(domain)
	if err != nil {
		return nil, err
	}

	if response.DomainDNSGetHostsResult == nil {
		return nil, nil
	}

	return response.DomainDNSGetHostsResult.Hosts, nil
}

func FindAndRemoveRecord(currentRecords *[]nc.DomainsDNSHostRecordDetailed, record *NamecheapDNSRecord) (*[]nc.DomainsDNSHostRecordDetailed, bool) {
	updatedRecords := make([]nc.DomainsDNSHostRecordDetailed, 0, len(*currentRecords))
	found := false

	for _, currentRecord := range *currentRecords {
		// Check if the record matches the one to delete
		if *currentRecord.Name == record.Host && *currentRecord.Type == record.Type {
			found = true
		} else {
			updatedRecords = append(updatedRecords, currentRecord)
		}
	}

	return &updatedRecords, found
}

func PrintRecords(records *[]nc.DomainsDNSHostRecordDetailed) {
	for _, record := range *records {
		if record.Name != nil && record.Type != nil && record.Address != nil && record.TTL != nil {
			fmt.Printf("Name: %s, Type: %s, Address: %s, TTL: %d\n", *record.Name, *record.Type, *record.Address, *record.TTL)
		}
	}
}

func (svc *NamecheapService) CreateRecord(domain string, record *NamecheapDNSRecord, dryRun bool) error {
	// Get the current records for the domain
	currentRecords, err := GetCurrentRecords(svc.Client, domain)
	if err != nil {
		return err
	}

	// Prepare updatedRecords as a slice of values instead of pointers
	updatedRecords := make([]nc.DomainsDNSHostRecordDetailed, 0, len(*currentRecords)+1)

	for _, currentRecord := range *currentRecords {
		if currentRecord.Name != nil && currentRecord.Type != nil && currentRecord.Address != nil {
			if *currentRecord.Name == record.Host && *currentRecord.Type == record.Type && *currentRecord.Address == record.Value {
				return errors.New("record already exists")
			}
		}
		updatedRecords = append(updatedRecords, currentRecord)
	}

	// Add the new record as a value (not pointer)
	updatedRecords = append(updatedRecords, nc.DomainsDNSHostRecordDetailed{
		Name:    &record.Host,
		Type:    &record.Type,
		Address: &record.Value,
		TTL:     &record.TTL,
	})

	// If dryRun is enabled, skip the actual API call
	if dryRun {
		// Print the updated records
		fmt.Println("Performing dry run. The following records would be updated:")
		for _, record := range updatedRecords {
			if record.Name != nil && record.Type != nil && record.Address != nil && record.TTL != nil {
				fmt.Printf("Name: %s, Type: %s, Address: %s, TTL: %d\n", *record.Name, *record.Type, *record.Address, *record.TTL)
			}
		}
		return nil
	}

	// Convert DomainsDNSHostRecordDetailed to DomainsDNSHostRecord
	hostRecords := make([]nc.DomainsDNSHostRecord, len(updatedRecords))
	for i, record := range updatedRecords {
		hostRecords[i] = nc.DomainsDNSHostRecord{
			HostName:   record.Name,
			RecordType: record.Type,
			Address:    record.Address,
			TTL:        record.TTL,
		}
	}

	// Now use hostRecords in the SetHosts call
	_, err = svc.Client.DomainsDNS.SetHosts(&nc.DomainsDNSSetHostsArgs{
		Domain:  nc.String(domain),
		Records: &hostRecords,
	})

	if err != nil {
		return err
	}

	return nil
}

func (svc *NamecheapService) DeleteRecord(domain string, record *NamecheapDNSRecord, dryRun bool) error {
	fmt.Println("Deleting record...")
	// Get the current records for the domain
	currentRecords, err := GetCurrentRecords(svc.Client, domain)
	if err != nil {
		return fmt.Errorf("failed to retrieve current records for domain '%s': %w", domain, err)
	}

	// Remove the record from current records
	updatedRecords, found := FindAndRemoveRecord(currentRecords, record)
	if !found {
		return fmt.Errorf("record not found: %s %s", record.Host, record.Type)
	}

	// Dry run - Print the remaining records and exit early
	if dryRun {
		fmt.Println("Performing dry run. The following records would remain after deletion:")
		PrintRecords(updatedRecords)
		return nil
	}

	// Update the DNS records
	hostRecords := make([]nc.DomainsDNSHostRecord, len(*updatedRecords))
	for i, record := range *updatedRecords {
		hostRecords[i] = nc.DomainsDNSHostRecord{
			HostName:   record.Name,
			RecordType: record.Type,
			Address:    record.Address,
			TTL:        record.TTL,
		}
	}

	_, err = svc.Client.DomainsDNS.SetHosts(&nc.DomainsDNSSetHostsArgs{
		Domain:  nc.String(domain),
		Records: &hostRecords,
	})
	if err != nil {
		return fmt.Errorf("failed to update DNS records: %w", err)
	}

	fmt.Println("Record successfully deleted.")
	return nil
}

func (svc *NamecheapService) ListRecords(domain string) ([]*NamecheapDNSRecord, error) {
	records := make([]*NamecheapDNSRecord, 0)
	response, err := svc.Client.DomainsDNS.GetHosts(domain)
	if err != nil {
		return nil, err
	}

	if response.DomainDNSGetHostsResult == nil || response.DomainDNSGetHostsResult.Hosts == nil {
		return records, nil
	}

	for _, host := range *response.DomainDNSGetHostsResult.Hosts {
		record := &NamecheapDNSRecord{
			Host:  *host.Name,
			Type:  *host.Type,
			Value: *host.Address,
			TTL:   *host.TTL,
		}
		records = append(records, record)
	}

	return records, nil
}
