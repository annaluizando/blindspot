.PHONY: build clean install

BINARY_NAME=blindspot
MAIN_PATH=./cmd/game
GOBIN=$(shell go env GOPATH)/bin

# Build the application
build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

# Clean build files
clean:
	go clean
	rm -f $(BINARY_NAME)

# Install the application to go/bin
install: build
	@mkdir -p $(GOBIN)
	mv $(BINARY_NAME) $(GOBIN)/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) to $(GOBIN)/"
	@echo "Make sure $(GOBIN) is in your PATH to run '$(BINARY_NAME)' from anywhere."