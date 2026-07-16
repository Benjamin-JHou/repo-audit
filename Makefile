#!/bin/bash
set -e

BINARY_NAME="ctxqa"
VERSION="0.1.0"
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

echo "Building ctxqa v$VERSION ..."

for platform in "${PLATFORMS[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$platform"
    OUTPUT="dist/${BINARY_NAME}-${GOOS}-${GOARCH}"

    echo "  Building for $GOOS/$GOARCH ..."
    env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build -o "$OUTPUT" ./cmd/ctxqa/
    chmod +x "$OUTPUT"

    if [ "$GOOS" = "windows" ]; then
        mv "$OUTPUT" "${OUTPUT}.exe"
    fi
done

echo ""
echo "Build complete. Binaries in dist/"
ls -la dist/
