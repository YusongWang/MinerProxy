VERSION ?= $(shell cat version.txt)
AUTHER ?= $(shell cat auther.txt)
COMMIT = $(shell git rev-parse HEAD)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

all :
	go build -a -ldflags '-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.branch=${BRANCH} -X main.auther=${AUTHER}' -o "miner_proxy" main.go
linux :
	CGO_ENABLED=0 GOOS="linux" GOARCH="amd64" go build -a -ldflags '-extldflags "-static" -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.branch=${BRANCH} -X main.auther=${AUTHER}'  -o "build/MinerProxy-linux-amd64" main.go
prod :
	sh ./build-all.sh

vue :
	go-bindata-assetfs -o=asset/asset.go -pkg=asset dist/...
