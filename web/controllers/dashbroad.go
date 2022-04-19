package controllers

import (
	"miner_proxy/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"go.uber.org/zap"
)

// 首页展示数据接口
func Home(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
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
