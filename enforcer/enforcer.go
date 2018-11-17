package enforcer

import (
	"github.com/dhtech/dnsenforcer/enforcer/ipplan"
)

// Enforcer is used to update DNS servers with new data
type Enforcer struct {
	IPPlan *ipplan.IPPlan
	Vars   *Vars
}

// Vars hold values needed for enforcer
type Vars struct {
	DBFile  string
	Static  string
	HostTTL int
	DryRun  bool
}

// New returns a new DNS Enforcer
func New(vars *Vars) (*Enforcer, error) {
	p, err := ipplan.Open(vars.DBFile)
	if err != nil {
		return nil, err
	}
	return &Enforcer{
		IPPlan: p,
		Vars:   vars,
	}, nil
}

// Close finalizes and releases resources held by the enforcer
func (e *Enforcer) Close() {
	e.IPPlan.Close()
}
