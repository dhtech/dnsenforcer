package enforcer

import (
	"io"

	"github.com/dhtech/dnsenforcer/enforcer/ipplan"
)

// Enforcer is used to update DNS servers with new data
type Enforcer struct {
	Vars   *Vars
	IPPlan *ipplan.IPPlan

	static io.Reader
}

// Vars hold values needed for enforcer
type Vars struct {
	Endpoint    string
	Certificate string
	Key         string
	Zones       []string
	HostTTL     int
	DryRun      bool
	IgnoreTypes []string
}

// New returns a new DNS Enforcer
func New(vars *Vars, ipp *ipplan.IPPlan, static io.Reader) (*Enforcer, error) {
	return &Enforcer{
		Vars:   vars,
		IPPlan: ipp,
		static: static,
	}, nil
}

// Close finalizes and releases resources held by the enforcer
func (e *Enforcer) Close() {
	e.IPPlan.Close()
}
