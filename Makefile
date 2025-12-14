.PHONY: build install clean test

# Build the binary
build:
	go build -o forge-deploy main.go

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o forge-deploy-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o forge-deploy-linux-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build -o forge-deploy-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o forge-deploy-darwin-arm64 main.go

# Install locally
install:
	go install

# Clean build artifacts
clean:
	rm -f forge-deploy forge-deploy-*

# Download dependencies
deps:
	go mod download
	go mod tidy
