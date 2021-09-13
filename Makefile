GOOS ?= linux
GOARCH ?= amd64
BIN_OUTPUT ?= bin/api-server
MAIN_PATH ?= cmd/server.go
DSN ?= mysql://root:root@(mysqld)/api-server
SQLBOILER_OUTPUT ?= app/models

.PHONY: build
build: 
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BIN_OUTPUT) $(MAIN_PATH)

.PHONY: migrate-up
migrate-up:
	docker-compose exec app migrate -database "$(DSN)" -path migrate/. up

.PHONY: migrate-down
migrate-down:
	docker-compose exec app migrate -database "$(DSN)" -path migrate/. down

.PHONY: gen
gen:
	go generate ./...

