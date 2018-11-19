
GOPATH=$(shell pwd)
BINARIES=bin/dugo bin/rmgo
LDFLAGS=-ldflags="-s -w"
GOOS=linux

.PHONY: all clean test

.EXPORT_ALL_VARIABLES:

all: bin $(BINARIES)

bin/dugo: $(wildcard src/fstools/dugo/*.go)
	@echo $(GOFILES)
	go install $(LDFLAGS) fstools/dugo

bin/rmgo: $(wildcard src/fstools/rmgo/*.go)
	go install $(LDFLAGS) fstools/rmgo

bin:
	mkdir $@

clean:
	rm $(BINARIES)

test:
	./scripts/test.sh
