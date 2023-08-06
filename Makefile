BIN = podzol
MODULE := $(shell go list -m)
VERSION := $(shell git describe --tags --always --dirty --broken)

LDFLAGS := -s -w -X $(MODULE)/pkg.Version=$(VERSION)

.PHONY: all test $(BIN)

all: $(BIN)

$(BIN):
	go build -o $@ -ldflags='$(LDFLAGS)'

test:
	go test -v ./...
