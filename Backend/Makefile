.PHONY: build_api
build_api:
		go build -v ./cmd/api

.PHONY: build_worker
build_worker:
		go build -v ./cmd/worker

.PHONY: test
test:
		go test -v -race -timeout 30s ./internal/...

.DEFAULT_GOAL := build_worker