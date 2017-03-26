image:
	docker run --rm -v $(CURDIR)/../:/go/src/github.com/xtraclabs/cqrs-sample-pub -w /go/src/github.com/xtraclabs/cqrs-sample-pub xtracdev/goora bash -c "cd gen-sample-events; make -f Makefile.docker"
	docker build -t xtracdev/genevents .