.PHONY: build build-container test fmt clean

BINARY_NAME := azure-dyndns2
GO_FLAGS := -mod=vendor

build: ## Build the binary
	go build $(GO_FLAGS) -o $(BINARY_NAME) .

build-container:
	podman build -f Dockerfile -t jannotjarks/azure-dyndns2:latest

test: ## Run tests with coverage and race detector
	go test $(GO_FLAGS) -race -coverprofile=coverage.out -v $$(go list ./... | grep -v "/cmd$$")
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

fmt: ## Format all go files
	go fmt ./...

clean: ## Remove built files
	rm -f $(BINARY_NAME) coverage.out coverage.html

integration: build
	AZURE_DYNDNS_DNS_ZONE_NAME=test ./$(BINARY_NAME) one-shot --hostname integration-test.tjarks.dev --myip 100.64.0.1 --dns-zone-name tjarks.dev --dns-resource-group-name rg-publicdns --dns-subscription-id cf4ea9ae-cfef-4132-a5a1-c507a07a3371
