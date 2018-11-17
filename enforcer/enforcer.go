package enforcer

import (
	"github.com/dhtech/dnsenforcer/enforcer/ipplan"
	log "github.com/sirupsen/logrus"
)

// Enforcer is used to update DNS servers with new data
type Enforcer struct {
	IPPlan *ipplan.IPPlan
}

// New returns a new DNS Enforcer
func New(dbfile string) (*Enforcer, error) {
	p, err := ipplan.Open(dbfile)
	defer p.Close()
	if err != nil {
		return nil, err
	}
	hosts, err := p.Hosts()
	if err != nil {
		return nil, err
	}
	for _, host := range hosts {
		log.Infof("%s -> %s, %s, %v", host.Name, host.IPv4, host.IPv6, host.Options)
	}
	return &Enforcer{
		IPPlan: p,
	}, nil
}
