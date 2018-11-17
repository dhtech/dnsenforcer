package enforcer

import (
	"github.com/dhtech/dnsenforcer/enforcer/ipplan"
	log "github.com/sirupsen/logrus"
)

// Enforcer is used to update DNS servers with new data
type Enforcer struct {
	IPPlan *ipplan.IPPlan
	Static string
}

// New returns a new DNS Enforcer
func New(dbfile, staticfile string) (*Enforcer, error) {
	p, err := ipplan.Open(dbfile)
	if err != nil {
		return nil, err
	}
	return &Enforcer{
		IPPlan: p,
		Static: staticfile,
	}, nil
}

// Close finalizes and releases resources held by the enforcer
func (e *Enforcer) Close() {
	e.IPPlan.Close()
}

// PrintAllRecords logs all records to stdout
func (e *Enforcer) PrintAllRecords() {
	records, err := e.GetAllRecords()
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range records {
		log.WithFields(log.Fields{
			"name": record.Name,
			"type": record.Type,
			"data": record.Data,
		}).Println()
	}
}
