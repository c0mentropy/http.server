BINARY_NAME = server

LDFLAGS = -s -w -extldflags "-static"

OS_ARCH = \
	linux/amd64 \
	linux/386 \
	linux/arm64 \
	linux/arm \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/386

.PHONY: all build clean cross compress help

all: build

build:
	@echo "Building for local platform..."
	CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o $(BINARY_NAME) .

clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME) $(foreach o,$(OS_ARCH),bin/$(subst /,_,$(o))/$(BINARY_NAME)$(if $(filter windows%,$(o)),.exe,)) $(foreach o,$(OS_ARCH),bin/$(subst /,_,$(o))/upx_$(BINARY_NAME)$(if $(filter windows%,$(o)),.exe,))

cross:
	@echo "Cross compiling for multiple platforms..."
	@for target in $(OS_ARCH); do \
		OS=$${target%/*}; \
		ARCH=$${target#*/}; \
		outdir=bin/$${OS}_$${ARCH}; \
		mkdir -p $${outdir}; \
		echo "Building for $$OS/$$ARCH..."; \
		CGO_ENABLED=0 GOOS=$$OS GOARCH=$$ARCH go build -ldflags '$(LDFLAGS)' -o $${outdir}/$(BINARY_NAME)$(if $(filter windows,$$OS),.exe,) . || exit 1; \
	done

compress:
	@echo "Compressing binaries with upx..."
	@if ! command -v upx > /dev/null 2>&1; then \
		echo "Error: upx not found, please install upx first."; \
		exit 1; \
	fi
	@for dir in bin/*/*; do \
		if [ -f $$dir/$(BINARY_NAME) ] || [ -f $$dir/$(BINARY_NAME).exe ]; then \
			f=$$(basename $$dir); \
			if [ -f $$dir/$(BINARY_NAME) ]; then \
				echo "Compressing $$dir/$(BINARY_NAME)..."; \
				upx -o $$dir/upx_$(BINARY_NAME) $$dir/$(BINARY_NAME); \
			elif [ -f $$dir/$(BINARY_NAME).exe ]; then \
				echo "Compressing $$dir/$(BINARY_NAME).exe..."; \
				upx -o $$dir/upx_$(BINARY_NAME).exe $$dir/$(BINARY_NAME).exe; \
			fi \
		fi \
	done

help:
	@echo "Usage:"
	@echo "  make build      # Build for local platform (static, minimized)"
	@echo "  make cross      # Cross compile for all supported OS/ARCH"
	@echo "  make compress   # Compress all cross-compiled binaries with upx"
	@echo "  make clean      # Remove all built binaries and compressed files"
	@echo "  make help       # Show this help message"
