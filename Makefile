ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

DB_USER ?= postgres
DB_PASS ?= postgres
DB_HOST ?= 127.0.0.1
DB_PORT ?= 5432
DB_NAME ?= pvz_avito_db
DB_NAME_TEST ?= pvz_avito_test_db

DB_DSN = postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
DB_DSN_TEST = postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME_TEST)?sslmode=disable

APP_NAME = pvz-avito

.PHONY: build-run run-exe run-go run-tests migrate-up migrate-down migrate-up-test migrate-down-test gen-mocks

build-run:
	go build -o $(APP_NAME) cmd/app/main.go
	./$(APP_NAME)

run-exe:
	./$(APP_NAME)

run-go:
	go run cmd/app/main.go

run-tests:
	go test -v ./internal/auth/application > ./tests/auth_tests.log 2>&1
	go test -v ./internal/pvz/application > ./tests/pvz_tests.log 2>&1
	go test -v ./internal/reception/application > ./tests/reception_tests.log 2>&1
	go test -v ./internal/product/application > ./tests/product_tests.log 2>&1

	go test -v ./tests/integration > ./tests/integration_tests.log 2>&1

migrate-up:
	migrate -database "$(DB_DSN)" -path ./migrations up

migrate-down:
	migrate -database "$(DB_DSN)" -path ./migrations down

migrate-up-test:
	migrate -database "$(DB_DSN_TEST)" -path ./migrations up

migrate-down-test:
	migrate -database "$(DB_DSN_TEST)" -path ./migrations down

gen-mocks:
	mockgen -source=internal/auth/domain/repository.go -destination=internal/auth/mocks/auth_repository_mock.go -package=mocks
	mockgen -source=internal/pvz/domain/repository.go -destination=internal/pvz/mocks/pvz_repository_mock.go -package=mocks
	mockgen -source=internal/reception/domain/repository.go -destination=internal/reception/mocks/reception_repository_mock.go -package=mocks
	mockgen -source=internal/product/domain/repository.go -destination=internal/product/mocks/product_repository_mock.go -package=mocks

	mockgen -source=internal/auth/delivery/http/handler.go -destination=internal/auth/mocks/auth_service_mock.go -package=mocks
	mockgen -source=internal/pvz/delivery/http/handler.go -destination=internal/pvz/mocks/pvz_service_mock.go -package=mocks
	mockgen -source=internal/reception/delivery/http/handler.go -destination=internal/reception/mocks/reception_service_mock.go -package=mocks
	mockgen -source=internal/product/delivery/http/handler.go -destination=internal/product/mocks/product_service_mock.go -package=mocks