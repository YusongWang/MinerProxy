all :
	go build -o MinerProxy main.go
prod :
	sh ./build-all.sh
