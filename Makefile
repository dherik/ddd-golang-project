# Define variables
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST := $(GOCMD) test
BINARY_NAME := main
MAIN_FILE := main.go

dependency:
	go get -v ./...

unit-test: dependency
	$(GOTEST) -v -short  ./...

integration-test: dependency
	$(GOTEST) -coverpkg=./... -coverprofile=coverage.out -v ./...

build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_FILE)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	DB_HOST=127.0.0.1 ./$(BINARY_NAME)

.PHONY: build test clean
