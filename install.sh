#!/bin/bash

# Kairoa CLI Installation Script
# Supports: macOS, Linux
# Architectures: amd64, arm64, 386

set -e

REPO="covoyage/kairoa-cli"
BINARY_NAME="kairoa"
INSTALL_DIR="/usr/local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Detect OS
 detect_os() {
    case "$(uname -s)" in
        Linux*)     OS=linux;;
        Darwin*)    OS=darwin;;
        CYGWIN*|MINGW*|MSYS*) OS=windows;;
        *)          echo "${RED}Unsupported operating system${NC}"; exit 1;;
    esac
    echo "Detected OS: $OS"
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)   ARCH=amd64;;
        arm64|aarch64)  ARCH=arm64;;
        i386|i686)      ARCH=386;;
        *)              echo "${RED}Unsupported architecture: $(uname -m)${NC}"; exit 1;;
    esac
    echo "Detected architecture: $ARCH"
}

# Get latest release version
get_latest_version() {
    echo "Fetching latest version..."
    VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        echo "${RED}Failed to fetch latest version${NC}"
        exit 1
    fi
    echo "Latest version: $VERSION"
}

# Download and install
download_and_install() {
    local platform="${OS}_${ARCH}"
    local ext="tar.gz"
    if [ "$OS" = "windows" ]; then
        ext="zip"
    fi
    
    local filename="${BINARY_NAME}_${VERSION#v}_${platform}.${ext}"
    local download_url="https://github.com/$REPO/releases/download/$VERSION/$filename"
    
    echo "Downloading $filename..."
    echo "URL: $download_url"
    
    # Create temp directory
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"
    
    # Download
    if command -v curl &> /dev/null; then
        curl -L -o "$filename" "$download_url"
    elif command -v wget &> /dev/null; then
        wget -O "$filename" "$download_url"
    else
        echo "${RED}curl or wget is required${NC}"
        exit 1
    fi
    
    # Extract
    echo "Extracting..."
    if [ "$ext" = "zip" ]; then
        unzip -q "$filename"
    else
        tar -xzf "$filename"
    fi
    
    # Install
    echo "Installing to $INSTALL_DIR..."
    if [ -w "$INSTALL_DIR" ]; then
        mv "$BINARY_NAME" "$INSTALL_DIR/"
    else
        echo "${YELLOW}Requesting sudo access to install to $INSTALL_DIR${NC}"
        sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
    fi
    
    # Cleanup
    cd -
    rm -rf "$TMP_DIR"
    
    echo "${GREEN}Successfully installed $BINARY_NAME $VERSION${NC}"
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" &> /dev/null; then
        echo ""
        echo "${GREEN}Installation successful!${NC}"
        echo ""
        echo "Run '$BINARY_NAME --help' to get started"
        echo ""
        "$BINARY_NAME" --version 2>/dev/null || "$BINARY_NAME" --help | head -1
    else
        echo "${RED}Installation failed. Please check your PATH.${NC}"
        exit 1
    fi
}

# Main
main() {
    echo "========================================"
    echo "  Kairoa CLI Installer"
    echo "========================================"
    echo ""
    
    detect_os
    detect_arch
    get_latest_version
    download_and_install
    verify_installation
}

main "$@"
