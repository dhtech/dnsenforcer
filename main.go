package main

import (
	"flag"
	"io/ioutil"

	"github.com/dhtech/dnsenforcer/enforcer"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func main() {
	// Pars values
	vars := &enforcer.Vars{}
	flag.StringVar(&vars.Endpoint, "endpoint", "dns.net.dreamhack.se:443", "gRPC endpoint for DNS server")
	flag.StringVar(&vars.Certificate, "cert", "./client.pem", "Client certificate to use")
	flag.StringVar(&vars.Key, "key", "./key.pem", "Key to use")
	flag.StringVar(&vars.DBFile, "ipplan", "./ipplan.db", "Path to ipplan file to use")
	flag.StringVar(&vars.Static, "static", "./static.prod.yaml", "Path to static file to use")
	flag.IntVar(&vars.HostTTL, "host-ttl", 1337, "Default TTL to use for host records")
	flag.BoolVar(&vars.DryRun, "dry-run", false, "Do not actually update records on the DNS server")
	export := flag.Bool("export", false, "Will write to static file rather than read all the records present in the server and not ipplan")
	flag.Parse()

	// Get data from zones file
	b, err := ioutil.ReadFile(*flag.String("zones-file", "./zones.prod.yaml", "YAML fail with DNS zones to manage"))
	if err != nil {
		log.Error("You need to create a zone config file")
		log.Fatal(err)
	}
	var zones struct {
		Zones []string `yaml:"zones"`
	}
	err = yaml.Unmarshal(b, &zones)
	if err != nil {
		log.Error("You need to create a zone config file")
		log.Fatal(err)
	}
	vars.Zones = zones.Zones

	log.Info("Generating DNS records...")

	// Create new enforcer
	e, err := enforcer.New(vars)
	defer e.Close()
	if err != nil {
		log.Fatal(err)
	}

	if *export {
		err = e.ExportStaticRecords()
		if err != nil {
			log.Fatal(err)
		}
		log.Info("Records exported")
	} else {
		err = e.UpdateRecords()
		if err != nil {
			log.Fatal(err)
		}
		log.Info("Records updated")
	}

}
