#!/bin/bash
build() {
    VERSION=$(cat version.txt)
    COMMIT=$(git rev-parse HEAD)
    AUTHER=$(cat auther.txt)
    BRANCH=$(git rev-parse --abbrev-ref HEAD)
    LDFLAGS="-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.auther=${AUTHER} -X main.branch=${BRANCH}"

    os="$1"
    if [ "$os" == "darwin" ]; then os="macos"; fi

    arch="$2"
    if [ "$arch" == "386" ]; then arch="x86"; fi
    if [ "$arch" == "amd64" ]; then arch="amd64"; fi
    if [ "$arch" == "loong64" ]; then arch="loongarch64"; fi
    if [ "$arch" == "mipsle" ]; then arch="mipsel"; fi
    if [ "$arch" == "mips64le" ]; then arch="mips64el"; fi

    ext=""
    if [ "$1" == "windows" ]; then ext=".exe"; fi

    echo "build for $os $arch..."
    CGO_ENABLED=0 GOOS="$1" GOARCH="$2" go build -a -ldflags "-extldflags \"-static\" $LDFLAGS" -o "build/MinerProxy-$os-$arch$ext" main.go
}

cd "$(dirname "$0")"
mkdir -p build

#build linux 386
build linux amd64
#build linux arm
#build linux arm64
#build linux mipsle
#build linux mips64le

#build windows 386
build windows amd64

#build darwin amd64
#build darwin arm64
