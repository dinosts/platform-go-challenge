# Variables
BINARY_NAME = api
MAIN_PATH = ./cmd/api
GO = go

.PHONY: test run_no_build

test:
	$(GO) test ./...

run_no_build:
	$(GO) run $(MAIN_PATH)
