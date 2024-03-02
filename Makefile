LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.2

lint:
	golangci-lint run
.PHONY: install-deps

build:
	docker buildx build --no-cache -t cr.selcloud.ru/registry/sender:latest .
.PHONY: build

