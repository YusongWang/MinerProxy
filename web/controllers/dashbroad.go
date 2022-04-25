package controllers

import (
	"math/big"
	"miner_proxy/global"
	"miner_proxy/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"go.uber.org/zap"
)

type Dashboard struct {
	PoolLength    int      `json:"pool_length"`
	OnlineWorker  int      `json:"online_worker"`
	OfflineWorker int      `json:"offline_worker"`
	TotalHash     *big.Int `json:"total_hash"`
	OnlineTime    string   `json:"online_time"`
	TotalShare    int64    `json:"total_shares"`
	TotalDiff     *big.Int `json:"total_diff"`
	FeeShares     int64    `json:"fee_shares"`
	FeeDiff       *big.Int `json:"fee_diff"`
	DevShares     int64    `json:"dev_shares"`
	DevDiff       *big.Int `json:"dev_diff"`
}

func newDashborad() *Dashboard {

	return &Dashboard{
		PoolLength:    0,
		OnlineWorker:  0,
		OfflineWorker: 0,
		TotalHash:     new(big.Int).SetInt64(0),
		OnlineTime:    "",
		TotalShare:    0,
		TotalDiff:     new(big.Int).SetInt64(0),
		FeeShares:     0,
		FeeDiff:       new(big.Int).SetInt64(0),
		DevDiff:       new(big.Int).SetInt64(0),
	}
}

// 首页展示数据接口
func Home(c *gin.Context) {

	eth := newDashborad()
	etc := newDashborad()

	for _, app := range global.ManageApp.Config {
		if app.Coin == "ETH" {
			eth.PoolLength++
			for _, w := range global.OnlinePools[app.ID] {
				if w.IsOnline() {
					eth.OnlineWorker++
					eth.TotalShare = eth.TotalShare + int64(w.Worker_share)
					eth.FeeShares = eth.FeeShares + int64(w.Fee_idx)
					eth.DevShares = eth.DevShares + int64(w.Dev_idx)
					eth.TotalHash = new(big.Int).Add(eth.TotalHash, w.Report_hash)
					eth.TotalDiff = new(big.Int).Div(new(big.Int).Add(eth.TotalDiff, w.Worker_diff), new(big.Int).SetInt64(2))
				} else {
					eth.OfflineWorker++
				}
			}
		}

		if app.Coin == "ETC" {
			etc.PoolLength++
			for _, w := range global.OnlinePools[app.ID] {

				if w.IsOnline() {
					etc.OnlineWorker++
					etc.TotalShare = etc.TotalShare + int64(w.Worker_share)
					etc.FeeShares = etc.FeeShares + int64(w.Fee_idx)
					etc.DevShares = etc.DevShares + int64(w.Dev_idx)
					etc.TotalHash = new(big.Int).Add(etc.TotalHash, w.Report_hash)
					etc.TotalDiff = new(big.Int).Div(new(big.Int).Add(etc.TotalDiff, w.Worker_diff), new(big.Int).SetInt64(2))
				} else {
					etc.OfflineWorker++
				}

			}
		}
	}

	var data = map[string]*Dashboard{"ETH": eth, "ETC": etc}

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
