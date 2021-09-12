GOOS ?= linux
GOARCH ?= amd64
BIN_OUTPUT ?= bin/api-server

.PHONY: build
build: 
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BIN_OUTPUT)
