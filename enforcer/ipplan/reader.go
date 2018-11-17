package ipplan

import (
	"database/sql"
	"fmt"
	"net"

	_ "github.com/mattn/go-sqlite3" // Driver for SQLite3
	log "github.com/sirupsen/logrus"
)

// IPPlan is used to get structured data from ipplan database
type IPPlan struct {
	db *sql.DB
}

// Host is a network host
type Host struct {
	Name    string
	IPv4    net.IP
	IPv6    net.IP
	Options map[string][]string
}

// Open and return an IPPlan instance used to read data from ipplan database
func Open(dbfile string) (*IPPlan, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}
	return &IPPlan{
		db: db,
	}, nil
}

// Close the IPPlan database
func (p *IPPlan) Close() error {
	return p.db.Close()
}

// Hosts returns all network hosts
func (p *IPPlan) Hosts() ([]*Host, error) {
	// Fetch hosts
	rows, err := p.db.Query(`SELECT node_id,name,ipv4_addr_txt,ipv6_addr_txt FROM host;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	hosts := make([]*Host, 0)
	for rows.Next() {
		// Parse host data
		var id, name, ipv4, ipv6 string
		rows.Scan(&id, &name, &ipv4, &ipv6)

		// Extract host options
		optRows, err := p.db.Query(fmt.Sprintf(`SELECT name,value FROM option WHERE node_id = '%s';`, id))
		if err != nil {
			return nil, err
		}
		defer optRows.Close()

		opts := make(map[string][]string)

		for optRows.Next() {
			var name, value string
			optRows.Scan(&name, &value)
			if opt, exists := opts[name]; exists {
				opts[name] = append(opt, value)
			} else {
				opts[name] = []string{value}
			}
		}

		// Add host to result
		hosts = append(hosts, &Host{
			Name:    name,
			IPv4:    net.ParseIP(ipv4),
			IPv6:    net.ParseIP(ipv6),
			Options: opts,
		})
	}
	return hosts, rows.Err()
}

// Dump will dump ipplan host data to logs
func (p *IPPlan) Dump() error {
	rows, err := p.db.Query(`SELECT name,ipv4_addr_txt,ipv6_addr_txt FROM host;`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var ipv4 string
		var ipv6 string

		rows.Scan(&name, &ipv4, &ipv6)
		log.Infof("%s %s %s", name, net.ParseIP(ipv4), net.ParseIP(ipv6))
	}
	return rows.Err()
}
