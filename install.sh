#!/usr/bin/env bash

# tmail installation script
set -eu

REPO="jacobbanks/tmail"
BINARY="tmail"
DEFAULT_INSTALL_DIR="$HOME/.local/bin"
INSTALL_DIR="${INSTALL_DIR:-$DEFAULT_INSTALL_DIR}"
GITHUB_API="https://api.github.com/repos/$REPO/releases/latest"
TMP_DIR=""

cleanup() {
  if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ]; then
    echo "🧹 Cleaning up temporary files..."
    cd /
    rm -rf "$TMP_DIR"
  fi
}

trap cleanup EXIT
trap 'echo "❌ Installation failed"; exit 1' ERR

check_dependencies() {
  for cmd in curl grep tar; do
    if ! command -v "$cmd" >/dev/null; then
      echo "❌ Required dependency '$cmd' not found."
      exit 1
    fi
  done
}

parse_args() {
  VERSION=""
  QUIET=false

  while [[ $# -gt 0 ]]; do
    case "$1" in
      -v|--version)
        VERSION="$2"
        shift 2
        ;;
      -q|--quiet)
        QUIET=true
        shift
        ;;
      -h|--help)
        show_help
        exit 0
        ;;
      *)
        echo "❌ Unknown argument: $1"
        show_help
        exit 1
        ;;
    esac
  done
}

show_help() {
  echo "Usage: $0 [options]"
  echo
  echo "Options:"
  echo "  -v, --version VERSION   Install specific version"
  echo "  -q, --quiet             Quiet mode, minimal output"
  echo "  -h, --help              Show this help message"
  echo
  echo "Environment variables:"
  echo "  INSTALL_DIR             Installation directory (default: $DEFAULT_INSTALL_DIR)"
}

# Check for existing installation
check_existing() {
  if command -v "$BINARY" >/dev/null; then
    CURRENT_VERSION=$("$BINARY" version 2>/dev/null || echo "unknown")
    echo "ℹ️ Found existing installation: $CURRENT_VERSION"
    if [ -n "$VERSION" ]; then
      echo "📦 Installing version: $VERSION"
    fi
    echo "⬆️ Would you like to continue? (y/n)"
    read -r response
    if [[ ! "$response" =~ ^[Yy] ]]; then
      echo "🛑 Installation aborted"
      exit 0
    fi
  fi
}

check_permissions() {
  if [ ! -d "$INSTALL_DIR" ]; then
    mkdir -p "$INSTALL_DIR" || {
      echo "❌ Cannot create directory $INSTALL_DIR"
      echo "👉 Try running with sudo or set INSTALL_DIR to a writable location"
      exit 1
    }
  fi

  if [ ! -w "$INSTALL_DIR" ]; then
    echo "❌ No write permission to $INSTALL_DIR"
    echo "👉 Try running with sudo or set INSTALL_DIR to a writable location"
    exit 1
  fi
}

detect_os() {
  OS=$(uname -s)
  case "$OS" in
    Linux)   OS="Linux" ;;
    Darwin)  OS="Darwin" ;;
    *) 
      echo "❌ Unsupported OS: $OS" 
      exit 1 
      ;;
  esac
}

detect_arch() {
  ARCH=$(uname -m)
  case "$ARCH" in
    x86_64) ARCH="x86_64" ;;
    arm64|aarch64)  ARCH="arm64" ;;
    i386|i686) ARCH="i386" ;;
    *) 
      echo "❌ Unsupported architecture: $ARCH" 
      exit 1 
      ;;
  esac
}

get_version() {
  if [ -z "$VERSION" ]; then
    echo "🔍 Fetching latest version..."
    # Use a more portable way to extract version that works on both GNU and BSD grep
    VERSION=$(curl -sL "$GITHUB_API" | grep -o '"tag_name":[^"]*"[^"]*"' | sed 's/"tag_name":[^"]*"\(.*\)"/\1/')
    if [ -z "$VERSION" ]; then
      echo "❌ Could not fetch latest release version."
      exit 1
    fi
    echo "📦 Latest version: $VERSION"
  fi
}

download_and_install() {
  echo "⬇️ Downloading $BINARY $VERSION for $OS $ARCH..."

  # Construct download URL
  TARBALL="${BINARY}_${VERSION}_${OS}_${ARCH}.tar.gz"
  URL="https://github.com/$REPO/releases/download/$VERSION/$TARBALL"

  TMP_DIR=$(mktemp -d)
  cd "$TMP_DIR"

  if [ "$QUIET" = true ]; then
    curl -sL "$URL" -o "$TARBALL"
  else
    curl -#L "$URL" -o "$TARBALL"
  fi

  if [ ! -s "$TARBALL" ]; then
    echo "❌ Failed to download $TARBALL"
    echo "🔗 URL: $URL"
    exit 1
  fi

  # Try to download and verify checksum if available
  if curl -sL "$URL.sha256" -o checksums.txt 2>/dev/null; then
    echo "🔐 Verifying checksum..."
    if command -v sha256sum >/dev/null; then
      SHA256=$(sha256sum "$TARBALL" | cut -d ' ' -f 1)
    elif command -v shasum >/dev/null; then
      SHA256=$(shasum -a 256 "$TARBALL" | cut -d ' ' -f 1)
    else
      echo "⚠️ Skipping checksum verification (no sha256sum or shasum command)"
      SHA256=""
    fi
    
    if [ -n "$SHA256" ] && ! grep -q "$SHA256" checksums.txt; then
      echo "❌ Checksum verification failed"
      exit 1
    fi
    [ -n "$SHA256" ] && echo "✅ Checksum verified"
  fi

  echo "📦 Extracting..."
  tar -xzf "$TARBALL"

  if [ ! -f "$BINARY" ]; then
    echo "❌ Binary not found in archive: $BINARY"
    exit 1
  fi

  echo "🚚 Installing to $INSTALL_DIR..."
  chmod +x "$BINARY"
  mv "$BINARY" "$INSTALL_DIR/"
}

setup_path() {
  if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo "⚠️ $INSTALL_DIR is not in your PATH."
    
    for profile in "$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.profile"; do
      if [ -f "$profile" ]; then
        echo "👉 Would you like to add $INSTALL_DIR to your PATH in $profile? (y/n)"
        read -r response
        if [[ "$response" =~ ^[Yy] ]]; then
          echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$profile"
          echo "✅ Updated $profile - restart your shell or run 'source $profile'"
          return 0
        fi
      fi
    done
    
    echo "👉 Add the following line to your shell profile:"
    echo "export PATH=\"\$PATH:$INSTALL_DIR\""
  else
    echo "🚀 You can now run '$BINARY'!"
  fi
}

main() {
  check_dependencies
  parse_args "$@"
  check_existing
  check_permissions
  detect_os
  detect_arch
  get_version
  download_and_install
  echo "✅ Installed $BINARY $VERSION to $INSTALL_DIR"
  setup_path
}

main "$@"
