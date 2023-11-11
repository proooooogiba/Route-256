ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=test password=test dbname=test host=localhost port=5432 sslmode=disable
endif

INTERNAL_PKG_PATH=$(CURDIR)/internal/pkg
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/db/migrations

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

.PHONY: test-migration-up
test-migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: test-migration-down
test-migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down

.PHONY: integration-test
integration-test:
	go test -tags=integration -count=1 -v -cover ./tests

.PHONY: unit-test
unit-test:
	go test -v ./... -count=1

.PHONY: unit-test-coverage
unit-test-coverage:
	go test -v ./... -coverprofile=coverage.out

build:
	docker-compose build

up-all:
	docker-compose up -d postgres zookeeper kafka1 kafka2 kafka3

down:
	docker-compose down

run:
	go run ./cmd/web/*.go
gen-proto:
	protoc --go_out=. --go-grpc_out=. api/hotel.proto

gen-proto-REST:
	protoc -I . --grpc-gateway_out ./internal/pkg/pb/ \
        --grpc-gateway_opt paths=source_relative \
        --grpc-gateway_opt generate_unbound_methods=true \
        ./api/hotel.proto