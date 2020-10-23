SHELL = /bin/bash

.PHONY: setup
setup:
	go get github.com/goreleaser/goreleaser

.PHONY: lint
lint: generate
	go vet ./...
	goreleaser check

.PHONY: test
test: generate
	go test ./...

.PHONY: integration-test
integration-test:
	go test -tags=integration ./...

.PHONY: coverage
coverage: generate
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: codecov
codecov:  coverage
	bash <(curl -s https://codecov.io/bash)

.PHONY: generate
generate:
	go generate ./...

.PHONY: build
build: generate
	go build

.PHONY: cross-build-snapshot
cross-build:
	goreleaser --rm-dist --snapshot

.PHONY: install
install:
	go install
