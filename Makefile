.PHONY: build clean test install

# Go settings
GOPATH ?= $(HOME)/go
GOMODCACHE ?= $(GOPATH)/pkg/mod

export GOPATH
export GOMODCACHE

BINARY_NAME=kairoa
BUILD_DIR=.

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

clean:
	rm -f $(BUILD_DIR)/$(BINARY_NAME)

test:
	go test -v ./...

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

run: build
	./$(BINARY_NAME)
