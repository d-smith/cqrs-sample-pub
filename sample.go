package main

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/xtraclabs/es-atom-feed-data"
	dp "github.com/xtracdev/es-data-pub"
)

func main() {
	if err := dp.ProcessEventRecords(); err != nil {
		log.Fatal(err.Error())
	}
}
