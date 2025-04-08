.PHONY: build clean test install uninstall release

# Default target
all: build

# Build the tmail binary with version information
build:
	@echo "Building tmail..."
	@go build -o tmail \
		-ldflags "-X main.BuildDate=`date -u +%Y-%m-%dT%H:%M:%SZ` \
		-X main.GitCommit=`git rev-parse --short HEAD` \
		-X main.GitState=`if [ -n "$$(git status --porcelain 2>/dev/null)" ]; then echo "dirty"; else echo "clean"; fi`" \
		main.go

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f tmail

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Install to $GOPATH/bin or ~/go/bin
install: build
	@echo "Installing tmail..."
	@go install
	@echo "tmail installed to $(shell go env GOPATH)/bin/tmail"
	@echo "Make sure $(shell go env GOPATH)/bin is in your PATH"

# Install to /usr/local/bin (requires sudo)
install-global: build
	@echo "Installing tmail globally..."
	@sudo cp tmail /usr/local/bin/
	@echo "tmail installed to /usr/local/bin/tmail"

# Uninstall from $GOPATH/bin
uninstall:
	@echo "Uninstalling tmail..."
	@rm -f $(shell go env GOPATH)/bin/tmail

# Uninstall from /usr/local/bin
uninstall-global:
	@echo "Uninstalling tmail globally..."
	@sudo rm -f /usr/local/bin/tmail

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run linter
lint:
	@echo "Linting code..."
	@golint ./...

# Generate release binaries for multiple platforms
release:
	@echo "Building release binaries..."
	@mkdir -p build
	@VERSION=$$(grep -oP 'Version   = "\K[^"]+' main.go)
	@BUILD_LDFLAGS="-X main.BuildDate=`date -u +%Y-%m-%dT%H:%M:%SZ` -X main.GitCommit=`git rev-parse --short HEAD` -X main.GitState=`if [ -n "$$(git status --porcelain 2>/dev/null)" ]; then echo "dirty"; else echo "clean"; fi`"
	
	@echo "Building for macOS (arm64)..."
	@GOOS=darwin GOARCH=arm64 go build -ldflags "$$BUILD_LDFLAGS" -o build/tmail-$$VERSION-darwin-arm64 main.go
	
	@echo "Building for macOS (amd64)..."
	@GOOS=darwin GOARCH=amd64 go build -ldflags "$$BUILD_LDFLAGS" -o build/tmail-$$VERSION-darwin-amd64 main.go
	
	@echo "Building for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 go build -ldflags "$$BUILD_LDFLAGS" -o build/tmail-$$VERSION-linux-amd64 main.go
	
	@echo "Building for Linux (arm64)..."
	@GOOS=linux GOARCH=arm64 go build -ldflags "$$BUILD_LDFLAGS" -o build/tmail-$$VERSION-linux-arm64 main.go
	
	@echo "Building for Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 go build -ldflags "$$BUILD_LDFLAGS" -o build/tmail-$$VERSION-windows-amd64.exe main.go
	
	@echo "Creating checksums..."
	@cd build && sha256sum tmail-$$VERSION-* > checksums-$$VERSION.txt
	
	@echo "Release binaries created in ./build directory"