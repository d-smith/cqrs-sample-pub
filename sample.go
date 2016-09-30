package main

import (
	_ "github.com/xtraclabs/es-atom-feed-data"
	dp "github.com/xtracdev/es-data-pub"
)

func main() {
	dp.ProcessEventRecords()
}
