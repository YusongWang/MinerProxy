package main

import (
	"miner_proxy/cmd"
	_ "miner_proxy/utils"
	"runtime"
)

func main() {
	ballast := make([]byte, 0.5*1024*1024*1024)
	runtime.KeepAlive(ballast)
	//runtime.GOMAXPROCS(4)

	cmd.Execute()
}
