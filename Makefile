# Commands
GO ?= `which go`
DOCKER ?= `which docker`

.PHONY: help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "make \033[36m%-10s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## Remove temporary resources
	@rm -f tee-mock-server
	@rm -f ./pki/*.pem

.PHONY: build
build: ## Build the TEE Mock Server static binary (for Linux only)
	CGO_ENABLED=0 GOOS=linux @$(GO) build -o tee-mock-server

.PHONY: docker-build
docker-build: ## Build the TEE Mock Server docker image
	@$(DOCKER) build --tag tee-server-mock .

.PHONY: docker-run
docker-run: ## Start the TEE Mock Server using docker
	@$(DOCKER) run --name "tee-mock-server" --rm -it -v /run/container_launcher:/run/container_launcher tee-mock-server

.PHONY: chain
chain: ## Generate a new certificate chain
	(cd ./pki && ./create-chain.sh)