package enforcer

import (
	"io"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// GetStaticRecords returns records that are specified in static YAML file
func (e *Enforcer) GetStaticRecords(hostr []*Record) ([]*Record, error) {
	m := make(map[string][]*Record)
	for _, r := range hostr {
		if d, e := m[r.Name]; e {
			m[r.Name] = append(d, r)
		} else {
			m[r.Name] = []*Record{r}
		}
	}
	var records []*Record
	reader := yaml.NewDecoder(e.static)
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
				if record.Type == "ALIAS" {
					if dsts, e := m[d]; e {
						for _, dst := range dsts {
							for _, ds := range dst.Data {
								records = append(records, &Record{
									Name: record.Name,
									TTL:  record.TTL,
									Type: dst.Type,
									Data: []string{ds},
								})
							}
						}
						log.Infof("Added ALIAS %s for %s", record.Name, d)
					} else {
						log.Warningf("Found ALIAS %s pointing to non existing host %s", record.Name, d)
					}
				} else {
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
}
