package enforcer

import (
	"fmt"
	"net"
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
		if v4r {
			records = append(records, &Record{
				Name: fmt.Sprintf("%s.", reverseaddr(host.IPv4)),
				Type: "PTR",
				TTL:  e.Vars.HostTTL,
				Data: []string{host.Name},
			})
		}
		if v4f {
			records = append(records, &Record{
				Name: fmt.Sprintf("%s.", host.Name),
				Type: "A",
				TTL:  e.Vars.HostTTL,
				Data: []string{host.IPv4.String()},
			})
		}
		if v6r {
			records = append(records, &Record{
				Name: fmt.Sprintf("%s.", reverseaddr(host.IPv6)),
				Type: "PTR",
				TTL:  e.Vars.HostTTL,
				Data: []string{host.Name},
			})
		}
		if v6f {
			records = append(records, &Record{
				Name: fmt.Sprintf("%s.", host.Name),
				Type: "A",
				TTL:  e.Vars.HostTTL,
				Data: []string{host.IPv6.String()},
			})
		}
	}
	return records, nil
}

// Helpers to construct PTR records

const hexDigit = "0123456789abcdef"

func reverseaddr(ip net.IP) string {
	if ip.To4() != nil {
		return uitoa(uint(ip[15])) + "." + uitoa(uint(ip[14])) + "." + uitoa(uint(ip[13])) + "." + uitoa(uint(ip[12])) + ".in-addr.arpa."
	}
	buf := make([]byte, 0, len(ip)*4+len("ip6.arpa."))
	for i := len(ip) - 1; i >= 0; i-- {
		v := ip[i]
		buf = append(buf, hexDigit[v&0xF])
		buf = append(buf, '.')
		buf = append(buf, hexDigit[v>>4])
		buf = append(buf, '.')
	}
	return string(append(buf, "ip6.arpa."...))
}

func uitoa(val uint) string {
	if val == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf) - 1
	for val >= 10 {
		q := val / 10
		buf[i] = byte('0' + val - q*10)
		i--
		val = q
	}
	buf[i] = byte('0' + val)
	return string(buf[i:])
}
