#!/bin/bash

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Lumine version
VERSION="v1.0.0"
REPO="Araryarch/lumine"

echo -e "${GREEN}⚡ Lumine Installer${NC}"
echo -e "${YELLOW}Docker Development Manager - Laragon for Linux${NC}"
echo ""

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}✗ Docker is not installed${NC}"
    echo "Please install Docker first: https://docs.docker.com/get-docker/"
    exit 1
fi

echo -e "${GREEN}✓ Docker is installed${NC}"

# Check if Docker is running
if ! docker ps &> /dev/null; then
    echo -e "${RED}✗ Docker is not running${NC}"
    echo "Please start Docker and try again"
    exit 1
fi

echo -e "${GREEN}✓ Docker is running${NC}"

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
    armv7l)
        ARCH="armv7"
        ;;
    *)
        echo -e "${RED}✗ Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo -e "${YELLOW}Detected: $OS/$ARCH${NC}"

# Installation directory
INSTALL_DIR="${DIR:-$HOME/.local/bin}"
mkdir -p "$INSTALL_DIR"

# Download URL
BINARY_NAME="lumine_${VERSION}_${OS}_${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}.tar.gz"

echo -e "${YELLOW}Downloading Lumine ${VERSION}...${NC}"

# Download and extract
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

if command -v curl &> /dev/null; then
    curl -sL "$DOWNLOAD_URL" -o lumine.tar.gz
elif command -v wget &> /dev/null; then
    wget -q "$DOWNLOAD_URL" -O lumine.tar.gz
else
    echo -e "${RED}✗ Neither curl nor wget is installed${NC}"
    exit 1
fi

if [ ! -f lumine.tar.gz ]; then
    echo -e "${RED}✗ Download failed${NC}"
    echo "Building from source instead..."
    
    # Build from source
    if ! command -v go &> /dev/null; then
        echo -e "${RED}✗ Go is not installed${NC}"
        echo "Please install Go: https://golang.org/doc/install"
        exit 1
    fi
    
    echo -e "${YELLOW}Cloning repository...${NC}"
    git clone "https://github.com/${REPO}.git" lumine-src
    cd lumine-src
    
    echo -e "${YELLOW}Building Lumine...${NC}"
    go build -o lumine main.go
    
    mv lumine "$INSTALL_DIR/lumine"
else
    tar -xzf lumine.tar.gz
    mv lumine "$INSTALL_DIR/lumine"
fi

chmod +x "$INSTALL_DIR/lumine"

# Cleanup
cd ~
rm -rf "$TMP_DIR"

echo -e "${GREEN}✓ Lumine installed to $INSTALL_DIR/lumine${NC}"

# Check if install dir is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo -e "${YELLOW}⚠ $INSTALL_DIR is not in your PATH${NC}"
    echo "Add this line to your ~/.bashrc or ~/.zshrc:"
    echo ""
    echo -e "${GREEN}export PATH=\"\$PATH:$INSTALL_DIR\"${NC}"
    echo ""
fi

# Create config directory
CONFIG_DIR="$HOME/.lumine"
mkdir -p "$CONFIG_DIR"

# Create default config if it doesn't exist
if [ ! -f "$CONFIG_DIR/config.yaml" ]; then
    echo -e "${YELLOW}Creating default configuration...${NC}"
    cat > "$CONFIG_DIR/config.yaml" << 'EOF'
# Lumine Configuration

default_php: "8.2"

php_versions:
  - "7.4"
  - "8.0"
  - "8.1"
  - "8.2"
  - "8.3"

services:
  nginx:
    image: "nginx:alpine"
    port: 80
    enabled: true
  
  mysql:
    image: "mysql:8.0"
    port: 3306
    enabled: true
    environment:
      MYSQL_ROOT_PASSWORD: "root"
      MYSQL_DATABASE: "lumine"
  
  redis:
    image: "redis:alpine"
    port: 6379
    enabled: true
  
  mailhog:
    image: "mailhog/mailhog"
    port: 8025
    enabled: true
  
  phpmyadmin:
    image: "phpmyadmin:latest"
    port: 8080
    enabled: true

projects_dir: "~/Projects"
refresh_interval: 2

ui:
  show_icons: true
  color_scheme: "minimal"
EOF
    echo -e "${GREEN}✓ Configuration created at $CONFIG_DIR/config.yaml${NC}"
fi

echo ""
echo -e "${GREEN}✓ Installation complete!${NC}"
echo ""
echo -e "${YELLOW}Quick Start:${NC}"
echo "  1. Run: ${GREEN}lumine${NC}"
echo "  2. Use j/k to navigate"
echo "  3. Press Enter to select"
echo "  4. Press ? for help"
echo ""
echo -e "${YELLOW}Documentation:${NC} https://github.com/${REPO}"
echo ""
echo -e "${GREEN}Happy coding! ⚡${NC}"
