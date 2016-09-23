# Event Publishing - CQRS Sample

This project shows how to update a CQRS read side table as events
in a domain are published from an event store.

To run this you need a database set up with the [oraeventstore](https://github.com/xtracdev/oraeventstore)
schema, plus the following:

<pre>
create table app_summary (
    client_id varchar2(100) not null,
    name varchar2(100) not null,
    created timestamp not null,
    primary key(client_id)
)
 
create or replace synonym esusr.app_summary for esdbo.app_summary;
 
grant insert, delete, select on app_summary to esusr;
</pre>

With this in place, the published can be run vi `go run sample.go`

To publish events, `go run genevents.go` in the gen-sample-events 
directory.