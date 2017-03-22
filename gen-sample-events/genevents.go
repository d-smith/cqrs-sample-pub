package main

import (
	"github.com/xtracdev/oraeventstore"
	"log"
	"os"
	"strconv"
	"time"
	"strings"
	"fmt"
	"github.com/xtraclabs/appreg/domain"
	"database/sql"
	_ "github.com/mattn/go-oci8"
)

var user, password, dbhost, dbPort, dbSvc string

func init() {
	var configErrors []string

	user = os.Getenv("DB_USER")
	if user == "" {
		configErrors = append(configErrors, "Configuration missing DB_USER env variable")
	}

	password = os.Getenv("DB_PASSWORD")
	if password == "" {
		configErrors = append(configErrors, "Configuration missing DB_PASSWORD env variable")
	}

	dbhost = os.Getenv("DB_HOST")
	if dbhost == "" {
		configErrors = append(configErrors, "Configuration missing DB_HOST env variable")
	}

	dbPort = os.Getenv("DB_PORT")
	if dbPort == "" {
		configErrors = append(configErrors, "Configuration missing DB_PORT env variable")
	}

	dbSvc = os.Getenv("DB_SVC")
	if dbSvc == "" {
		configErrors = append(configErrors, "Configuration missing DB_SVC env variable")
	}

	if len(configErrors) != 0 {
		log.Fatal(strings.Join(configErrors, "\n"))
	}

}



func main() {

	if len(os.Args) != 3 {
		log.Fatalf("Usage: go run genevents.go <num aggregates> <delay ms>")
	}

	numAggregates, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}

	delay, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err.Error())
	}

	os.Setenv("ES_PUBLISH_EVENTS", "1")

	var connectStr = fmt.Sprintf("%s/%s@//%s:%s/%s", user, password, dbhost, dbPort, dbSvc)
	db, err := sql.Open("oci8", connectStr)
	if err != nil {
		log.Fatalf("Error connecting to oracle: %s", err.Error())
	}

	//Are we really in an ok state for starters?
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to oracle: %s", err.Error())
	}

	eventStore, err := oraeventstore.NewOraEventStore(db)
	if err != nil {
		log.Fatalf("Error connecting to oracle: %s", err.Error())
	}

	for i := 0; i < numAggregates; i++ {

		app, _ := domain.NewApplicationReg(fmt.Sprintf("app %d", i), "app desc")

		err = app.Store(eventStore)
		if err != nil {
			log.Fatalf("Error storing events: %s", err.Error())
		}

		log.Println("sleep")
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
}
