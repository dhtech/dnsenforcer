package enforcer

import (
	"fmt"
	"io"
	"os"
	"strings"

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
		if record != nil && !e.ignoredType(record.Type) { // TODO (rctl): Make ignored types dynamic
			if !e.inZones(record.Name) {
				log.Warningf("%s is not a member of any of the enforced zones", record.Name)
				continue
			}
			for _, d := range record.Data { // TODO (rctl): Make it so that records does not have to be single data entry
				if record.Type == "CNAME" && !strings.HasSuffix(d, ".") {
					d = fmt.Sprintf("%s.", d)
				}
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
