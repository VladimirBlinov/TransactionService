.PHONY: build api
build:
		go build -v ./cmd/api

.PHONY: build worker
build:
		go build -v ./cmd/worker

.PHONY: test
test:
		go test -v -race -timeout 30s ./internal/...

.DEFAULT_GOAL := build worker