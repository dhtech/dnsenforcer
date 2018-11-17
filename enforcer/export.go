package enforcer

import (
	"context"
	"crypto/tls"
	"net"
	"os"
	"regexp"
	"sync"

	"github.com/dhtech/proto/dns"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/yaml.v2"
)

// ExportStaticRecords will take all server side records and export them to static YAML records
func (e *Enforcer) ExportStaticRecords() error {
	// Client Auth
	certificate, err := tls.LoadX509KeyPair(e.Vars.Certificate, e.Vars.Key)
	if err != nil {
		return err
	}

	host, _, err := net.SplitHostPort(e.Vars.Endpoint)
	if err != nil {
		return err
	}

	creds := credentials.NewTLS(&tls.Config{
		ServerName:   host,
		Certificates: []tls.Certificate{certificate},
	})

	// gRPC connection
	conn, err := grpc.Dial(e.Vars.Endpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	ctx := context.Background()
	c := dns.NewDynamicDnsServiceClient(conn)

	// Convert all remote records to local style records
	onlineRecords := make([]*Record, 0)
	var wg sync.WaitGroup
	for _, zone := range e.Vars.Zones {
		wg.Add(1)
		go func(zone string) {
			res, err := c.GetZone(ctx, &dns.GetZoneRequest{
				Zone: zone,
			})
			if err != nil {
				log.Error(err)
			}

			for _, record := range res.Record {
				onlineRecords = append(onlineRecords, &Record{
					Name: record.Domain,
					TTL:  int(record.Ttl),
					Type: record.Type,
					Data: []string{record.Data},
				})
			}
			wg.Done()
		}(zone)
	}

	wg.Wait()

	// Open data file
	data, err := os.Create(e.Vars.Static)
	if err != nil {
		return err
	}
	writer := yaml.NewEncoder(data)

	// Filter generated records
	localRecords, err := e.GetAllRecords()
	if err != nil {
		return err
	}

outer:
	for _, online := range onlineRecords {
		for _, local := range localRecords {
			// Filter out forward and reverse
			re := regexp.MustCompile(`(\d+-\d+-\d+-\d+.*|.*\.in-addr\.arpa\.)`)
			if online == nil || local == nil {
				continue outer
			}
			if online.Name == local.Name || online.Type == "SOA" || online.Type == "PTR" || re.MatchString(online.Name) {
				continue outer
			}
		}
		writer.Encode(online)
	}

	data.Close()

	return nil
}
