
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
	go install github.com/pav5000/smartimports/cmd/smartimports@v0.2.0
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

.PHONY: generate
generate:
	buf generate --template buf.gen.yaml --path proto/v1/common.proto --path proto/v1/photo_library_common.proto --path proto/v1/photo_groups.proto --path proto/v1/photo_metadata.proto --path proto/v1/photo_tags.proto --path proto/v1/auth.proto
	move swagger\api.swagger.json swagger\photo_library.swagger.json

	buf generate --template buf.gen.yaml --path proto/v1/common.proto --path proto/v1/sync_photos_service.proto
	move swagger\api.swagger.json swagger\sync_photos.swagger.json

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

.PHONY: migrate-up
migrate-up:
	goose -dir=./migrations postgres "postgresql://postgres@localhost:5432/photo_library?sslmode=disable" up

.PHONY: run-sync-server
run-sync-server:
	go run cmd/sync_photos_server/main.go

.PHONY: processing
processing:
	go run cmd/photo_processing/main.go

.PHONY: run-photos-server
run-photos-server:
	go run cmd/photos_server/main.go

.PHONY: format
format:
	smartimports -local "github.com/kkiling/photo-library/backend/api/"

.PHONY: db
db:
	psql -U postgres -h localhost -p 5432 -d postgres -c "DROP DATABASE IF EXISTS photo_library;" -c "CREATE DATABASE photo_library;"
	make migrate-up

.PHONY: db-test
db-test:
	psql -U postgres -h localhost -p 5432 -d postgres -c "DROP DATABASE IF EXISTS photo_library_test;" -c "CREATE DATABASE photo_library_test;"
	goose -dir=./migrations postgres "postgresql://postgres@localhost:5432/photo_library_test?sslmode=disable" up

.PHONY: schema
schema:
	pg_dump -U postgres -h localhost -p 5432 -d photo_library --no-owner --no-privileges --no-tablespaces --no-security-labels --no-comments -s  > schema.sql
	sqlc generate

.PHONY: test
test:
	go test ./...
