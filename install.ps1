# Lumine Installation Script for Windows

$ErrorActionPreference = "Stop"

$REPO = "Araryarch/lumine"
$INSTALL_DIR = "$env:LOCALAPPDATA\lumine"
$BINARY_NAME = "lumine.exe"

Write-Host "🌟 Installing Lumine..." -ForegroundColor Cyan

# Create install directory
if (-not (Test-Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Path $INSTALL_DIR | Out-Null
}

# Get latest release
Write-Host "📥 Downloading latest release..." -ForegroundColor Cyan
$LATEST_RELEASE = (Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest").tag_name

if (-not $LATEST_RELEASE) {
    Write-Host "❌ Failed to get latest release" -ForegroundColor Red
    exit 1
}

$DOWNLOAD_URL = "https://github.com/$REPO/releases/download/$LATEST_RELEASE/lumine-windows-amd64.exe"

# Download binary
$TMP_FILE = "$env:TEMP\lumine.exe"
Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $TMP_FILE

# Install
Write-Host "📦 Installing to $INSTALL_DIR..." -ForegroundColor Cyan
Move-Item -Path $TMP_FILE -Destination "$INSTALL_DIR\$BINARY_NAME" -Force

# Add to PATH if not already there
$PATH = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($PATH -notlike "*$INSTALL_DIR*") {
    Write-Host "➕ Adding to PATH..." -ForegroundColor Cyan
    [Environment]::SetEnvironmentVariable("PATH", "$PATH;$INSTALL_DIR", "User")
    $env:PATH = "$env:PATH;$INSTALL_DIR"
}

Write-Host "✅ Lumine installed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Run 'lumine' to get started." -ForegroundColor Yellow
Write-Host ""
Write-Host "📚 Documentation: https://github.com/$REPO" -ForegroundColor Cyan
Write-Host ""
Write-Host "⚠️  You may need to restart your terminal for PATH changes to take effect." -ForegroundColor Yellow
