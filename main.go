package main

import (
	"github.com/dhtech/dnsenforcer/enforcer"
	log "github.com/sirupsen/logrus"
)

func main() {
	_, err := enforcer.New("./ipplan.db")
	if err != nil {
		log.Fatal(err)
	}
}
