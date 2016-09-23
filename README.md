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

## Atom Event Publishing

For demo purposes this also publishes atom events as well - see the es-atom-feed-data
project for details.

## Build

This code has a cgo dependency based on the oracle driver it imports. Because of this, we
need to built it in a container to avoid the cross compiler limitation with go lang.

<pre>
docker run --rm -v "$PWD":/go/src/github.com/xtraclabs/cqrs-sample-pub -w /go/src/github.com/xtraclabs/cqrs-sample-pub xtracdev/goora bash -c "make -f Makefile.docker"
</pre>

After building the binary, create the docker image using `make`, then `docker push` it to docker hub.