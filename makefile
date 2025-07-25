BINARY_NAME=twaybar
GO_FILES=$(shell find . -name "*.go" -type f)
BUILD_DIR=.

.PHONY: all
all: build

.PHONY: build
build:
	go build -o $(BINARY_NAME) .

.PHONY: run
run:
	go run .

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

.PHONY: deps
deps:
	go mod download
	go mod tidy

.PHONY: vet
vet:
	go vet ./...

.PHONY: dev
dev: fmt vet build

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build    - Build the binary"
	@echo "  run      - Run the application"
	@echo "  clean    - Remove build artifacts"
	@echo "  deps     - Download and tidy dependencies"
	@echo "  vet      - Run go vet"
	@echo "  dev      - Format, vet, and build"
	@echo "  help     - Show this help message"

