GOFMT_FILES?=$$(find . -name '*.go')

default: build

tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run ./...

.PHONY: build test testacc vet fmt fmtcheck lint tools test-compile
