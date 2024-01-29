
.PHONY: install
install:
	go mod tidy
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/golang/mock/mockgen@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
	go install github.com/pressly/goose/v3/cmd/goose@latest

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
	goose -dir=./migrations postgres "postgresql://root:q9ckMfi6xQUc1@10.10.10.201:5501/photos?sslmode=disable" up

.PHONY: run_sync_server
run-sync-server:
	go run cmd/sync_photos_server/main.go