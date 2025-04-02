.PHONY: build clean install

BINARY_NAME=blindspot
MAIN_PATH=./cmd/game

# Build the application
build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

# Clean build files
clean:
	go clean
	rm -f $(BINARY_NAME)

# Install the application to go/bin (using GOPATH if defined, otherwise $HOME/go/bin)
install: build
	@mkdir -p $(shell go env GOPATH)/bin
	mv $(BINARY_NAME) $(shell go env GOPATH)/bin/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) to $(shell go env GOPATH)/bin/"
	@echo "Make sure $(shell go env GOPATH)/bin is in your PATH to run '$(BINARY_NAME)' from anywhere."