
.PHONY: install
install:
	go mod tidy
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go get github.com/GaijinEntertainment/go-exhaustruct/v3/cmd/exhaustruct
	go install github.com/golang/mock/mockgen@latest

.PHONY: generate
generate: ## Compile message protobuf and gRPC service files.
	buf generate

.PHONY: lint-full
lint-full:
	goimports -w ./internal/..
	goimports -w ./pkg/common/..
	golangci-lint run --config=.golangci.yaml ./...

.PHONY: mock
mock:
	# Common
	mockgen -source=pkg/common/log/logger.go -destination=pkg/common/log/mock/logger_mock.go
	mockgen -source=pkg/common/uid/generator.go -destination=pkg/common/uid/mock/generator_mock.go
	# Service
	mockgen -source=internal/service/auth/service.go -destination=internal/service/auth/mock/service_mock.go

migrate-create\:%:
	goose -dir=./migrations create $* sql

migrate-create\:%:
	goose -dir=./migrations create $* sql

.PHONY: migrate-up
migrate-up:
	goose -dir=./migrations postgres "postgresql://localhost:5432/uritmix?sslmode=disable" up

.PHONY: run
run:
	go run cmd/app/main.go
