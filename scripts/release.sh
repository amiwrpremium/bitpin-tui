#!/bin/bash

# Set the app name
APP_NAME="bitpin-tui"

# Set the output directory
OUTPUT_DIR="$(pwd)/build"
# Clean up the previous builds
echo "Cleaning up previous builds..."
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Common build flags for size reduction
# Use an array to properly handle multiple arguments
BUILD_FLAGS=(-ldflags="-s -w" -trimpath)

# Build for macOS (Apple Silicon)
echo "Building for macOS (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build "${BUILD_FLAGS[@]}" -o "$OUTPUT_DIR/${APP_NAME}-mac-arm64"

# Build for macOS (Intel)
echo "Building for macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build "${BUILD_FLAGS[@]}" -o "$OUTPUT_DIR/${APP_NAME}-mac-amd64"

# Build for Windows (64-bit)
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build "${BUILD_FLAGS[@]}" -o "$OUTPUT_DIR/${APP_NAME}-windows.exe"

# Build for Linux (64-bit)
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build "${BUILD_FLAGS[@]}" -o "$OUTPUT_DIR/${APP_NAME}-linux"

# Compress the binaries
echo "Compressing binaries..."
upx --best --lzma "$OUTPUT_DIR/${APP_NAME}-mac-arm64"
upx --best --lzma "$OUTPUT_DIR/${APP_NAME}-mac-amd64"
upx --best --lzma "$OUTPUT_DIR/${APP_NAME}-windows.exe"
upx --best --lzma "$OUTPUT_DIR/${APP_NAME}-linux"

# Display the results
echo "Build completed:"
ls -lh "$OUTPUT_DIR"
