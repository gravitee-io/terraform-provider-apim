VERSION=$(shell cat VERSION)

APIM_USERNAME ?= admin
APIM_PASSWORD ?= admin
APIM_SERVER_URL ?= http://localhost:30083/automation

.PHONY: speakeasy
speakeasy: ## Run speakeasy generation with curated examples and docs
	@rm -rf examples/data-sources
	speakeasy run --set-version $(VERSION) --output console --minimal
	@go mod tidy
	# @go generate .
	# @git checkout -- README.md
	@git clean -fd examples/data-sources docs/data-sources examples/README.md USAGE.md > /dev/null

.PHONY: acceptance-tests
acceptance-tests: ## Run acceptance tests against a running environment
	@APIM_USERNAME=${APIM_USERNAME} APIM_PASSWORD=${APIM_PASSWORD} APIM_SERVER_URL=${APIM_SERVER_URL} TF_ACC=1 go test -count=1 -timeout=10m -v ./internal/provider


.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
