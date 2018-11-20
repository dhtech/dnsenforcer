package enforcer

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// GetStaticRecords returns records that are specified in static YAML file
func (e *Enforcer) GetStaticRecords() ([]*Record, error) {
	data, err := os.Open(e.Vars.Static)
	if err != nil {
		return nil, err
	}
	var records []*Record
	reader := yaml.NewDecoder(data)
	for {
		var record *Record
		err := reader.Decode(&record)
		if err == io.EOF {
			return records, nil
		} else if err != nil {
			return nil, err
		}
		if record != nil {
			if !e.inZones(record.Name) {
				log.Warningf("%s is not a member of any of the enforced zones", record.Name)
				continue
			}
			if e.ignoredType(record.Type) {
				log.Warningf("Found ignored type %s for %s in static file", record.Type, record.Name)
				continue
			}
			for _, d := range record.Data {
				records = append(records, &Record{
					Name: record.Name,
					TTL:  record.TTL,
					Type: record.Type,
					Data: []string{d},
				})
			}
		}
	}
}
