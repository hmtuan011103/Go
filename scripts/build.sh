#!/bin/bash
set -e

# Build script for GoStructure App

APP_NAME="app"
BUILD_DIR="./build/bin"
MAIN_PATH="./cmd/app"

# Create build directory
mkdir -p $BUILD_DIR

# Get version from git or default
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

echo "Building ${APP_NAME}..."
echo "Version: ${VERSION}"
echo "Build time: ${BUILD_TIME}"

# Build for current platform
go build -ldflags="${LDFLAGS}" -o "${BUILD_DIR}/${APP_NAME}" ${MAIN_PATH}

echo "Build complete: ${BUILD_DIR}/${APP_NAME}"
