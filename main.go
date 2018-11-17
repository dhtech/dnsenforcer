package main

import (
	"flag"

	"github.com/dhtech/dnsenforcer/enforcer"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Pars values
	vars := &enforcer.Vars{}
	flag.StringVar(&vars.DBFile, "ipplan", "./ipplan.db", "Path to ipplan file to use")
	flag.StringVar(&vars.Static, "static", "./static.yaml", "Path to static file to use")
	flag.IntVar(&vars.HostTTL, "host-ttl", 60, "Default TTL to use for host records")
	flag.BoolVar(&vars.DryRun, "dry-run", false, "Do not actually update records on the DNS server")
	flag.Parse()

	log.Info("Generating DNS records...")

	// Create new enforcer
	e, err := enforcer.New(vars)
	defer e.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = e.UpdateRecords()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Records updated")
}
