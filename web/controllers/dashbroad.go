package controllers

import (
	"miner_proxy/global"
	"miner_proxy/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"go.uber.org/zap"
)

type Dashboard struct {
	PoolLength    int
	OnlineWorker  int
	OfflineWorker int
	TotalHash     string
	OnlineTime    string
	TotalShare    int64
	TotalDiff     int64
	FeeShares     int64
	FeeDiff       int64
	DevShares     int64
	DevDiff       int64
}

// 首页展示数据接口
func Home(c *gin.Context) {
	var data map[string]Dashboard

	//data["ETH"].PoolLength = len(global.ManageApp.Config)
	//data.PoolLength = len(global.ManageApp.Config)

	var eth Dashboard
	var etc Dashboard
	for _, app := range global.ManageApp.Config {
		if app.Coin == "ETH" {
			eth.PoolLength++
			eth.OnlineWorker = eth.OnlineWorker + len(global.OnlinePools[app.ID])
			//eth.OnlineWorker = eth.OnlineWorker + len(global.OnlinePools[app.ID])
		}

		if app.Coin == "ETC" {
			etc.PoolLength++
			etc.OnlineWorker = etc.OnlineWorker + len(global.OnlinePools[app.ID])
			//eth.OnlineWorker = eth.OnlineWorker + len(global.OnlinePools[app.ID])
		}
	}

	data["ETH"] = eth
	data["ETC"] = eth

	c.JSON(200, gin.H{
		"data":    data,
		"message": "",
		"code":    200,
	})
}

// 系统性能接口
func System(c *gin.Context) {
	utils.Logger.Info("cpu", zap.Any("Cpu", GetCpuPercent()))
	utils.Logger.Info("mem", zap.Any("mem", GetMemPercent()))
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func GetMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}
