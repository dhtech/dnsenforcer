package enforcer

// Record holds a DNS record
type Record struct {
	Name    string   `yaml:"name"`
	Data    []string `yaml:"rrdatas"`
	Type    string   `yaml:"type"`
	TTL     int      `yaml:"ttl"`
	Comment string   `yaml:"comment"`
	Owner   string   `yaml:"owner"`
}

// GetAllRecords will return the full DNS record dataset
func (e *Enforcer) GetAllRecords() ([]*Record, error) {
	hosts, err := e.GetHostRecords()
	if err != nil {
		return nil, err
	}
	static, err := e.GetStaticRecords()
	if err != nil {
		return nil, err
	}
	return append(hosts, static...), nil
}
