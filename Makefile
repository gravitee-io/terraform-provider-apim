APIM_USERNAME ?= admin
APIM_PASSWORD ?= admin
APIM_SERVER_URL ?= http://localhost:30083/automation

.PHONY: speakeasy
speakeasy: ## Run speakeasy generation with curated examples and docs
	speakeasy run --output console --minimal
	@go mod tidy
	@rm -rf examples/data-sources docs/data-sources examples/README.md USAGE.md > /dev/null
	@go build
	@git restore docs/guides

.PHONY: lint
lint: ## Run speakeasy lint accepting no error or warning
	@speakeasy lint openapi --schema schemas/automation-api-oas.yaml --max-validation-errors 0 --max-validation-warnings 0 --non-interactive
	@git diff --quiet HEAD docs/guides || { \
		echo "Error: Documentation generation produced changes. Please run 'make doc-gen' and commit the updated documentation."; \
		exit 1; \
	}


.PHONY: acceptance-tests
acceptance-tests: ## Run acceptance tests against a running environment
	@APIM_USERNAME=${APIM_USERNAME} APIM_PASSWORD=${APIM_PASSWORD} APIM_SERVER_URL=${APIM_SERVER_URL} TF_ACC=1 go test -v ./internal/provider


.PHONY: doc-gen
doc-gen: ## Generate Terraform examples docs
	@docker run --rm -v ./.docgen/config:/config -v ./:/plugin graviteeio/doc-gen

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)