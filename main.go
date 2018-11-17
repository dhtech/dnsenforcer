package main

import (
	"github.com/dhtech/dnsenforcer/enforcer"
	log "github.com/sirupsen/logrus"
)

func main() {
	e, err := enforcer.New("./ipplan.db")
	defer e.Close()
	if err != nil {
		log.Fatal(err)
	}
	e.PrintAllRecords()
}
