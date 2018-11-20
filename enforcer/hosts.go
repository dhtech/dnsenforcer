package enforcer

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// GetHostRecords will use IPPlan to return all reverse and forward records
func (e *Enforcer) GetHostRecords() ([]*Record, error) {
	records := make([]*Record, 0)
	hosts, err := e.IPPlan.Hosts()
	if err != nil {
		return nil, err
	}
	for _, host := range hosts {
		// Using options to determine which records to create for this host
		v4r, v4f, v6r, v6f := false, false, false, false
		_, v4r = host.Options["ipv4r"]
		_, v4f = host.Options["ipv4f"]
		_, v6r = host.Options["ipv6r"]
		_, v6f = host.Options["ipv6f"]
		if !(v4r || v4f || v6r || v6f) {
			// No DNS options will create all records
			v4r, v4f, v6r, v6f = true, true, true, true
		}

		// Create records
		if v4r && host.IPv4 != nil {
			h := reverseaddr(host.IPv4)
			if !e.inZones(h) {
				log.Warningf("%s is not a member of any of the enforced zones", h)
				continue
			}
			records = append(records, &Record{
				Name: h,
				Type: "PTR",
				TTL:  e.Vars.HostTTL,
				Data: []string{fmt.Sprintf("%s.", host.Name)},
			})
		}
		if v4f && host.IPv4 != nil {
			h := fmt.Sprintf("%s.", host.Name)
			if !e.inZones(h) {
				log.Warningf("%s is not a member of any of the enforced zones", h)
				continue
			}
			records = append(records, &Record{
				Name: h,
				Type: "A",
				TTL:  e.Vars.HostTTL,
				Data: []string{host.IPv4.String()},
			})
		}
		if v6r && host.IPv6 != nil {
			h := reverseaddr(host.IPv6)
			if !e.inZones(h) {
				log.Warningf("%s is not a member of any of the enforced zones", h)
				continue
			}
			records = append(records, &Record{
				Name: h,
				Type: "PTR",
				TTL:  e.Vars.HostTTL,
				Data: []string{fmt.Sprintf("%s.", host.Name)},
			})
		}
		if v6f && host.IPv6 != nil {
			h := fmt.Sprintf("%s.", host.Name)
			if !e.inZones(h) {
				log.Warningf("%s is not a member of any of the enforced zones", h)
				continue
			}
			records = append(records, &Record{
				Name: h,
				Type: "AAAA",
				TTL:  e.Vars.HostTTL,
				Data: []string{host.IPv6.String()},
			})
		}
	}
	return records, nil
}
