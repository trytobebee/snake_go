#!/bin/bash

# Snake Game Build Script
# Creates executables for multiple platforms

echo "🐍 Building Snake Game..."

# Create dist directory
mkdir -p dist

# Build for Mac (Apple Silicon)
echo "Building for Mac (ARM64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/snake_game_mac_arm64 ./cmd/snake

# Build for Mac (Intel)
echo "Building for Mac (AMD64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/snake_game_mac_amd64 ./cmd/snake

# Build for Windows
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/snake_game_windows.exe ./cmd/snake

# Build for Linux
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/snake_game_linux ./cmd/snake

echo ""
echo "✅ Build complete! Files in ./dist:"
ls -lh dist/
