APP_NAME := filecli
SRC := ./main.go
OUTDIR := bin

.PHONY: help build clean all build-linux-amd64 build-windows-386 build-windows-amd64 build-darwin-arm64

help:
	@echo "Makefile commands:"
	@echo "  make build OS=<os> ARCH=<arch>   # Build for specified OS and ARCH"
	@echo "  make build-linux-amd64            # Build for Linux amd64"
	@echo "  make build-windows-386            # Build for Windows 386"
	@echo "  make build-windows-amd64            # Build for Windows amd64"
	@echo "  make build-darwin-arm64           # Build for macOS arm64"
	@echo "  make all                        # Build for all common platforms"
	@echo "  make clean                       # Remove built binaries"

build:
	@if [ -z "$(OS)" ] || [ -z "$(ARCH)" ]; then \
		echo "Error: OS and ARCH variables must be set. Example:"; \
		echo "  make build OS=linux ARCH=amd64"; \
		exit 1; \
	fi
	@mkdir -p $(OUTDIR)
	@echo "Building $(APP_NAME) for OS=$(OS) ARCH=$(ARCH)..."
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -trimpath -ldflags="-s -w" -o $(OUTDIR)/$(APP_NAME)-$(OS)-$(ARCH) $(SRC)

build-linux-amd64:
	@mkdir -p $(OUTDIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o $(OUTDIR)/$(APP_NAME)-linux-amd64 $(SRC)

build-windows-386:
	@mkdir -p $(OUTDIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -trimpath -ldflags="-s -w" -o $(OUTDIR)/$(APP_NAME)-windows-386.exe $(SRC)

build-windows-amd64:
	@mkdir -p $(OUTDIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o $(OUTDIR)/$(APP_NAME)-windows-amd64.exe $(SRC)

build-darwin-arm64:
	@mkdir -p $(OUTDIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o $(OUTDIR)/$(APP_NAME)-darwin-arm64 $(SRC)

all: build-linux-amd64 build-windows-386 build-windows-amd64 build-darwin-arm64

clean:
	@echo "Cleaning binaries in $(OUTDIR)..."
	rm -f $(OUTDIR)/$(APP_NAME)-*
