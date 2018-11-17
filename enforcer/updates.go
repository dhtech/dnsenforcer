package enforcer

import "log"

// UpdateRecords logs all records to stdout
func (e *Enforcer) UpdateRecords() error {
	records, err := e.GetAllRecords()
	if err != nil {
		log.Fatal(err)
	}
	if e.Vars.DryRun {
		for _, record := range records {
			log.Printf("Dry-run: %s,%s,%d,%v", record.Name, record.Type, record.TTL, record.Data)
		}
		return nil
	}
	// TODO: Implement
	return nil
}
