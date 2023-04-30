#!/bin/bash

echo 'Cleaning build directory'
rm -r build

echo '[1/3] Building linux/amd64...'
GOARCH=amd64 GOOS=linux go build -ldflags="-X main.version=$(git describe --always --tags --dirty)" -o ./build/dyngo-amd64 .

echo '[2/3] Building linux/arm64...'
GOARCH=arm64 GOOS=linux go build -ldflags="-X main.version=$(git describe --always --tags --dirty)" -o ./build/dyngo-arm64 .

echo '[3/3] Aliasing dyngo-arm64 to dyngo-aarch64...'
cp ./build/dyngo-arm64 ./build/dyngo-aarch64

echo 'Done!'
