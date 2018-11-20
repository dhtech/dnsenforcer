package enforcer

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"sync"

	"github.com/dhtech/proto/dns"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// UpdateRecords logs all records to stdout
func (e *Enforcer) UpdateRecords() error {
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

	// Get and convert all remote records to local style records
	onlineRecords := make([]*Record, 0)
	var mutex = &sync.Mutex{}
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

			if res.Record != nil {
				for _, record := range res.Record {
					if !e.ignoredType(record.Type) {
						mutex.Lock()
						onlineRecords = append(onlineRecords, &Record{
							Name: record.Domain,
							TTL:  int(record.Ttl),
							Type: record.Type,
							Data: []string{record.Data},
						})
						mutex.Unlock()
					}
				}
			}
			wg.Done()
		}(zone)
	}

	wg.Wait()

	// Get localally constructed records
	localRecords, err := e.GetAllRecords()
	if err != nil {
		return err
	}

	// Find which records to remove
	remove := make([]*dns.Record, 0)
	for _, r := range onlineRecords {
		if !contains(localRecords, r) {
			// Delete
			for _, d := range r.Data {
				remove = append(remove, &dns.Record{
					Domain: r.Name,
					Type:   r.Type,
					Data:   d,
				})
			}
		}
	}

	// Remove records that are present on server but no locally
	if !e.Vars.DryRun {
		log.Infof("Deleting %d records", len(remove))
		for _, r := range remove {
			if _, err := c.Remove(ctx, &dns.RemoveRequest{Record: []*dns.Record{r}}); err != nil {
				log.Errorf("Remove of %s failed with %v", r.Domain, err)
			} else {
				log.Infof("Removed %s", r.Domain)
			}
		}
	} else {
		for _, r := range remove {
			fmt.Printf("-%s:%s:%d:%s\n", r.Domain, r.Type, r.Ttl, r.Data)
		}
	}

	// Find which records to insert
	insert := make([]*dns.Record, 0)
	for _, r := range localRecords {
		if !contains(onlineRecords, r) {
			// Add
			for _, d := range r.Data {
				insert = append(insert, &dns.Record{
					Domain: r.Name,
					Ttl:    uint32(r.TTL),
					Class:  "IN",
					Type:   r.Type,
					Data:   d,
				})
			}
		}
	}

	// Insert records that are missing on the server
	if !e.Vars.DryRun {
		log.Infof("Inserting %d records", len(insert))
		for _, r := range insert {
			if _, err := c.Insert(ctx, &dns.InsertRequest{Record: []*dns.Record{r}}); err != nil {
				log.Errorf("Insert of %s failed with %v", r.Domain, err)
			} else {
				log.Infof("Added %s", r.Domain)
			}
		}
	} else {
		for _, r := range insert {
			fmt.Printf("+%s:%s:%d:%s\n", r.Domain, r.Type, r.Ttl, r.Data)
		}
	}

	return nil
}
