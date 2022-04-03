#!/bin/bash
build() {
    # VERSION ?= $(shell cat version.txt)
    # COMMIT = $(shell git rev-parse HEAD)
    # LDFLAGS = -ldflags "-X cmd.version=${VERSION} -X cmd.commit=${COMMIT}"

    os="$1"
    if [ "$os" == "darwin" ]; then os="macos"; fi

    arch="$2"
    if [ "$arch" == "386" ]; then arch="x86"; fi
    if [ "$arch" == "amd64" ]; then arch="x64"; fi
    if [ "$arch" == "loong64" ]; then arch="loongarch64"; fi
    if [ "$arch" == "mipsle" ]; then arch="mipsel"; fi
    if [ "$arch" == "mips64le" ]; then arch="mips64el"; fi

    ext=""
    if [ "$1" == "windows" ]; then ext=".exe"; fi

    echo "build for $os $arch..."
    CGO_ENABLED=0 GOOS="$1" GOARCH="$2" go build -a -ldflags '-extldflags "-static"' -o "build/MinerProxy-$os-$arch$ext" main.go
}

cd "$(dirname "$0")"
mkdir -p build

build linux 386
build linux amd64
build linux arm
build linux arm64
build linux mipsle
build linux mips64le

build windows 386
build windows amd64

build darwin amd64
build darwin arm64
