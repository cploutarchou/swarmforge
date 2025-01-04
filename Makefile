.PHONY: all build install clean test

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary name
BINARY_NAME=infra
BUILD_DIR=build

# Build information
VERSION=1.0.0
BUILD_TIME=$(shell date +%FT%T%z)
GIT_COMMIT=$(shell git rev-parse --short HEAD || echo "unknown")

# Linker flags
LDFLAGS=-ldflags "-X github.com/cploutarchou/swarmforge/cmd.Version=${VERSION} \
                  -X github.com/cploutarchou/swarmforge/cmd.BuildTime=${BUILD_TIME} \
                  -X github.com/cploutarchou/swarmforge/cmd.GitCommit=${GIT_COMMIT}"

all: clean deps test build

build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

test:
	$(GOTEST) -v ./...

deps:
	$(GOMOD) download

# Cross compilation
build-linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64

build-all: build-linux 

# Development
dev: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Release
release: build-all
	tar -czf $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-linux-amd64

.DEFAULT_GOAL := build
