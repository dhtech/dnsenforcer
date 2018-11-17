package enforcer

// Record holds a DNS record
type Record struct {
	Name string
	Data []string
	Type string
	TTL  int
}

// GetAllRecords will return the full DNS record dataset
func (e *Enforcer) GetAllRecords() ([]*Record, error) {
	return e.GetHostRecords()
}
