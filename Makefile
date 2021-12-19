GOOS ?= linux
GOARCH ?= amd64
BIN_OUTPUT ?= bin/api-server
MAIN_PATH ?= cmd/server.go
SQLBOILER_OUTPUT ?= app/models

.PHONY: setup
setup:
	go mod tidy 
	docker compose build
	docker compose up -d && docker compose stop

.PHONY: build
build: 
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BIN_OUTPUT) $(MAIN_PATH)

.PHONY: gen
gen:
	docker compose start
	go generate ./...

