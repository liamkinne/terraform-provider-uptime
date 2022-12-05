default: testacc

.PHONY: help
help: ## self documenting help output
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: gen
gen: ## generate docs and openapi client code
	go generate ./... -v

# Run acceptance tests
.PHONY: testacc
testacc: ## run tests
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: gettoken
gettoken: ## open browser to get token
	open https://uptime.com/api/tokens
