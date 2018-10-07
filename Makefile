.DEFAULT_GOAL := build

# Build app
build:
	go build
.PHONY: build

# Install app
install:
	go install
.PHONY: install

# Run tests
test:
	@go test -v -race $(shell go list ./...)
.PHONY: test
