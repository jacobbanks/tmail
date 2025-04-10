#!/usr/bin/env bash

set -e

REPO="jacobbanks/tmail"
BINARY="tmail"
INSTALL_DIR="$HOME/.local/bin"
GITHUB_API="https://api.github.com/repos/$REPO/releases/latest"

# Detect OS
OS=$(uname -s)
case "$OS" in
  Linux)   OS="Linux" ;;
  Darwin)  OS="Darwin" ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect Arch
ARCH=$(uname -m)
case "$ARCH" in
  x86_64) ARCH="x86_64" ;;
  arm64)  ARCH="arm64" ;;
  aarch64) ARCH="arm64" ;;
  i386) ARCH="i386" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest version tag from GitHub
echo "🔍 Fetching latest version..."
VERSION=$(curl -sL $GITHUB_API | grep -Po '"tag_name": "\K.*?(?=")')
if [ -z "$VERSION" ]; then
  echo "❌ Could not fetch latest release version."
  exit 1
fi

echo "⬇️ Downloading $BINARY $VERSION for $OS $ARCH..."

# Construct download URL
TARBALL="${BINARY}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$VERSION/$TARBALL"

TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Download and extract
curl -sL "$URL" -o "$TARBALL"
tar -xzf "$TARBALL"

# Move binary to install directory
mkdir -p "$INSTALL_DIR"
chmod +x "$BINARY"
mv "$BINARY" "$INSTALL_DIR/"

echo "✅ Installed $BINARY to $INSTALL_DIR"

# Check if $INSTALL_DIR is in PATH
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
  echo "⚠️ $INSTALL_DIR is not in your PATH."
  echo "👉 Add the following line to your shell profile:"
  echo "export PATH=\"\$PATH:$INSTALL_DIR\""
else
  echo "🚀 You can now run '$BINARY'!"
fi

# Clean up
cd /
rm -rf "$TMP_DIR"
