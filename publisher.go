package main

import (
	"database/sql"
	"github.com/golang/protobuf/proto"
	"github.com/xtracdev/goes"
	"github.com/xtracdev/orapub"
	"github.com/xtraclabs/appreg/domain"
	"time"
)

var domainEventProcessor orapub.EventProcessor

func init() {
	domainEventProcessor = orapub.EventProcessor{
		Initialize: func(db *sql.DB) error {
			return nil
		},
		Processor: func(db *sql.DB, event *goes.Event) error {
			if event.TypeCode == domain.AppRegCreatedCode {
				var payload domain.ApplicationRegistrationCreated
				err := proto.Unmarshal(event.Payload.([]byte), &payload)
				if err != nil {
					return err
				}

				ts := time.Unix(0, payload.CreateTimestamp)
				_, err = db.Exec("insert into app_summary (client_id, name, created) values(:1,:2,:3)",
					payload.AggregateId, payload.Name, ts)
				return err
			}
			return nil
		},
	}

	orapub.RegisterEventProcessor("appreg-created", domainEventProcessor)
}
