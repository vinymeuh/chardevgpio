.DEFAULT_GOAL := help
.PHONY: build

GOFLAGS := -trimpath
GOBUILDFLAGS := ${GOFLAGS} -o ../../build
SUBDIRS := $(wildcard cmd/*/.)

builddir:
	@rm -rf build
	@mkdir build

build: builddir ## Build for current architecture
	$(foreach cmd, $(SUBDIRS), $(shell cd $(cmd) && GOOS=linux go build ${GOBUILDFLAGS}))
	@:

build-arm6: builddir ## Build for ARM6
	$(foreach cmd, $(SUBDIRS), $(shell cd $(cmd) && GOOS=linux GOARCH=arm GOARM=6 go build ${GOBUILDFLAGS}))
	@:

build-arm7: builddir ## Build for ARM7
	$(foreach cmd, $(SUBDIRS), $(shell cd $(cmd) && GOOS=linux GOARCH=arm GOARM=7 go build ${GOBUILDFLAGS}))
	@:

coverage: ## Show test coverage
	go tool cover -func=coverage.txt
	go tool cover -html=coverage.txt

test: ## Run tests
	go test -race -coverprofile=coverage.txt -covermode=atomic .

help: ## Show Help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@exit 0
