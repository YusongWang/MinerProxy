package main

import (
	"miner_proxy/cmd"
	_ "miner_proxy/utils"
	"runtime"
)

func main() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	ballast := make([]byte, mem.Alloc)
	runtime.KeepAlive(ballast)

	cmd.Execute()
}
