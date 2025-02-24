GOFMT_FILES?=$$(find . -name '*.go')

default: build

tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

build:
	go build .

install:
	go install .

fmt:
	gofmt -s -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	golangci-lint run ./...

.PHONY: build install fmt fmtcheck lint tools
