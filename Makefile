APIM_USERNAME ?= admin
APIM_PASSWORD ?= admin
APIM_API1_USERNAME ?= api1
APIM_API1_PASSWORD ?= api1
APIM_SERVER_URL ?= http://localhost:30083/automation

.PHONY: speakeasy
speakeasy: ## Run speakeasy generation with curated examples and docs
	speakeasy run --output console --minimal --skip-versioning
	@go mod tidy
	@rm -rf examples/data-sources docs/data-sources examples/README.md USAGE.md > /dev/null
	go build
	cd examples && terraform fmt -recursive

.PHONY: lint
lint: ## Run speakeasy lint check run terraform resource formatting
	speakeasy lint openapi --schema schemas/automation-api-oas.yaml --max-validation-errors 0 --max-validation-warnings 0 --non-interactive
	@terraform fmt -recursive -check || (echo "Error: Above terraform files are not properly formatted. Please run 'terraform fmt -recursive' to fix formatting issues" && exit 1)

.PHONY: lint-fix
lint-fix: ## Fix issues that can be found
	@terraform fmt -recursive

.PHONY: sync-oas
sync-oas: ## Copy OAS from APIM assuming the project is in ../gravitee-apim-management
	@cp ../gravitee-api-management/gravitee-apim-rest-api/gravitee-apim-rest-api-automation/gravitee-apim-rest-api-automation-rest/src/main/resources/open-api.yaml schemas/automation-api-oas.yaml

PRE_TEST_DIR = "$(shell pwd)/examples/use-cases/application-simple"

.PHONY: pre-test
pre-test:
# First user (APIM_USERNAME)
	@echo "Validating resource creation with user ${APIM_USERNAME}"
	@cd $(PRE_TEST_DIR) && rm -rf terraform.state terraform.state.backup .terraform && \
	APIM_USERNAME=${APIM_USERNAME} \
	APIM_PASSWORD=${APIM_PASSWORD} \
	APIM_SERVER_URL=${APIM_SERVER_URL} \
	terraform apply -auto-approve 2>&1 > /tmp/tf.log || \
    (echo "Can't do terraform apply with ${APIM_USERNAME}" && cat /tmp/tf.log && exit 1)

	@echo "Validating resource destruction with user ${APIM_USERNAME}"
	@cd $(PRE_TEST_DIR) && \
	APIM_USERNAME=${APIM_USERNAME} \
	APIM_PASSWORD=${APIM_PASSWORD} \
	APIM_SERVER_URL=${APIM_SERVER_URL} \
	terraform apply -auto-approve -destroy 2>&1 > /tmp/tf.log || \
    (echo "Can't do terraform destroy with ${APIM_USERNAME}" && cat /tmp/tf.log && exit 1)

# Second user (APIM_API1_USERNAME)
	@echo "Validating resource creation with user ${APIM_API1_USERNAME}"
	@cd $(PRE_TEST_DIR) && rm -rf terraform.state terraform.state.backup .terraform && \
	APIM_USERNAME=${APIM_API1_USERNAME} \
	APIM_PASSWORD=${APIM_API1_PASSWORD} \
	APIM_SERVER_URL=${APIM_SERVER_URL} \
	terraform apply -auto-approve 2>&1 > /tmp/tf.log || \
    (echo "Can't do terraform apply with ${APIM_API1_USERNAME}" && cat /tmp/tf.log && exit 1)

	@echo "Validating resource destruction with user ${APIM_API1_USERNAME}"
	@cd $(PRE_TEST_DIR) && \
	APIM_USERNAME=${APIM_API1_USERNAME} \
	APIM_PASSWORD=${APIM_API1_PASSWORD} \
	APIM_SERVER_URL=${APIM_SERVER_URL} \
	terraform apply -auto-approve -destroy 2>&1 > /tmp/tf.log || \
    (echo "Can't do terraform destroy with ${APIM_API1_USERNAME}" && cat /tmp/tf.log && exit 1)

.PHONY: acceptance-tests
acceptance-tests: pre-test ## Run acceptance tests
	@APIM_USERNAME=${APIM_USERNAME} APIM_PASSWORD=${APIM_PASSWORD} APIM_SERVER_URL=${APIM_SERVER_URL} TF_ACC=1 go test -v ./tests/acceptance


.PHONY: examples-tests
examples-tests: pre-test ## Run acceptance tests using examples
	@APIM_USERNAME=${APIM_USERNAME} APIM_PASSWORD=${APIM_PASSWORD} APIM_SERVER_URL=${APIM_SERVER_URL} TF_ACC=1 go test -v ./tests/examples


.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)