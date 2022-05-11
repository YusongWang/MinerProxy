package main

import (
	"miner_proxy/cmd"
	_ "miner_proxy/global"
	_ "miner_proxy/utils"
	"runtime"
)

var (
	version string
	commit  string
	branch  string
	auther  string
)

func main() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	ballast := make([]byte, mem.Alloc)
	runtime.KeepAlive(ballast)

	cmd.Execute()
}
