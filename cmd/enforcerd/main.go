package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/dhtech/dnsenforcer/enforcer"
	"github.com/dhtech/dnsenforcer/enforcer/ipplan"
	"github.com/ulikunitz/xz"
	log "github.com/sirupsen/logrus"
	pb "github.com/dhtech/proto/dns"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v2"
)

var (
	listenAddress = flag.String("listen", ":1215", "address to listen to")
)

type enforcerServer struct {
	v *enforcer.Vars
}

func (s *enforcerServer) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	ipplanXz, err := http.Get(os.Getenv("IPPLAN_XZ_URL"))
	if err != nil {
		log.Errorf("Unable to fetch IPPLAN_XZ_URL: %v", err)
		return nil, fmt.Errorf("failed to fetch ipplan.db.xz")
	}
	defer ipplanXz.Body.Close()

	ipplanDb, err := xz.NewReader(ipplanXz.Body)
	if err != nil {
		log.Errorf("Failed to decompress ipplan.db: %v", err)
		return nil, fmt.Errorf("failed to decompress ipplan.db")
	}

	ipplanFile, err := ioutil.TempFile("", "enforcerd")
	if err != nil {
		return nil, err
	}
	defer os.Remove(ipplanFile.Name())
	if _, err := io.Copy(ipplanFile, ipplanDb); err != nil {
		log.Errorf("Failed to write ipplan.db: %v", err)
		return nil, err
	}

	ipp, err := ipplan.Open(ipplanFile.Name())
	if err != nil {
		return nil, err
	}

	static, err := http.Get(os.Getenv("STATIC_URL"))
	if err != nil {
		log.Errorf("Unable to fetch STATIC_URL: %v", err)
		return nil, fmt.Errorf("failed to fetch static record map")
	}
	defer static.Body.Close()

	// Create new enforcer
	e, err := enforcer.New(s.v, ipp, static.Body)
	defer e.Close()
	if err != nil {
		return nil, err
	}

	added, removed, err := e.UpdateRecords()
	if err != nil {
		return nil, err
	}

	rev, err := ipp.Revision()
	if err != nil {
		log.Errorf("Could not get revision of ipplan: %v", err)
		rev = "<unknown>"
	}
	log.Info("Records updated to revision %s", rev)
	resp := &pb.RefreshResponse{
		Version: rev,
		Added:   uint32(added),
		Removed: uint32(removed),
	}
	return resp, nil
}

func main() {
	// Parse values
	vars := &enforcer.Vars{}
	flag.StringVar(&vars.Endpoint, "endpoint", "dns.net.dreamhack.se:443", "gRPC endpoint for DNS server")
	flag.StringVar(&vars.Certificate, "cert", "./client.pem", "Client certificate to use")
	flag.StringVar(&vars.Key, "key", "./key.pem", "Key to use")
	flag.IntVar(&vars.HostTTL, "host-ttl", 1337, "Default TTL to use for host records")
	vars.IgnoreTypes = strings.Split(*flag.String("ignore-types", "SOA,NS", "Do not remove or add these types of records"), ",")
	zonefile := flag.String("zones-file", "./zones.prod.yaml", "YAML fail with DNS zones to manage")
	flag.Parse()

	l, err := net.Listen("tcp", *listenAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Get data from zones file
	b, err := ioutil.ReadFile(*zonefile)
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

	s := &enforcerServer{vars}
	g := grpc.NewServer()
	pb.RegisterEnforcerServiceServer(g, s)
	reflection.Register(g)
	g.Serve(l)
}
