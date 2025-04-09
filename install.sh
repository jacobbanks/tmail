#!/bin/bash

set -e

OWNER="jacobbanks"
REPO="tmail"
BINARY_NAME="tmail"

# Where to install the binary
INSTALL_DIR="/usr/local/bin"
FALLBACK_DIR="$HOME/.local/bin"

# Detect OS
OS="$(uname -s)"
ARCH="$(uname -m)"

# Normalize OS and ARCH
case "$OS" in
    Linux*)     OS="linux" ;;
    Darwin*)    OS="darwin" ;;
    *)          echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
    x86_64)     ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *)          echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# GitHub API call to get latest release tag
LATEST_TAG=$(curl -s "https://api.github.com/repos/${OWNER}/${REPO}/releases/latest" | grep tag_name | cut -d '"' -f 4)
if [ -z "$LATEST_TAG" ]; then
    echo "Could not fetch the latest release tag."
    exit 1
fi

# Compose download URL
FILENAME="${BINARY_NAME}-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/${OWNER}/${REPO}/releases/download/${LATEST_TAG}/${FILENAME}"

# Download binary
echo "📦 Downloading ${BINARY_NAME} ${LATEST_TAG} for ${OS}/${ARCH}..."
curl -L "$DOWNLOAD_URL" -o "$BINARY_NAME"
chmod +x "$BINARY_NAME"

# Try to install to INSTALL_DIR
if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    INSTALL_PATH="$INSTALL_DIR/$BINARY_NAME"
else
    echo "⚠️ No write permission for $INSTALL_DIR. Installing to $FALLBACK_DIR..."
    mkdir -p "$FALLBACK_DIR"
    mv "$BINARY_NAME" "$FALLBACK_DIR/$BINARY_NAME"
    INSTALL_PATH="$FALLBACK_DIR/$BINARY_NAME"
fi

# Ensure it's in PATH
if ! command -v "$BINARY_NAME" &> /dev/null; then
    echo "🔧 Adding $FALLBACK_DIR to PATH in your shell profile..."

    SHELL_NAME=$(basename "$SHELL")
    PROFILE_FILE=""

    case "$SHELL_NAME" in
        bash) PROFILE_FILE="$HOME/.bashrc" ;;
        zsh)  PROFILE_FILE="$HOME/.zshrc" ;;
        fish) PROFILE_FILE="$HOME/.config/fish/config.fish" ;;
        *)    PROFILE_FILE="$HOME/.profile" ;;
    esac

    if ! grep -q "$FALLBACK_DIR" "$PROFILE_FILE"; then
        echo "export PATH=\"\$PATH:$FALLBACK_DIR\"" >> "$PROFILE_FILE"
        echo "✅ Added $FALLBACK_DIR to PATH in $PROFILE_FILE. Please restart your shell or run:"
        echo "source $PROFILE_FILE"
    fi
fi

echo "✅ Installed $BINARY_NAME to $INSTALL_PATH"
echo "🚀 Run '$BINARY_NAME --help' to get started!"
