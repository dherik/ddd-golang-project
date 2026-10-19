# Define variables
GOCMD := go
GOINSTALL := $(GOCMD) install
GOGET := $(GOCMD) get
GOBUILD := $(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST := $(GOCMD) test
BINARY_NAME := main
MAIN_FILE := main.go

load-test-create-tasks:
	K6_WEB_DASHBOARD=true k6 run tests/load/k6_create_tasks.js

# upgrade all dependencies to the latest version available
# 'go mod tidy' ensures go.mod is sync with source code and clean it up unused deps
upgrade-deps-latest-version:
	$(GOGET) -u all
	go mod tidy

# upgrade all dependencies only to the latest patch version. This is useful when you 
# want to ensure that your project is up-to-date with the latest security patches and bug fixes, 
# without introducing any breaking changes.
upgrade-deps-latest-patch:
	$(GOGET) -u=patch all
	go mod tidy

vulnerability-check:
	GOTOOLCHAIN=go1.27.2 $(GOCMD) run golang.org/x/vuln/cmd/govulncheck@latest ./...

# apply the modernizer fixes suggested by 'go fix'. run 'go fix -diff ./...'
# first to preview the changes without applying them.
fix:
	$(GOCMD) fix ./...

dependency:
	$(GOCMD) mod download

unit-test: dependency
	$(GOTEST) -v -short ./...

integration-test: dependency
	$(GOTEST) -coverpkg=./... -coverprofile=coverage.out -v ./...

build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_FILE)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run: build
	DB_HOST=127.0.0.1 BROKER_HOST=127.0.0.1 ./$(BINARY_NAME)

.PHONY: build test clean fix
