package enforcer

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/dhtech/proto/dns"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// UpdateRecords logs all records to stdout
func (e *Enforcer) UpdateRecords() error {
	records, err := e.GetAllRecords()
	if err != nil {
		return err
	}

	if e.Vars.DryRun {
		for _, record := range records {
			fmt.Printf("Dry-run: %s,%s,%d,%v\n", record.Name, record.Type, record.TTL, record.Data)
		}
		return nil
	}

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

	// Get Zone
	c := dns.NewDynamicDnsServiceClient(conn)

	res, err := c.GetZone(ctx, &dns.GetZoneRequest{
		Zone: "tech.dreamhack.se.",
	})
	if err != nil {
		return err
	}

outer:
	for _, server := range res.Record {
		for _, local := range records {
			if server.Domain == local.Name {
				continue outer
			}
		}
		fmt.Printf("-%s,%s\n", server.Domain, server.Type)
	}

	return nil
}
