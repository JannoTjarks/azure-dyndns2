.PHONY: build test fmt clean

BINARY_NAME := azure-dyndns2
GO_FLAGS := -mod=vendor

build: ## Build the binary
	go build $(GO_FLAGS) -o $(BINARY_NAME) .

test: ## Run tests with coverage and race detector
	go test $(GO_FLAGS) -race -coverprofile=coverage.out -v $$(go list ./... | grep -v "/cmd$$")
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

fmt: ## Format all go files
	go fmt ./...

clean: ## Remove built files
	rm -f $(BINARY_NAME) coverage.out
