package main

import (
	"miner_proxy/cmd"
	_ "miner_proxy/utils"
	"runtime"
)

func main() {
	ballast := make([]byte, 1*1024*1024*1024)
	runtime.KeepAlive(ballast)

	cmd.Execute()
}
