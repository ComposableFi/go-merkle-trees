#!/usr/bin/env bash


install_tooling: ## Install linters
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.39.0
	go install gotest.tools/gotestsum@v1.6.4

lint: ## Run linters.
	which golangci-lint || ( \
		make install_tooling \
	)
	golangci-lint run --deadline=3m --config .golangci.yml ./...