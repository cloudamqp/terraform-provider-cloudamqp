GOFMT_FILES?=$$(find . -name '*.go')

build: terraform-provider-cloudamqp

tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

terraform-provider-cloudamqp:
	go build -o terraform-provider-cloudamqp

install:
	go install .

fmt:
	gofmt -s -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	golangci-lint run ./...

clean:
	rm -f terraform-provider-cloudamqp

.PHONY: clean install fmt fmtcheck lint tools
