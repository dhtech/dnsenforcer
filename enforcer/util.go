package enforcer

import (
	"net"
	"strings"
)

// Check if domain is a part of a zone managed by enforcer
func (e *Enforcer) inZones(domain string) bool {
	for _, z := range e.Vars.Zones {
		if strings.HasSuffix(domain, z) {
			return true
		}
	}
	return false
}

// Check if a list of records contains a given record
func contains(list []*Record, r *Record) bool {
	for _, record := range list {
		if compare(r, record) {
			return true
		}
	}
	return false
}

// Check if two records are the same
func compare(a, b *Record) bool {
	if a == nil || b == nil {
		return false
	}
	if a.Name != b.Name {
		return false
	}
	if a.Type != b.Type {
		return false
	}
	if a.TTL != b.TTL {
		return false
	}
	if len(a.Data) != len(b.Data) {
		return false
	}
	for i := 0; i < len(a.Data); i++ {
		if a.Data[i] != b.Data[i] {
			return false
		}
	}
	return true
}

// PTR

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
