GOFILES = $(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES = $(shell go list ./...  | grep -v /amazonian/)

default: build

workdir:
	mkdir -p workdir

build: workdir/amazonian

workdir/amazonian: $(GOFILES)
	go build -o workdir/amazonian .

test: test-all

test-all:
	#@go test -v $(GOFILES)
	@go test -v ./...

test-all-min:
	@go test ./...