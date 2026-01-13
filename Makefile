build:
	go build -o ./bin/trace ./main.go

lint: ## Run the linter
	@golangci-lint run
	@go fmt ./...
	@go vet ./...