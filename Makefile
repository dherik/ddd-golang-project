# Define variables
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
BINARY_NAME := ddd-golang-project
MAIN_FILE := cmd/main.go

# Define targets and their commands
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_FILE)

test:
	$(GOTEST) ./...

clean:
	rm -f $(BINARY_NAME)

.PHONY: build test clean
