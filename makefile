all :
	go build -o MinerProxy ./src/main.go
prod :
	sh ./build-all.sh