package sampledomain

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/xtracdev/goes"
	"log"
	"time"
)

type ApplicationReg struct {
	*goes.Aggregate
	Name        string
	Description string
	Created     int64 //Unix time stamp serialized as an int64
}

func (ar *ApplicationReg) String() string {
	ts := time.Unix(0, ar.Created).Format(time.RFC3339Nano)
	return fmt.Sprintf("ID: %s, Name: %s, Description: %s, Created: %s",
		ar.ID, ar.Name, ar.Description, ts)
}

const (
	AppRegCreatedCode = "ARCRE"
)

var (
	ErrUnknownType = errors.New("Unknown event type")
)

func NewApplicationReg(name, description string) (*ApplicationReg, error) {
	var appReg = new(ApplicationReg)
	appReg.Aggregate = goes.NewAggregate()
	appReg.Version = 1

	appRegCreated := ApplicationRegistrationCreated{
		AggregateId:     appReg.ID,
		Name:            name,
		Description:     description,
		CreateTimestamp: time.Now().UnixNano(),
	}

	appReg.Apply(
		goes.Event{
			Source:  appReg.ID,
			Version: appReg.Version,
			Payload: appRegCreated,
		})

	return appReg, nil
}

func NewApplicationRegFromHistory(events []goes.Event) *ApplicationReg {

	log.Printf("New application reg from history - %d events\n", len(events))
	appReg := new(ApplicationReg)
	appReg.Aggregate = goes.NewAggregate()

	unmarshalledEvents, err := unmarshallEvents(events)
	if err != nil {
		return nil
	}

	for _, e := range unmarshalledEvents {
		log.Println("apply event", e)
		appReg.Version += 1
		appReg.Route(e)
	}

	return appReg
}

func (ar *ApplicationReg) Apply(event goes.Event) {
	ar.Route(event)
	ar.Events = append(ar.Events, event)
}

func (ar *ApplicationReg) Route(event goes.Event) {
	event.Version = ar.Version
	switch event.Payload.(type) {
	case ApplicationRegistrationCreated:
		ar.handleApplicationRegistrationCreated(event.Payload.(ApplicationRegistrationCreated))
	default:
		log.Printf("unexpected type handled: %t", event.Payload)
	}
}

func (ar *ApplicationReg) handleApplicationRegistrationCreated(event ApplicationRegistrationCreated) {
	ar.ID = event.AggregateId
	ar.Name = event.Name
	ar.Description = event.Description
	ar.Created = event.CreateTimestamp
}

func (ar *ApplicationReg) Store(eventStore goes.EventStore) error {
	marshalled, err := marshallEvents(ar.Events)
	if err != nil {
		return nil
	}

	log.Println("Storing ", len(ar.Events), " events.")

	aggregateToStore := &goes.Aggregate{
		ID:      ar.ID,
		Version: ar.Version,
		Events:  marshalled,
	}

	err = eventStore.StoreEvents(aggregateToStore)
	if err != nil {
		return err
	}

	ar.Events = make([]goes.Event, 0)

	return nil
}

func marshallEvents(events []goes.Event) ([]goes.Event, error) {

	var updatedEvents []goes.Event

	for _, e := range events {

		var err error
		var newEvent goes.Event
		newEvent.Source = e.Source
		newEvent.Version = e.Version

		switch e.Payload.(type) {
		case ApplicationRegistrationCreated:
			newEvent.TypeCode = AppRegCreatedCode
			newEvent.Payload, err = marshallCreate(e.Payload.(ApplicationRegistrationCreated))
			if err != nil {
				return nil, err
			}

		default:
			return nil, ErrUnknownType
		}

		updatedEvents = append(updatedEvents, newEvent)
	}

	return updatedEvents, nil
}

func marshallCreate(create ApplicationRegistrationCreated) ([]byte, error) {
	return proto.Marshal(&create)
}

func unmarshallCreated(bytes []byte) (ApplicationRegistrationCreated, error) {
	var payload ApplicationRegistrationCreated
	err := proto.Unmarshal(bytes, &payload)
	return payload, err
}

func unmarshallEvents(events []goes.Event) ([]goes.Event, error) {
	var unmarshalled []goes.Event

	for _, e := range events {

		var err error
		var newEvent goes.Event
		newEvent.Source = e.Source
		newEvent.Version = e.Version
		newEvent.TypeCode = e.TypeCode

		switch e.TypeCode {
		case AppRegCreatedCode:
			newEvent.Payload, err = unmarshallCreated(e.Payload.([]byte))
			if err != nil {
				return nil, err
			}
		default:
			log.Println("Warning: unknown type code in event history", e.TypeCode)
		}

		unmarshalled = append(unmarshalled, newEvent)
	}

	return unmarshalled, nil
}
