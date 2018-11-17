package enforcer

import (
	"io"
	"os"

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
		records = append(records, record)
	}
}
