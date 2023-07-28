#!make

-include .env
export

PROJECT_DIR=$(shell pwd)
BUILD_PATH=$(PROJECT_DIR)/build

PATH:=$(PATH):$(BUF_PATH)/bin

generate-ent:
	@echo "\033[32mGenerate Entities\033[m"
	rm -rf ./internal/infra/repo/entity
	sqlboiler mysql --no-tests -d -o ./internal/infra/repo/entity -c .sqlboiler.toml

.PHONY: build
build:
	@echo "\033[32mBuild\033[m"
	go build -o $(BUILD_PATH)/allocator ./cmd/allocator/main.go

.PHONY: lint
lint:
	golangci-lint run ./... -v

.PHONY: test
test:
	@echo "\033[32mTest\033[m"
	@go test ./...

generate-api: ## generate artifacts from swagger files
	@echo "\033[32mGenerate Open API\033[m"
	@rm -rf $(PROJECT_DIR)/pkg/api/openapi/*
	oapi-codegen -generate types -package openapi -o $(PROJECT_DIR)/pkg/api/openapi/types.gen.go $(PROJECT_DIR)/api/openapi/allocator.yaml
	oapi-codegen -generate spec -package openapi -o $(PROJECT_DIR)/pkg/api/openapi/spec.gen.go $(PROJECT_DIR)/api/openapi/allocator.yaml
	oapi-codegen -generate server -package openapi -o $(PROJECT_DIR)/pkg/api/openapi/server.gen.go $(PROJECT_DIR)/api/openapi/allocator.yaml
