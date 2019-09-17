UNAME_P := $(shell uname -p)
ifeq ($(UNAME_P),i386)
    GOARCH += 386
else
    ifeq ($(UNAME_P),AMD64)
        GOARCH += amd64
    endif

endif

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
    GOOS += linux
endif
ifeq ($(UNAME_S),Darwin)
    GOOS += darwin
endif

help:
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean:  ## Clean files
	rm -f ~/.terraform.d/plugins/terraform-provider-cloudamqp

depupdate: clean  ## Update all vendored dependencies
	dep ensure -update

build:  ## Build cloudamqp provider
	@echo $(GOOS);
	@echo $(GOARCH);
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o terraform-provider-cloudamqp

install: build  ## Install cloudamqp provider into terraform plugin directory
	cp $(CURDIR)/terraform-provider-cloudamqp ~/.terraform.d/plugins/
	mv $(CURDIR)/terraform-provider-cloudamqp $(CURDIR)/bin/

init: install  ## Run terraform init for local testing
	terraform init

.PHONY: help build install init
.DEFAULT_GOAL := help
