package enforcer

// Record holds a DNS record
type Record struct {
	Name string
	Data []string
	Type string
	TTL  int
}
