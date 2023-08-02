default: testacc

.PHONY: help testacc install

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
PLATFORM := $(GOOS)_$(GOARCH)
GOPATH := $(shell go env GOPATH)

# Export make variables (and those defined in the local.env file) into the environment of sub-processes
include env/local.env
.EXPORT_ALL_VARIABLES:

## install: builds the provider project and puts it in the relevant terraform plugin directory
install:
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/resourcely-inc/resourcely/0.0.1/$(PLATFORM)
	rm -f ~/.terraform.d/plugins/registry.terraform.io/resourcely-inc/resourcely/0.0.1/$(PLATFORM)/terraform-provider-resourcely_v0.0.1
	go install
	cp $(GOPATH)/bin/terraform-provider-resourcely ~/.terraform.d/plugins/registry.terraform.io/resourcely-inc/resourcely/0.0.1/$(PLATFORM)/terraform-provider-resourcely_v0.0.1

## testacc: Run acceptance tests (with env variables from env/local.env)
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

## codegen: Update auto-generated code
codegen:
	go mod tidy
	go generate ./...

## help: Prints this help message.
help:
	@echo "Usage: "
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
