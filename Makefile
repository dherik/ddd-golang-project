# Define variables
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST := $(GOCMD) test
BINARY_NAME := main
MAIN_FILE := main.go

# Define targets and their commands
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_FILE)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	./$(BINARY_NAME)

.PHONY: build test clean
