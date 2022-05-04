all :
	go build -o MinerProxy main.go
linux :
	CGO_ENABLED=0 GOOS="linux" GOARCH="amd64" go build -a -ldflags '-extldflags "-static"' -o "build/MinerProxy-linux-amd64" main.go
prod :
	sh ./build-all.sh
vue :
	go-bindata-assetfs -o=asset/asset.go -pkg=asset dist/...