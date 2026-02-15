.PHONY: help deps build deploy clean test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Install Go dependencies
	@echo "Installing dependencies..."
	@go get github.com/aws/aws-lambda-go/lambda
	@go get github.com/aws/aws-lambda-go/events
	@go get github.com/aws/aws-sdk-go-v2/service/kms
	@go mod tidy
	@echo "Dependencies installed!"

build: deps ## Build Lambda functions
	@echo "Building Lambda functions..."
	@./build-lambda.sh

deploy: build ## Build and deploy to AWS
	@echo "Deploying to AWS..."
	@./deploy.sh

test: ## Run tests
	@go test ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf lambda/
	@rm -f tfplan
	@echo "Clean complete!"

tf-init: ## Initialize Terraform
	@terraform init

tf-plan: build ## Run Terraform plan
	@terraform plan -out=tfplan

tf-apply: ## Apply Terraform plan
	@terraform apply tfplan

tf-destroy: ## Destroy infrastructure
	@terraform destroy

local-run: ## Run locally for testing
	@go run cmd/main.go