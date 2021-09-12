GOOS ?= linux
GOARCH ?= amd64
BIN_OUTPUT ?= bin/api-server
DSN ?= mysql://root:root@(mysqld)/api-server

.PHONY: build
build: 
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BIN_OUTPUT)

.PHONY: migrate-up
migrate-up:
	docker-compose run --rm migrate migrate -database "$(DSN)" -path . up

.PHONY: migrate-down
migrate-down:
	docker-compose run --rm migrate migrate -database "$(DSN)" -path . down
