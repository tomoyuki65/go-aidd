APP_NAME := aidd
CMD_PATH := cmd/aidd/main.go
BIN_PATH := src/bin/$(APP_NAME)

# Declaration of command-only targets
.PHONY: build-docker build-mac build-linux build-windows run-mac run-linux run-windows

# Command for building
build-docker:
	docker compose build --no-cache

build-mac:
	docker compose run --rm app sh -c \
	  "GOOS=darwin GOARCH=$(shell uname -m) go build -o bin/$(APP_NAME)-mac $(CMD_PATH)"

build-linux:
	docker compose run --rm app sh -c \
	  "GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME)-linux $(CMD_PATH)"

build-windows:
	docker compose run --rm app sh -c \
	  "GOOS=windows GOARCH=amd64 go build -o bin/$(APP_NAME).exe $(CMD_PATH)"

# Command for running
run-mac:
	@test -f $(BIN_PATH)-mac || { echo "Please run 'make build-mac' first."; exit 1; }
	$(BIN_PATH)-mac

run-linux:
	@test -f $(BIN_PATH)-linux || { echo "Please run 'make build-linux' first."; exit 1; }
	$(BIN_PATH)-linux

run-windows:
	@test -f $(BIN_PATH).exe || { echo "Please run 'make build-windows' first."; exit 1; }
	$(BIN_PATH).exe
