#!/bin/bash

echo 'Cleaning build directory'
rm -r build

echo '[1/2] Building linux/amd64...'
GOARCH=amd64 GOOS=linux go build -ldflags="-X main.version=$(git describe --always --tags --dirty)" -o ./build/dyngo-amd64 .

echo '[2/2] Building linux/arm64...'
GOARCH=arm64 GOOS=linux go build -ldflags="-X main.version=$(git describe --always --tags --dirty)" -o ./build/dyngo-arm64 .

echo 'Done!'
