.PHONY: build run install clean test fmt vet

build:
	@echo "Building Lumine..."
	@go build -ldflags="-s -w" -o bin/lumine main.go
	@echo "✓ Build complete: bin/lumine"

run:
	@go run main.go

install:
	@go install

clean:
	@rm -rf bin/
	@echo "✓ Cleaned"

test:
	@go test -v ./...

fmt:
	@go fmt ./...

vet:
	@go vet ./...

lint:
	@golangci-lint run

deps:
	@go mod download
	@go mod tidy

dev: fmt vet
	@go run main.go

release: clean fmt vet test build
	@echo "✓ Release build complete"
