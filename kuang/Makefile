.PHONY: build install test lint lint-fix security coverage clean help check ci install-tools

.DEFAULT_GOAL := help

help: ## Display available commands
	@echo "kuang Development Commands"
	@echo "=============================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

build: ## Build the kuang binary
	@go build -o kuang ./cmd

install: ## Build and install kuang to $GOPATH/bin
	@go install ./cmd

test: ## Run all tests with race detector
	@go test -v -race ./...

lint: ## Run linters
	@golangci-lint run --config=.golangci.yml --timeout=5m

lint-fix: ## Run linters with auto-fix
	@golangci-lint run --config=.golangci.yml --fix

security: ## Run security scanner
	@$(shell go env GOPATH)/bin/gosec -quiet ./...

coverage: ## Generate coverage report (HTML)
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out | tail -1
	@echo "Coverage report: coverage.html"

clean: ## Remove generated files
	@rm -f kuang coverage.out coverage.html
	@find . -name "*.test" -delete
	@find . -name "*.prof" -delete
	@find . -name "*.out" -delete

install-tools: ## Install development tools
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.7.2
	@go install github.com/securego/gosec/v2/cmd/gosec@latest

check: lint test security ## Run lint, tests, and security scan
	@echo "All checks passed!"

ci: clean check coverage ## Full CI simulation
	@echo "CI simulation complete!"
