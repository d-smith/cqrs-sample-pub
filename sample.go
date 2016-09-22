package main

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/xtraclabs/cqrs-sample-pub/sampledomain"
	dp "github.com/xtraclabs/es-data-pub"
)

func main() {
	if err := dp.ProcessEventRecords(); err != nil {
		log.Fatal(err.Error())
	}
}
