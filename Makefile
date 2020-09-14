export GOPATH := $(shell go env GOPATH)

default: install

build: export CGO_ENABLED = 0
build:
	go build -tags netgo -o build

install: export CGO_ENABLED = 0
install:
	@rm -f ${GOPATH}/bin/interpolate
	@go install -tags netgo
	@echo installed to ${GOPATH}/bin/interpolate

test: export CGO_ENABLED = 0
test:
	@go test -tags netgo --count 1 ./...

race: export CGO_ENABLED = 1
race:
	go test -tags netgo -race --count 1 ./...

bench: export CGO_ENABLED = 1
bench:
	go test -tags netgo --count 1 -run Bench -bench=. ./...

lint:
	@golangci-lint run

clean:
	@rm -f ./build

.PHONY: build test race bench lint install clean
