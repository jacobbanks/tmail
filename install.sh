#!/bin/bash

set -e

OWNER="jacobbanks"
REPO="tmail"
BINARY_NAME="tmail"

# Where to install the binary
INSTALL_DIR="/usr/local/bin"
FALLBACK_DIR="$HOME/.local/bin"

# Detect OS and ARCH
OS="$(uname -s)"
ARCH="$(uname -m)"

# Normalize for asset naming
case "$OS" in
    Linux*)     OS_NAME="Linux" ;;
    Darwin*)    OS_NAME="Darwin" ;;
    MINGW*|MSYS*|CYGWIN*) OS_NAME="Windows" ;;
    *)          echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
    x86_64)     ARCH_NAME="x86_64" ;;
    arm64|aarch64) ARCH_NAME="arm64" ;;
    i386)       ARCH_NAME="i386" ;;
    *)          echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Determine archive extension
if [[ "$OS_NAME" == "Windows" ]]; then
    EXT="zip"
else
    EXT="tar.gz"
fi

# Get latest release tag
LATEST_TAG=$(curl -s "https://api.github.com/repos/${OWNER}/${REPO}/releases/latest" | grep tag_name | cut -d '"' -f 4)
if [ -z "$LATEST_TAG" ]; then
    echo "Could not fetch the latest release tag."
    exit 1
fi

# Build download URL
ARCHIVE_NAME="${BINARY_NAME}_${OS_NAME}_${ARCH_NAME}.${EXT}"
DOWNLOAD_URL="https://github.com/${OWNER}/${REPO}/releases/download/${LATEST_TAG}/${ARCHIVE_NAME}"

# Temp directory
TMP_DIR=$(mktemp -d)

echo "📦 Downloading $ARCHIVE_NAME from $DOWNLOAD_URL..."
curl -sL "$DOWNLOAD_URL" -o "$TMP_DIR/$ARCHIVE_NAME"

# Extract archive
cd "$TMP_DIR"
echo "📂 Extracting..."

if [[ "$EXT" == "tar.gz" ]]; then
    tar -xzf "$ARCHIVE_NAME"
elif [[ "$EXT" == "zip" ]]; then
    unzip -q "$ARCHIVE_NAME"
else
    echo "❌ Unknown archive format: $EXT"
    exit 1
fi

# Move binary to destination
BIN_PATH="./${BINARY_NAME}"
if [[ ! -f "$BIN_PATH" ]]; then
    # Try to find it if it's inside a folder
    BIN_PATH=$(find . -type f -name "$BINARY_NAME" | head -n 1)
    if [[ -z "$BIN_PATH" ]]; then
        echo "❌ Could not find binary in archive."
        exit 1
    fi
fi

chmod +x "$BIN_PATH"

if [ -w "$INSTALL_DIR" ]; then
    mv "$BIN_PATH" "$INSTALL_DIR/$BINARY_NAME"
    INSTALL_PATH="$INSTALL_DIR/$BINARY_NAME"
else
    echo "⚠️ No write permission for $INSTALL_DIR. Installing to $FALLBACK_DIR..."
    mkdir -p "$FALLBACK_DIR"
    mv "$BIN_PATH" "$FALLBACK_DIR/$BINARY_NAME"
    INSTALL_PATH="$FALLBACK_DIR/$BINARY_NAME"
fi

# Add to PATH if needed
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
