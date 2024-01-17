TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=cloudamqp
PROVIDER_VERSION = 1.29.2

default: build

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
PROVIDER_ARCH := $(GOOS)_$(GOARCH)

tools:
	GO111MODULE=on go install github.com/client9/misspell/cmd/misspell
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint

build: fmtcheck
	go install -ldflags "-X 'main.version=$(PROVIDER_VERSION)'"

local-clean:
	rm -rf ~/.terraform.d/plugins/localhost/cloudamqp/cloudamqp/$(PROVIDER_VERSION)/$(PROVIDER_ARCH)/terraform-provider-cloudamqp_v$(PROVIDER_VERSION)

local-build: local-clean
	@echo $(GOOS);
	@echo $(GOARCH);
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X 'main.version=$(PROVIDER_VERSION)'" -o terraform-provider-cloudamqp_v$(PROVIDER_VERSION)

local-install: local-build
	mkdir -p ~/.terraform.d/plugins/localhost/cloudamqp/cloudamqp/$(PROVIDER_VERSION)/$(PROVIDER_ARCH)
	cp $(CURDIR)/terraform-provider-cloudamqp_v$(PROVIDER_VERSION) ~/.terraform.d/plugins/localhost/cloudamqp/cloudamqp/$(PROVIDER_VERSION)/$(PROVIDER_ARCH)

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run ./...

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

.PHONY: build test testacc vet fmt fmtcheck lint tools test-compile
