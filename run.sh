#!/bin/bash

# Lumine Run Script

echo "🔨 Building Lumine..."

# Build with explicit output
go build -mod=mod -o lumine . 2>&1

BUILD_EXIT=$?

if [ $BUILD_EXIT -eq 0 ] && [ -f "lumine" ]; then
    echo "✅ Build successful!"
    echo ""
    echo "🚀 Starting Lumine..."
    echo ""
    ./lumine
elif [ $BUILD_EXIT -eq 0 ]; then
    echo "✅ Build completed but binary not found in current directory"
    echo "Trying to run with go run..."
    go run -mod=mod .
else
    echo "❌ Build failed with exit code: $BUILD_EXIT"
    exit 1
fi
