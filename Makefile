VERSION := 0.1.0
BINARY_NAME := ctxqa
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

.PHONY: build release clean

build:
	@echo "Building $(BINARY_NAME) v$(VERSION) ..."
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d/ -f1); \
		GOARCH=$$(echo $$platform | cut -d/ -f2); \
		OUTPUT="dist/$(BINARY_NAME)-$$GOOS-$$GOARCH"; \
		echo "  Building for $$GOOS/$$GOARCH ..."; \
		env GOOS=$$GOOS GOARCH=$$GOARCH CGO_ENABLED=0 go build -o "$$OUTPUT" ./cmd/ctxqa/; \
		chmod +x "$$OUTPUT"; \
	done
	@echo ""
	@echo "Build complete. Binaries in dist/"
	@ls -la dist/

clean:
	rm -rf dist/
