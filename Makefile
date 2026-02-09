APIM_USERNAME ?= admin
APIM_PASSWORD ?= admin
APIM_API1_USERNAME ?= api1
APIM_API1_PASSWORD ?= api1
APIM_SERVER_URL ?= http://localhost:30083/automation

.PHONY: speakeasy
speakeasy: ## Run speakeasy generation with curated examples and docs
	@rm -f terraform-provider-apim
	@mv ~/.terraformrc ~/.terraformrc.keep 2>/dev/null || true
	@terraform fmt -recursive > /dev/null
	@make doc-gen
	speakeasy run --output console --skip-versioning --skip-compile --skip-testing --skip-upload-spec
	@go mod tidy
	@rm -rf examples/data-sources docs/data-sources examples/README.md USAGE.md > /dev/null
	@mv ~/.terraformrc.keep ~/.terraformrc 2>/dev/null || true
	@go build

.PHONY: lint
lint: ## Run speakeasy lint accepting no error or warning
	@speakeasy lint openapi --schema schemas/automation-api-oas.yaml --max-validation-errors 0 --max-validation-warnings 0 --non-interactive
	@grep "// BEGIN GRAVITEE CLOUD INIT" internal/provider/provider.go > /dev/null || (echo "Cloud initializer code snippet appear to be missing" && exit 1)
	@terraform fmt -recursive -check || (echo "Error: Above terraform files are not properly formatted. Please run 'terraform fmt -recursive' to fix formatting issues" && exit 1)

.PHONY: lint-fix
lint-fix: ## Fix issues that can be found
	@terraform fmt -recursive

.PHONY: sync-oas
sync-oas: ## Copy OAS from APIM assuming the project is in ../gravitee-apim-management
		@curl -fsSL "https://raw.githubusercontent.com/gravitee-io/gravitee-api-management/refs/heads/master/gravitee-apim-rest-api/gravitee-apim-rest-api-automation/gravitee-apim-rest-api-automation-rest/src/main/resources/open-api.yaml" -o schemas/automation-api-oas.yaml

PRE_TEST_DIR = "$(shell pwd)/examples/use-cases/application-simple"

.PHONY: local-test-setup
local-test-setup:
	@./tests/examples/setup-local-provider.sh
	@go build

.PHONY: pre-test
pre-test: local-test-setup
	@echo "Validating resource creation with user ${APIM_USERNAME}"
	@cd $(PRE_TEST_DIR) && rm -rf terraform.state terraform.state.backup .terraform && \
	terraform apply -auto-approve 2>&1 > /tmp/tf.log || \
	(echo "Can't do terraform apply with ${APIM_USERNAME}" && cat /tmp/tf.log && exit 1)

	@echo "Validating resource destruction with user ${APIM_USERNAME}"
	@cd $(PRE_TEST_DIR) && \
	terraform apply -auto-approve -destroy 2>&1 > /tmp/tf.log || \
	(echo "Can't do terraform destroy with ${APIM_USERNAME}" && cat /tmp/tf.log && exit 1)

.PHONY: acceptance-tests
acceptance-tests: ## Run acceptance tests
	@APIM_USERNAME=${APIM_USERNAME} APIM_PASSWORD="$${APIM_PASSWORD}" APIM_SERVER_URL=${APIM_SERVER_URL} TF_ACC=1 go test -count=1 -v ./tests/acceptance


.PHONY: examples-tests
examples-tests: local-test-setup ## Run acceptance tests using examples
	@APIM_USERNAME=${APIM_USERNAME} APIM_PASSWORD="$${APIM_PASSWORD}" APIM_SERVER_URL=${APIM_SERVER_URL} TF_ACC=1 go test -count=1  -v ./tests/examples

.PHONY: unit-tests
unit-tests: ## Run unit tests
	@go test -count=1 ./internal/...


.PHONY: doc-gen
doc-gen: ## Generate Terraform examples docs
	@docker run --rm -u $$(id -u) -v ./.docgen/config:/config -v ./:/plugin graviteeio/doc-gen

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: start-cluster
start-cluster: ## Init and start a local cluster
	@cd hack/scripts && npx zx run-kind.mjs

.PHONY: stop-cluster
stop-cluster: ## Delete the local kind cluster
	@kind delete cluster --name gravitee
