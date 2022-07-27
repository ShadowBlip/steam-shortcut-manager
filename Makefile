# Docker image configuration
BUILDER_NAME ?= golang:1.17-buster
GO_FILES := $(shell find ./ -name '*.go')
BIN := steam-shortcut-manager
SRC_FILES := $(GO_FILES)
UID := $(shell id -u)
GID := $(shell id -g)

# Default target will run build
.PHONY: default
default: build

# Test will run tests against the project to ensure there are no errors.
.PHONY: test
test: build

# Build will compile the project using Docker.
.PHONY: build
build:
ifdef DOCKER
	@echo "Running build container..."
	docker run --rm -u $(UID):$(GID) -v $(shell pwd):/src --workdir /src $(BUILDER_NAME) make bin/$(BIN)
else
	make bin/$(BIN)
endif

# The main binary
bin/$(BIN): $(SRC_FILES)
	mkdir -p bin
	GOCACHE=/tmp go build -o bin/$(BIN) .

# Clean
clean:
	rm -rf bin

# Run the project
.PHONY: run
run:
	go run .

# Install/tidy go dependencies
.PHONY: dep
dep:
	go mod vendor
	go mod tidy
