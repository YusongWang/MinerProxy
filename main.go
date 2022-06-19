package main

import (
	"miner_proxy/cmd"
	"miner_proxy/global"
	"github.com/gin-gonic/gin"
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

	global.Commit = commit
	global.Version = version
	global.Branch = branch
	global.Auther = auther

	// 增大文件描述符上限
	//utils.IncreaseFDLimit()


	// 设置gin
	gin.SetMode(gin.ReleaseMode)
	
	cmd.Execute()
}
