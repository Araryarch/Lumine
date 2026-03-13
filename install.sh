#!/bin/bash
# Lumine Installation Script

set -e

REPO="Araryarch/lumine"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="lumine"

echo "🌟 Installing Lumine..."

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "❌ Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Get latest release
echo "📥 Downloading latest release..."
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo "❌ Failed to get latest release"
    exit 1
fi

DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/lumine-${OS}-${ARCH}"

if [ "$OS" = "darwin" ]; then
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/lumine-darwin-${ARCH}"
fi

# Download binary
TMP_FILE="/tmp/lumine"
curl -L "$DOWNLOAD_URL" -o "$TMP_FILE"

# Make executable
chmod +x "$TMP_FILE"

# Install
echo "📦 Installing to $INSTALL_DIR..."
sudo mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"

echo "✅ Lumine installed successfully!"
echo ""
echo "Run 'lumine' to get started."
echo ""
echo "📚 Documentation: https://github.com/$REPO"
