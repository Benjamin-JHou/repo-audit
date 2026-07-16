#!/bin/bash
set -e

VERSION="0.1.0"
REPO_OWNER="YOUR_GITHUB_USERNAME"
REPO_NAME="YOUR_REPO_NAME"
INSTALL_DIR=""

echo "=== ctxqa Installer ==="

detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$ARCH" in
        x86_64)  ARCH="amd64" ;;
        aarch64) ARCH="arm64" ;;
        armv7*)  ARCH="arm" ;;
    esac

    case "$OS" in
        darwin*)  OS="darwin" ;;
        linux*)   OS="linux" ;;
        *)        echo "Unsupported OS: $OS"; exit 1 ;;
    esac
}

find_install_dir() {
    if [ -d "$HOME/.local/bin" ]; then
        INSTALL_DIR="$HOME/.local/bin"
    elif [ -d "$HOME/bin" ]; then
        INSTALL_DIR="$HOME/bin"
    elif [ -w "/usr/local/bin" ]; then
        INSTALL_DIR="/usr/local/bin"
    else
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
    fi
}

download_binary() {
    BINARY_NAME="ctxqa-$OS-$ARCH"
    DOWNLOAD_URL="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/v$VERSION/$BINARY_NAME"
    TMPDIR=$(mktemp -d)
    TMPFILE="$TMPDIR/ctxqa"

    echo "Downloading from $DOWNLOAD_URL ..."

    if command -v curl &> /dev/null; then
        curl -fsSL "$DOWNLOAD_URL" -o "$TMPFILE"
    elif command -v wget &> /dev/null; then
        wget -q "$DOWNLOAD_URL" -O "$TMPFILE"
    else
        echo "Neither curl nor wget found. Please install one of them."
        exit 1
    fi

    HASH_URL="$DOWNLOAD_URL.sha256"
    EXPECTED_HASH=""

    if command -v curl &> /dev/null; then
        EXPECTED_HASH=$(curl -fsSL "$HASH_URL" 2>/dev/null || echo "")
    elif command -v wget &> /dev/null; then
        EXPECTED_HASH=$(wget -qO- "$HASH_URL" 2>/dev/null || echo "")
    fi

    if [ -n "$EXPECTED_HASH" ]; then
        ACTUAL_HASH=$(sha256sum "$TMPFILE" | awk '{print $1}')
        if [ "$ACTUAL_HASH" != "$EXPECTED_HASH" ]; then
            echo "ERROR: SHA256 checksum mismatch!"
            echo "Expected: $EXPECTED_HASH"
            echo "Actual:   $ACTUAL_HASH"
            rm -rf "$TMPDIR"
            exit 1
        fi
        echo "Checksum verified."
    else
        echo "No hash file found, skipping verification."
    fi

    chmod +x "$TMPFILE"
}

install_binary() {
    cp "$TMPFILE" "$INSTALL_DIR/ctxqa"
    rm -rf "$TMPDIR"

    if [[ "$INSTALL_DIR" != "$PATH"* ]]; then
        export PATH="$INSTALL_DIR:$PATH"
        echo "Added $INSTALL_DIR to PATH for this session."
        echo "Add 'export PATH=\"$INSTALL_DIR:\$PATH\"' to your shell profile for persistence."
    fi
}

verify_install() {
    if command -v ctxqa &> /dev/null || [ -x "$INSTALL_DIR/ctxqa" ]; then
        echo ""
        echo "=== Installation Complete ==="
        echo "Installed to: $INSTALL_DIR/ctxqa"
        echo "Version: $( $INSTALL_DIR/ctxqa --version 2>/dev/null || echo $VERSION )"
        echo ""
        echo "Run 'ctxqa config init' to configure your API key."
        echo "Run 'ctxqa audit' to start auditing your repository."
    else
        echo "ERROR: Installation failed. Could not verify binary."
        exit 1
    fi
}

main() {
    detect_platform
    find_install_dir
    echo "Platform: $OS/$ARCH"
    echo "Installing to: $INSTALL_DIR"
    download_binary
    install_binary
    verify_install
}

main
